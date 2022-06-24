package engine

import (
	"beeb/carcassonne/engine/tile"
	"image/color"
	"math/rand"

	"github.com/google/uuid"
)

type Player struct {
	Id         uuid.UUID
	Name       string
	Color      color.Color
	Score      int
	Meeples    []*Meeple
	Evaluation Evaluation
}

func NewPlayer(name string, color color.Color) *Player {
	player := &Player{}

	meepleCount := 7

	player.Id = uuid.New()
	player.Name = name
	player.Meeples = make([]*Meeple, meepleCount)
	player.Color = color

	for i := 0; i < meepleCount; i++ {
		player.Meeples[i] = &Meeple{
			Id:           uuid.New(),
			ParentPlayer: player,
			Power:        1,
			Feature:      nil,
		}
	}

	return player
}

func (p *Player) DeterminePlacement(options []Placement, e *Engine) *Placement {
	if len(options) == 0 {
		return nil
	}

	var bestFeatureEval FeatureEvaluation
	var bestPlacement *Placement
	for i, pl := range options {
		eval := p.EvaluatePlacement(pl, e)
		for _, fEval := range eval.EvaluatedFeatures {
			if fEval.Score > bestFeatureEval.Score {
				bestFeatureEval = fEval
				bestPlacement = &options[i]
			}
		}
	}

	if bestPlacement == nil {
		randN := rand.Int() % len(options)
		return &options[randN]
	}

	return bestPlacement
}

type FeatureEvaluation struct {
	Feature *tile.Feature
	Score   int
	Meeple  int
}

type Evaluation struct {
	EvaluatedFeatures map[*tile.Feature]FeatureEvaluation
}

func (p *Player) EvaluatePlacement(placement Placement, e *Engine) Evaluation {

	eval := Evaluation{}
	eval.EvaluatedFeatures = make(map[*tile.Feature]FeatureEvaluation)

	t := e.TileFactory.NewTileFromReference(placement.ReferenceTile)
	e.GameBoard.PlaceTile(placement.Position, t)

	visitedFeaturesOfTile := make(map[*tile.Feature]struct{})

	for _, f := range t.Features {
		//avoid re-evaluating features later
		if _, exists := visitedFeaturesOfTile[f]; exists {
			continue
		}

		chain := VisitFeatureLinks(f)

		//support to avoid re-evaluating features later
		for f := range chain.FeaturesVisited {
			if t.HasFeature(f) {
				visitedFeaturesOfTile[f] = struct{}{}
			}
		}

		featureEval := FeatureEvaluation{}

		//check if there is a meeple on the feature aleady

		//calculate potential score if no meeple

		chainLenTiles := len(chain.TilesVisited)
		featureEval.Feature = f
		featureEval.Score = f.Type.Score() * chainLenTiles

		//calculate score if meeple placed

		//calculate score if x meeples placed? Where X is remaining meeples

		//determine if the feature becomes complete (you get meeple back intantly)

		//estimate chance of meeple returning before the game ends

		if featureEval.Score > 0 {
			eval.EvaluatedFeatures[f] = featureEval
		}
	}

	e.GameBoard.RemoveTileAt(placement.Position)

	return eval
}

func VisitFeatureLinks(feature *tile.Feature) FeatureChain {

	fc := FeatureChain{}
	fc.FeaturesVisited = make(map[*tile.Feature]struct{})
	fc.TilesVisited = make(map[*tile.Tile]struct{})

	stack := make([]*tile.Feature, 0, 16)
	stack = append(stack, feature)

	var f *tile.Feature
	for len(stack) > 0 {
		f, stack = stack[len(stack)-1], stack[:len(stack)-1]

		for l := range f.Links {
			if _, exists := fc.FeaturesVisited[l]; !exists {
				stack = append(stack, l)
			}
		}

		fc.FeaturesVisited[f] = struct{}{}
		fc.TilesVisited[f.ParentTile] = struct{}{}
	}

	return fc
}

type FeatureChain struct {
	TilesVisited    map[*tile.Tile]struct{}
	FeaturesVisited map[*tile.Feature]struct{}
}
