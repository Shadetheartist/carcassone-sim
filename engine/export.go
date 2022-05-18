package engine

import (
	"beeb/carcassonne/tile"
	"beeb/carcassonne/util"
	"encoding/json"
	"image/color"
	"os"
)

type StateFeature struct {
	Id   string
	Type string
}

type StateTile struct {
	Id           string
	Name         string
	Position     util.Point[int]
	Features     []StateFeature
	EdgeFeatures []string
}

type StateMeeple struct {
	Id    string
	Power int
}

type StatePlayer struct {
	Id      string
	Name    string
	Color   color.Color
	Score   int
	Meeples []string
}

type StateDeck struct {
}

type EngineState struct {
	Players map[string]StatePlayer
	Meeples map[string]StateMeeple
	Tiles   map[string]StateTile
}

func NewEngineState(e *Engine) *EngineState {
	state := &EngineState{}

	state.Players = make(map[string]StatePlayer)
	state.Meeples = make(map[string]StateMeeple)

	for _, p := range e.Players {
		id := p.Id.String()

		player := StatePlayer{
			Id:    id,
			Name:  p.Name,
			Color: p.Color,
			Score: p.Score,
		}

		player.Meeples = make([]string, len(p.Meeples))

		for i, m := range p.Meeples {
			id := m.Id.String()
			state.Meeples[id] = StateMeeple{
				Id:    id,
				Power: m.Power,
			}

			player.Meeples[i] = id
		}

		state.Players[id] = player
	}

	state.Tiles = make(map[string]StateTile)

	e.GameBoard.TileMatrix.Iterate(func(t *tile.Tile, x int, y int, idx int) {

		if t == nil {
			return
		}

		id := t.Id.String()
		st := StateTile{
			Id:       id,
			Name:     t.Reference.Name,
			Position: t.Position,
		}

		st.Features = make([]StateFeature, len(t.Features))

		for i, f := range t.Features {
			st.Features[i] = StateFeature{
				Id:   f.Id.String(),
				Type: f.Type.String(),
			}
		}

		st.EdgeFeatures = make([]string, len(t.EdgeFeatures))

		for i, f := range t.EdgeFeatures {
			st.EdgeFeatures[i] = f.Id.String()
		}

		state.Tiles[st.Id] = st
	})

	return state
}

func (e *Engine) ExportEngineState() {

	state := NewEngineState(e)

	jsonBytes, err := json.MarshalIndent(state, "", "\t")

	if err != nil {
		panic(err)
	}

	os.WriteFile("./state.json", jsonBytes, 0644)
}
