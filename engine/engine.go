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

	TileFactory       *tile.TileFactory
	placementBuffer   []Placement
	orientationBuffer []*tile.ReferenceTile
	connectionsBuffer []Connection

	placementFunction PlacementFunction
}

func NewEngine(gameData *data.GameData, boardSize int, numPlayers int, deterministic bool) *Engine {
	if numPlayers > len(PLAYER_COLOR_LIST) {
		panic(fmt.Sprint("too many players, max ", len(PLAYER_COLOR_LIST)))
	}

	engine := &Engine{}

	engine.Deterministic = deterministic
	engine.BoardSize = boardSize
	engine.GameData = gameData
	engine.Players = make([]*Player, numPlayers)
	engine.TileFactory = &tile.TileFactory{}
	engine.placementBuffer = make([]Placement, 0, 128)
	engine.orientationBuffer = make([]*tile.ReferenceTile, 4)
	engine.connectionsBuffer = make([]Connection, 0, 4)

	if deterministic {
		engine.placementFunction = possibleTilePlacementsDeterministic
	} else {
		engine.placementFunction = possibleTilePlacementsNonDeterministic
	}

	engine.InitGame()

	return engine
}

func (e *Engine) InitGame() {
	for i := 0; i < len(e.Players); i++ {
		playerName := fmt.Sprint("Player ", i)
		e.Players[i] = NewPlayer(playerName, PLAYER_COLOR_LIST[i])
	}

	e.GameBoard = board.NewBoard(e.BoardSize, e.Deterministic)
	e.RiverDeck = deck.BuildRiverDeck(e.GameData)
	e.Deck = deck.BuildDeck(e.GameData)
	e.GameOver = false
	e.TurnCounter = 0
	e.TilePlacedThisTurn = nil
	e.HeldRefTileGroup = nil
	e.CurrentPossibleTilePlacements = nil
	e.CurrentPlayerIndex = 0
	e.TurnStage = turnStage.Draw
}

func (e *Engine) Step() {

	if e.GameOver {
		return
	}

	switch e.TurnStage {
	case turnStage.Draw:

		rtg, tileTakeErr := e.TakeNextTile()

		if tileTakeErr != nil {
			panic("Attempted to take a tile from an empty deck.")
		}

		e.HeldRefTileGroup = rtg
		e.CurrentPossibleTilePlacements = e.PossibleTilePlacements(rtg)

		if e.GameBoard.PlacedTileCount > 0 && (e.RiverDeck.Remaining() > 0 || rtg.IsRiverTile()) {
			e.CurrentPossibleTilePlacements, _ = e.restictRiverPlacement()
		}

		//here would be a 'shuffle tile into deck when not playable' clause

		e.TurnStage++

	case turnStage.PlaceTile:
		randomPlacement := RandomPlacement(e.CurrentPossibleTilePlacements)
		if randomPlacement == nil {
			return
		}

		e.PlaceTile(*randomPlacement)
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
			if cf.FeatureA.ParentRefenceTile.EdgeSignature.Curving() {

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
