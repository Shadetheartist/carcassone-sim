package engine

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine/board"
	"beeb/carcassonne/engine/deck"
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/engine/turnStage"
	"errors"
	"fmt"
	"image/color"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

var PLAYER_COLOR_LIST = [...]color.RGBA{
	colornames.White,
	colornames.Red500,
	colornames.Blue500,
	colornames.Green500,
	colornames.Black,
}

type Engine struct {
	Deterministic                 bool
	BoardSize                     int
	GameOver                      bool
	GameBoard                     *board.Board
	GameData                      *data.GameData
	Players                       []*Player
	CurrentPlayerIndex            int
	TilePlacedThisTurn            *tile.Tile
	HeldRefTileGroup              *tile.ReferenceTileGroup
	CurrentPossibleTilePlacements []Placement

	TurnCounter int
	TurnStage   turnStage.TurnStage

	RiverDeck *deck.Deck
	Deck      *deck.Deck

	TileFactory *tile.TileFactory

	TilePlacementManager *TilePlacementManager
}

func NewEngine(gameData *data.GameData, boardSize int, numPlayers int) *Engine {
	if numPlayers > len(PLAYER_COLOR_LIST) {
		panic(fmt.Sprint("too many players, max ", len(PLAYER_COLOR_LIST)))
	}

	engine := &Engine{}

	engine.BoardSize = boardSize
	engine.GameData = gameData
	engine.Players = make([]*Player, numPlayers)
	engine.TileFactory = &tile.TileFactory{}
	engine.TilePlacementManager = NewTilePlacementManager(engine)

	engine.InitGame()

	return engine
}

func (e *Engine) InitGame() {
	for i := 0; i < len(e.Players); i++ {
		playerName := fmt.Sprint("Player ", i)
		e.Players[i] = NewPlayer(playerName, PLAYER_COLOR_LIST[i])
	}

	e.GameBoard = board.NewBoard(e.BoardSize)
	e.RiverDeck = deck.BuildRiverDeck(e.GameData)
	e.Deck = deck.BuildDeck(e.GameData)
	e.GameOver = false
	e.TurnCounter = 0
	e.TilePlacedThisTurn = nil
	e.HeldRefTileGroup = nil
	e.CurrentPossibleTilePlacements = nil
	e.CurrentPlayerIndex = 0
	e.TurnStage = turnStage.Draw

	lastRiverTurn = 1
	lastRiverTile = nil
}

func (e *Engine) Step() {

	if e.GameOver {
		return
	}

	player := e.CurrentPlayer()

	switch e.TurnStage {
	case turnStage.Draw:

		//retry getting possible tiles a few times if we don't have a place to put one
		for i := 0; i < 4; i++ {
			if i == 4 {
				panic("Game cannot continue, nowhere to place tile, tried 3 times.")
			}

			rtg, tileTakeErr := e.TakeNextTile()

			if tileTakeErr != nil {
				panic("Attempted to take a tile from an empty deck.")
			}

			e.HeldRefTileGroup = rtg

			e.CurrentPossibleTilePlacements = e.TilePlacementManager.PossibleTilePlacements(rtg)

			if e.GameBoard.PlacedTileCount > 0 && (e.RiverDeck.Remaining() > 0 || rtg.IsRiverTile()) {
				e.CurrentPossibleTilePlacements, _ = e.restictRiverPlacement()
			}

			//this clause shuffles a tile back in when it is not playable

			// this should never happen to the river
			if len(e.CurrentPossibleTilePlacements) < 1 {
				//replace tile
				e.Deck.Append(e.HeldRefTileGroup)
				e.Deck.Shuffle()
				continue
			}

			break
		}

		e.TurnStage++

	case turnStage.PlaceTile:

		selectedPlacement := player.DeterminePlacement(e.CurrentPossibleTilePlacements, e)

		if selectedPlacement == nil {
			panic("No Placement Determined By Player")
		}

		e.PlaceTile(*selectedPlacement)

		e.CurrentPossibleTilePlacements = nil
		e.HeldRefTileGroup = nil
		e.TurnStage++

	case turnStage.PlaceMeeple:
		e.TurnStage++

	case turnStage.Score:
		e.TurnStage++

	case turnStage.Pass:
		err := e.GoToNextTurn()

		if err != nil {
			e.EndGame()
		}

		e.TurnStage = turnStage.Draw
	}
}

func (e *Engine) PlaceTile(placement Placement) {
	newTile := e.TileFactory.NewTileFromReference(placement.ReferenceTile)
	e.GameBoard.PlaceTile(placement.Position, newTile)

	if newTile.Reference.EdgeSignature.Contains(tile.River) {
		lastRiverTile = newTile

		if newTile.Reference.EdgeSignature.Curving() {
			lastRiverTurn = newTile.Reference.Orientation
		}
	}

}

func (e *Engine) TakeNextTile() (*tile.ReferenceTileGroup, error) {

	currentDeck := e.Deck

	//must draw from river deck until it's empty
	if e.RiverDeck.Remaining() > 0 {
		currentDeck = e.RiverDeck
	}

	rtg, err := currentDeck.Pop()

	if err != nil {
		return nil, err
	}

	return rtg, nil
}

func (e *Engine) GoToNextTurn() error {

	if e.Deck.Remaining() == 0 {
		return errors.New("no more tiles in the deck")
	}

	e.TurnCounter++
	e.CurrentPlayerIndex = (e.CurrentPlayerIndex + 1) % len(e.Players)
	e.TilePlacedThisTurn = nil

	return nil
}

func (e *Engine) EndGame() {
	e.GameOver = true
}

var lastRiverTurn = 1
var lastRiverTile *tile.Tile = nil

func (e *Engine) restictRiverPlacement() ([]Placement, error) {

	permittedPlacements := make([]Placement, 0)

	for _, placement := range e.CurrentPossibleTilePlacements {

		connectedFeatures := placement.ConnectedFeatures

		for _, cf := range connectedFeatures {
			dir := cf.EdgeA

			connectedTilePos := placement.Position.EdgePos(dir)

			//must be connected to the last tile placed
			if connectedTilePos != lastRiverTile.Position {
				continue
			}

			//must be a river connection
			if cf.FeatureA.Type != tile.River {
				continue
			}

			//dont let the river turn the same way twice
			//if the new tile is curving, keep track of which way it was oriented
			//the next curved tile must be oriented 180 deg different from this tile
			if cf.FeatureA.ParentRefenceTileGroup.Orientations[0].EdgeSignature.Curving() {

				//this is the first curve of the river, it can go either way
				if lastRiverTurn == 1 {
					permittedPlacements = append(permittedPlacements, placement)
					continue
				}

				//next piece must be 180 degrees out of phase with the last
				nextCurveOrientation := (lastRiverTurn + 180) % 360

				if placement.ReferenceTile.Orientation != nextCurveOrientation {
					continue
				}
			}

			permittedPlacements = append(permittedPlacements, placement)
		}
	}

	return permittedPlacements, nil
}

func (e *Engine) CurrentPlayer() *Player {
	return e.Players[e.CurrentPlayerIndex]
}
