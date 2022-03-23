package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/loader"
	"beeb/carcassonne/tile"
	"errors"
	"fmt"
	"image"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	TileInfo     loader.TileInfoFile
	Tiles        map[string]tile.Tile
	RiverDeck    []tile.Tile
	Deck         []tile.Tile
	Board        board.Board
	CameraOffset image.Point
}

func CreateGame() Game {
	fmt.Println("Creating Game")

	game := Game{}

	game.Tiles, game.TileInfo = loader.LoadTiles()
	game.Board = board.CreateBoard()
	game.RiverDeck = game.buildRiverDeck()
	game.Deck = game.buildDeck()

	//game.Deck = game.builDeck()

	return game
}

var deckCounter int = 0
var riverCounter int = 0
var lastRiverTile *tile.Tile = nil

//1 is not ever a valid orientation so it will not be a false positive
var lastRiverTurn uint16 = 1

func (g *Game) updateRiverBuild() error {
	start := time.Now()

	if riverCounter == 0 {
		//force add the river starter, which is always at the top of the river deck
		g.Board.AddTile(&g.RiverDeck[0], tile.Placement{
			Position:    tile.Position{X: 10, Y: 10},
			Orientation: 0,
		})
		riverCounter++
		lastRiverTile = &g.RiverDeck[0]

		return nil
	}

	riverTile := g.RiverDeck[riverCounter]

	riverPlacement, err := g.getRiverPlacement(&riverTile)

	if err != nil {
		return err
	}

	err = g.Board.AddTile(&riverTile, riverPlacement)

	if err != nil {
		panic("Err Placing Tile")
	}

	lastRiverTile = &riverTile

	riverCounter++

	elapsed := time.Since(start)
	fmt.Println(riverCounter, elapsed)

	return nil
}

var mouseDown bool = false
var mouseInitialX, mouseInitialY int

func (g *Game) handleMouseInput() {

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		g.Tiles, g.TileInfo = loader.LoadTiles()
		g.Board = board.CreateBoard()
		g.RiverDeck = g.buildRiverDeck()
		g.Deck = g.buildDeck()

		deckCounter = 0
		riverCounter = 0
		lastRiverTile = nil
	}

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseDown = false
	}

	if !mouseDown && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseInitialX, mouseInitialY = ebiten.CursorPosition()
		mouseDown = true
	}

	if mouseDown {
		newX, newY := ebiten.CursorPosition()

		g.CameraOffset.X += newX - mouseInitialX
		g.CameraOffset.Y += newY - mouseInitialY

		mouseInitialX = newX
		mouseInitialY = newY
	}
}

func (g *Game) Update() error {

	g.handleMouseInput()

	//build river first
	if riverCounter < len(g.RiverDeck) {
		if err := g.updateRiverBuild(); err != nil {
			return err
		}
	} else if deckCounter < len(g.Deck) {
		//now place tiles wherever
		start := time.Now()

		tile := g.Deck[deckCounter]

		possiblePlacements := g.Board.PossibleTilePlacements(&tile)

		if len(possiblePlacements) == 0 {
			return errors.New("No valid placement for river tile")
		}

		randomIndex := rand.Intn(len(possiblePlacements))
		randomlySelectedPlacement := possiblePlacements[randomIndex]

		err := g.Board.AddTile(&tile, randomlySelectedPlacement)

		if err != nil {
			panic("Error Placing Tile")
		}

		elapsed := time.Since(start)
		fmt.Println(deckCounter, elapsed)

		deckCounter++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.RenderBoard(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400
}

func (g *Game) getRiverPlacement(riverTile *tile.Tile) (tile.Placement, error) {

	possiblePlacements := g.Board.PossibleTilePlacements(riverTile)

	permittedPlacements := make([]tile.Placement, 0)

	isCurve := false
	for _, pl := range possiblePlacements {
		connectedFeatures := g.Board.ConnectedFeatures(riverTile, pl)

		for dir, cf := range connectedFeatures {

			connectedTilePos := pl.Position.EdgePos(dir)

			if connectedTilePos != lastRiverTile.Placement.Position {
				continue
			}

			//must be a river connection
			if cf == tile.River {
				//dont let the river turn the same way twice
				if riverTile.Feature(pl.TileDirection(dir)).Curve {

					isCurve = true

					//next piece must be 180 degrees out of phase with the last
					nextCurveOrientation := (lastRiverTurn + 180) % 360

					//!= 1 means this is the first curve
					if lastRiverTurn != 1 && pl.Orientation != nextCurveOrientation {
						break
					}
				}

				permittedPlacements = append(permittedPlacements, pl)
			}
		}
	}

	if len(permittedPlacements) == 0 {
		return tile.Placement{}, errors.New("No valid placement for river tile")
	}

	randomIndex := rand.Intn(len(permittedPlacements))
	randomlySelectedPlacement := permittedPlacements[randomIndex]

	if isCurve {
		lastRiverTurn = randomlySelectedPlacement.Orientation
	}

	return randomlySelectedPlacement, nil
}

func deckSize(deck map[string]int) int {
	var deckSize int = 0

	for _, v := range deck {
		deckSize += v
	}

	return deckSize
}

func (g *Game) buildRiverDeck() []tile.Tile {

	deckSize := deckSize(g.TileInfo.RiverDeck.Deck)

	riverTiles := make([]tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range g.TileInfo.RiverDeck.Deck {
		for i := 0; i < quantity; i++ {
			riverTiles[c] = g.Tiles[tileName]
			c++
		}
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(deckSize, func(i, j int) {
		riverTiles[i], riverTiles[j] = riverTiles[j], riverTiles[i]
	})

	//prepend begin tile
	riverTiles = append([]tile.Tile{g.Tiles[g.TileInfo.RiverDeck.Begin]}, riverTiles...)

	//append begin tile
	riverTiles = append(riverTiles, g.Tiles[g.TileInfo.RiverDeck.End])

	return riverTiles
}

func (g *Game) buildDeck() []tile.Tile {

	deckSize := deckSize(g.TileInfo.Deck)

	tiles := make([]tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range g.TileInfo.Deck {
		for i := 0; i < quantity; i++ {
			tiles[c] = g.Tiles[tileName]
			c++
		}
	}

	rand.Seed(time.Now().Unix())
	rand.Shuffle(deckSize, func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})

	return tiles
}
