package engine

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine/board"
	"beeb/carcassonne/engine/deck"
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/engine/turnStage"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
	"errors"
	"fmt"
	"image/color"
	"sort"

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
	Deterministic                  bool
	BoardSize                      int
	GameOver                       bool
	GameBoard                      *board.Board
	GameData                       *data.GameData
	Players                        []*Player
	CurrentPlayerIndex             int
	TilePlacedThisTurn             *tile.Tile
	DecidedMeeplePlacementThisTurn *MeeplePlacement
	HeldRefTileGroup               *tile.ReferenceTileGroup
	CurrentPossibleTilePlacements  []Placement

	TurnCounter int
	TurnStage   turnStage.TurnStage

	RiverDeck *deck.Deck
	Deck      *deck.Deck

	TileFactory *tile.TileFactory

	TilePlacementManager *TilePlacementManager

	isFirstRiverTurn bool
	lastRiverTurn    int
	lastRiverTile    *tile.Tile
}

func NewEngine(gameData *data.GameData, boardSize int, numPlayers int) *Engine {
	if numPlayers > len(PLAYER_COLOR_LIST) {
		panic(fmt.Sprint("too many players for the colors implemented, max ", len(PLAYER_COLOR_LIST)))
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
	e.DecidedMeeplePlacementThisTurn = nil
	e.HeldRefTileGroup = nil
	e.CurrentPossibleTilePlacements = nil
	e.CurrentPlayerIndex = 0
	e.TurnStage = turnStage.Draw

	e.isFirstRiverTurn = true
	e.lastRiverTurn = 1
	e.lastRiverTile = nil
}

func (e *Engine) Step() {

	if e.GameOver {
		return
	}

	player := e.CurrentPlayer()

	switch e.TurnStage {
	case turnStage.Draw:

		e.TilePlacedThisTurn = nil
		e.DecidedMeeplePlacementThisTurn = nil
		e.HeldRefTileGroup = nil
		e.CurrentPossibleTilePlacements = nil

		//retry getting possible tiles a few times if we don't have a place to put one
		for i := 0; i < 3; i++ {

			rtg, tileTakeErr := e.TakeNextTile()

			if tileTakeErr != nil {
				//attempted to take a tile from an empty deck, so end the game
				e.EndGame()
				return
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

		if len(e.CurrentPossibleTilePlacements) < 1 {
			//"nowhere to place tile, tried 3 times, just remove tile completely" (take and do not place)
			_, _ = e.TakeNextTile()
			return
		}

		e.TurnStage++

	case turnStage.PlaceTile:

		selectedTilePlacement, meeplePlacement := player.DeterminePlacement(e, e.CurrentPossibleTilePlacements)
		e.DecidedMeeplePlacementThisTurn = meeplePlacement

		if selectedTilePlacement == nil {
			panic("No Placement Determined By Player")
		}

		e.TilePlacedThisTurn = e.PlaceTile(*selectedTilePlacement)

		e.CurrentPossibleTilePlacements = nil
		e.HeldRefTileGroup = nil
		e.TurnStage++

	case turnStage.PlaceMeeple:
		e.PlaceMeepleOnFeature()
		e.TurnStage++

	case turnStage.Score:
		e.ScoreFinishedFeature()
		e.TurnStage++

	case turnStage.Pass:
		err := e.GoToNextTurn()

		if err != nil {
			e.EndGame()
		}

		e.TurnStage = turnStage.Draw
	}
}

func (e *Engine) ScoreFinishedFeature() {
	mp := e.DecidedMeeplePlacementThisTurn

	if mp == nil {
		return
	}

	playerMeepleCountMap := make(map[*Player]int)

	for _, m := range mp.ReturnedMeeples {

		//fast-remove meeple from feature list
		for i, am := range m.Feature.AttachedMeeples {
			if am == m {
				l := len(m.Feature.AttachedMeeples) - 1
				m.Feature.AttachedMeeples[i] = m.Feature.AttachedMeeples[l]
				m.Feature.AttachedMeeples = m.Feature.AttachedMeeples[:l]
			}
		}

		//setting the meeples feature to nil returns it to the pool
		m.Feature = nil

		playerMeepleCountMap[m.ParentPlayer]++
	}

	mostPlayersOnFeature := 0
	for _, c := range playerMeepleCountMap {
		if c > mostPlayersOnFeature {
			mostPlayersOnFeature = c
		}
	}

	if mp.ScoreGained > 0 {
		for p, c := range playerMeepleCountMap {
			if c == mostPlayersOnFeature {
				p.Score += mp.ScoreGained
			}
		}
	}

}

// the feature selected by the player is only theoretical, the actual tile will have a different feature entirely
func (e *Engine) PlaceMeepleOnFeature() {

	t := e.TilePlacedThisTurn
	mp := e.DecidedMeeplePlacementThisTurn

	if t == nil {
		return
	}

	if mp == nil {
		return
	}

	if mp.SelectedMeeple == nil {
		return
	}

	var newTileFeature *tile.Feature

	for _, f := range t.Features {
		if f.ParentFeature == mp.ParentFeature {
			newTileFeature = f
		}
	}

	if newTileFeature == nil {
		panic("new tile does not have a corresponding feature to place a meeple on")
	}

	mp.SelectedMeeple.Feature = newTileFeature
	newTileFeature.AttachedMeeples = append(newTileFeature.AttachedMeeples, mp.SelectedMeeple)
}

func (e *Engine) PlaceTile(placement Placement) *tile.Tile {
	newTile := e.TileFactory.NewTileFromReference(placement.ReferenceTile)
	e.GameBoard.PlaceTile(placement.Position, newTile)

	if newTile.Reference.EdgeSignature.Contains(tile.River) {
		e.lastRiverTile = newTile

		if newTile.Reference.EdgeSignature.IsRiverCurving() {
			e.lastRiverTurn = newTile.Reference.Orientation
			e.isFirstRiverTurn = false
		}
	}

	return newTile
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

	for _, p := range e.Players {
		fmt.Println(p.Name, ": ", p.Score)
	}
}

func (e *Engine) restictRiverPlacement() ([]Placement, error) {

	permittedPlacements := make([]Placement, 0)
	permittedCurvedPlacements := make([]Placement, 0)

	for _, placement := range e.CurrentPossibleTilePlacements {

		connectedFeatures := placement.ConnectedFeatures

		for _, cf := range connectedFeatures {
			dir := cf.EdgeA

			connectedTilePos := placement.Position.EdgePos(dir)

			//must be connected to the last tile placed
			if connectedTilePos != e.lastRiverTile.Position {
				continue
			}

			//must be a river connection
			if cf.FeatureA.Type != tile.River {
				continue
			}

			//dont let the river turn the same way twice
			//if the new tile is curving, keep track of which way it was oriented
			//the next curved tile must be oriented 180 deg different from this tile
			if cf.FeatureA.ParentRefenceTileGroup.Orientations[0].EdgeSignature.IsRiverCurving() {

				//this is the first curve of the river, it can go either way
				if e.lastRiverTurn == 1 {
					permittedCurvedPlacements = append(permittedCurvedPlacements, placement)
					continue
				}

				// for the first turn, let it turn whatever way it wants
				if !e.isFirstRiverTurn {
					//next piece must be 180 degrees out of phase with the last
					nextCurveOrientation := (e.lastRiverTurn + 180) % 360

					if placement.ReferenceTile.Orientation != nextCurveOrientation {
						continue
					}
				}

				permittedCurvedPlacements = append(permittedCurvedPlacements, placement)
				continue
			}

			permittedPlacements = append(permittedPlacements, placement)
		}
	}

	// second pass to preferentially pick inward facing curves
	// only relevant when there are multiple valid curve placements
	if len(permittedCurvedPlacements) > 1 {

		magnitudes := make([]struct {
			p Placement
			m float64
		}, 0, len(permittedCurvedPlacements))

		for _, placement := range permittedCurvedPlacements {
			size := e.GameBoard.TileMatrix.Size()
			middle := util.Point[int]{X: size / 2, Y: size / 2}

			// add the direction the river is pointing to the tile's position
			// then calculate distance to center
			dir := computeRiverDirection(e, placement)
			pos := placement.Position.Add(dir)

			dist := middle.Subtract(pos)
			mag := dist.Magnitude()

			magnitudes = append(magnitudes, struct {
				p Placement
				m float64
			}{
				p: placement,
				m: mag,
			})
		}

		sort.Slice(magnitudes, func(i, j int) bool {
			return magnitudes[i].m < magnitudes[j].m
		})

		permittedPlacements = append(permittedPlacements, magnitudes[0].p)
	}

	return permittedPlacements, nil
}

// computeRiverDirection
// a river tiles pointed direction can be determined by looking at the open-end of the river on the tile
func computeRiverDirection(eng *Engine, placement Placement) util.Point[int] {

	for edge, feature := range placement.ReferenceTile.EdgeFeatures {
		if feature.Type != tile.River {
			continue
		}

		// river feature

		edgeDir := directions.Direction(edge)
		neighborTilePos := placement.Position.EdgePos(edgeDir)
		neighborTile, err := eng.GameBoard.TileMatrix.GetPt(neighborTilePos)
		if err != nil {
			continue
		}

		if neighborTile != nil {
			continue
		}

		// river feature which is open

		diff := neighborTilePos.Subtract(placement.Position)

		return diff
	}

	return util.Point[int]{}
}

func (e *Engine) CurrentPlayer() *Player {
	return e.Players[e.CurrentPlayerIndex]
}
