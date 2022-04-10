package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/tile"
	"errors"
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {

	g.handleMouseInput()

	//build river first
	if g.RiverDeck.Remaining() > 0 {
		for g.RiverDeck.Remaining() > 0 {
			if err := g.updateRiverBuild(); err != nil {
				return err
			}
		}

	} else if g.Deck.Remaining() > 0 {
		fmt.Println("Placing", g.Deck.Remaining(), "Tiles")

		for g.Deck.Remaining() > 0 {
			tile, err := g.Deck.Pop()

			if err != nil {
				return err
			}

			possiblePlacements := g.Board.PossibleTilePlacements(&tile)

			if len(possiblePlacements) == 0 {
				return errors.New("No valid placement for tile")
			}

			randomIndex := rand.Intn(len(possiblePlacements))
			randomlySelectedPlacement := possiblePlacements[randomIndex]

			err = g.Board.AddTile(&tile, randomlySelectedPlacement)

			if err != nil {
				err = g.Board.AddTile(&tile, randomlySelectedPlacement)

				panic(fmt.Sprint("Error Placing Tile: ", err))
			}
		}

		fmt.Println("Done Placing Tiles")
	}

	return nil
}

func (g *Game) updateRiverBuild() error {

	//start the river with the river terminus piece (which is always located first in a new deck)
	if g.RiverDeck.Index == 0 {

		rt, err := g.RiverDeck.Pop()

		if err != nil {
			return err
		}

		g.Board.AddTile(&rt, tile.Placement{
			Position:    tile.Position{X: 10, Y: 10},
			Orientation: 0,
		})

		g.lastRiverTile = &rt

		return nil
	}

	riverTile, err := g.RiverDeck.Pop()

	if err != nil {
		return err
	}

	riverPlacement, err := g.getRiverPlacement(&riverTile)

	if err != nil {
		return err
	}

	err = g.Board.AddTile(&riverTile, riverPlacement)

	if err != nil {
		return err
	}

	g.lastRiverTile = &riverTile

	return nil
}

func selectRandomTile(tiles map[string]tile.Tile) tile.Tile {
	keys := make([]string, 0, len(tiles))
	for k := range tiles {
		keys = append(keys, k)
	}

	n := rand.Intn(len(keys))

	randKey := keys[n]

	randTile := tiles[randKey]

	return randTile
}

func (g *Game) findTileForPos(pos tile.Position) (tile.Tile, error) {

	for tileName, t := range g.TileFactory.ReferenceTiles() {
		if orientation, err := g.Board.IsTilePlaceable(&t, pos); err == nil {

			builtTile := g.TileFactory.BuildTile(tileName)
			builtTile.Placement.Position = pos
			builtTile.Placement.Orientation = orientation

			return builtTile, nil
		}
	}

	return tile.Tile{}, errors.New("No Tile Fits this Place")
}

//input state
var mouseDown bool = false
var mouseInitialX, mouseInitialY int

func (g *Game) handleMouseInput() {

	cX, cY := ebiten.CursorPosition()

	//cursor pos mapped to tile grid
	g.HoveredPosition.X = int(math.Floor(float64(cX-g.CameraOffset.X) / g.renderScale / float64(g.baseSize)))
	g.HoveredPosition.Y = int(math.Floor(float64(cY-g.CameraOffset.Y) / g.renderScale / float64(g.baseSize)))

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.HoveredPosition != g.SelectedPosition {
		g.SelectedPosition = g.HoveredPosition

		if _, exists := g.Board.Tiles[g.HoveredPosition]; exists {
			clickedTile := g.Board.Tiles[g.HoveredPosition]

			roads := make([]board.Road, 0)
			for _, rs := range clickedTile.UniqueRoadSegements() {
				rd := board.CompileRoadFromSegment(rs)
				roads = append(roads, rd)
			}

			g.HighlightedRoads = roads

		}

		if _, exists := g.Board.OpenPositions[g.HoveredPosition]; exists {
			t, err := g.findTileForPos(g.SelectedPosition)

			if err == nil {
				g.Board.AddTile(&t, t.Placement)
			}
		}
	}

	//mouse state change tracking
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseDown = false
	}

	//panning support
	if !mouseDown && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseInitialX, mouseInitialY = ebiten.CursorPosition()
		mouseDown = true
	}

	//panning
	if mouseDown {

		g.CameraOffset.X += cX - mouseInitialX
		g.CameraOffset.Y += cY - mouseInitialY

		mouseInitialX = cX
		mouseInitialY = cY
	}

}

func (g *Game) getRiverPlacement(riverTile *tile.Tile) (tile.Placement, error) {

	possiblePlacements := g.Board.PossibleTilePlacements(riverTile)

	permittedPlacements := make([]tile.Placement, 0)

	isCurve := false
	for _, pl := range possiblePlacements {
		connectedFeatures := g.Board.ConnectedFeatures(riverTile, pl)

		for dir, cf := range connectedFeatures {

			connectedTilePos := pl.Position.EdgePos(dir)

			if connectedTilePos != g.lastRiverTile.Placement.Position {
				continue
			}

			//must be a river connection
			if cf == tile.River {
				riverFeature := riverTile.Feature(pl.GridToTileDir(dir))
				//dont let the river turn the same way twice
				if riverFeature.Curve {

					isCurve = true

					if g.lastRiverTurn == 1 {
						g.lastRiverTurn = 0
					}

					//next piece must be 180 degrees out of phase with the last
					nextCurveOrientation := (g.lastRiverTurn + 180) % 360

					//!= 1 means this is the first curve
					if g.lastRiverTurn != 1 && pl.Orientation != nextCurveOrientation {
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
		g.lastRiverTurn = randomlySelectedPlacement.Orientation
	}

	return randomlySelectedPlacement, nil
}
