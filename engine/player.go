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

func (p *Player) remainingMeeples() int {
	remainingMeeples := 0

	for _, m := range p.Meeples {
		if m.Feature == nil {
			remainingMeeples++
		}
	}

	return remainingMeeples
}

func (p *Player) GetAvailableMeepleWithPower(power int) *Meeple {
	for _, m := range p.Meeples {
		if m.Feature == nil && m.Power == power {
			return m
		}
	}

	return nil
}

func (p *Player) scoreMeepleCostEval(meepleCostEval MeepleCostEvaluation) float64 {
	//remainingMeeples := p.remainingMeeples()

	//low meeples increases spite play
	//low meeples increases desire get meeples back
	//low meeples deter placement when points not great

	baseScore := meepleCostEval.Score

	return float64(baseScore)
}

type MeeplePlacement struct {
	ParentFeature  *tile.Feature
	SelectedMeeple *Meeple
}

func (p *Player) DeterminePlacement(placementOptions []Placement, e *Engine) (*Placement, *MeeplePlacement) {
	if len(placementOptions) == 0 {
		return nil, nil
	}

	var bestScore float64 = 0
	var bestScoreMeepleCost int = 0
	var bestParentFeature *tile.Feature
	var bestPlacement *Placement

	for i, pl := range placementOptions {
		eval := p.EvaluatePlacement(pl, e)
		for _, featureEval := range eval.EvaluatedFeatures {
			for _, meepleCostEval := range featureEval.EvaluatedMeepleCosts {
				calculatedScore := p.scoreMeepleCostEval(meepleCostEval)
				if calculatedScore > bestScore {
					bestPlacement = &placementOptions[i]
					bestScoreMeepleCost = meepleCostEval.MeepleCost
					bestParentFeature = featureEval.Feature.ParentFeature
					bestScore = calculatedScore
				}
			}
		}
	}

	if bestPlacement == nil {
		randN := rand.Int() % len(placementOptions)
		return &placementOptions[randN], nil
	}

	selectedMeeple := e.CurrentPlayer().GetAvailableMeepleWithPower(bestScoreMeepleCost)

	if selectedMeeple == nil {
		return bestPlacement, nil
	}

	return bestPlacement, &MeeplePlacement{
		SelectedMeeple: selectedMeeple,
		ParentFeature:  bestParentFeature,
	}
}

type MeepleCostEvaluation struct {
	Score      int
	MeepleCost int
}

type FeatureEvaluation struct {
	Feature              *tile.Feature
	EvaluatedMeepleCosts []MeepleCostEvaluation
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

		chain := visitFeatureLinks(f)

		//support to avoid re-evaluating features later
		for f := range chain.FeaturesVisited {
			if t.HasFeature(f) {
				visitedFeaturesOfTile[f] = struct{}{}
			}
		}

		chainScore := calculateChainScore(f, chain)

		//quick exit if the score is irrelevent
		if chainScore < 1 {
			continue
		}

		featureEval := FeatureEvaluation{}
		featureEval.Feature = f
		featureEval.EvaluatedMeepleCosts = make([]MeepleCostEvaluation, 0, 4)

		meeplesOnChain := meeplesOnFeatureChain(chain)

		if len(meeplesOnChain) < 1 {
			featureEval.EvaluatedMeepleCosts = append(featureEval.EvaluatedMeepleCosts, MeepleCostEvaluation{
				MeepleCost: 1,
				Score:      chainScore,
			})
		} else {
			otherPlayersMeeples := meeplesPerOtherPlayer(meeplesOnChain, p)
			mostMeeplesByOtherPlayers := mostPlayerMeeples(otherPlayersMeeples)

			featureEval.EvaluatedMeepleCosts = append(featureEval.EvaluatedMeepleCosts, MeepleCostEvaluation{
				MeepleCost: mostMeeplesByOtherPlayers + 1,
				Score:      chainScore,
			})
		}

		eval.EvaluatedFeatures[f] = featureEval

		//calculate score if a meeple exists

		//calculate score if x meeples placed? Where X is remaining meeples

		//determine if the feature becomes complete (you get meeple back intantly)

		//estimate chance of meeple returning before the game ends

	}

	e.GameBoard.RemoveTileAt(placement.Position)

	return eval
}

type FeatureChain struct {
	TilesVisited    map[*tile.Tile]struct{}
	FeaturesVisited map[*tile.Feature]struct{}
}

func visitFeatureLinks(feature *tile.Feature) FeatureChain {

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

func meeplesOnFeatureChain(fc FeatureChain) []*Meeple {
	meeplesOnChain := make([]*Meeple, 0, 4)
	for f := range fc.FeaturesVisited {
		for _, mi := range f.AttachedMeeples {
			m := mi.(*Meeple)
			meeplesOnChain = append(meeplesOnChain, m)
		}
	}
	return meeplesOnChain
}

func mostPlayerMeeples(meeplesPerPlayer map[*Player][]*Meeple) int {
	most := 0

	for _, m := range meeplesPerPlayer {
		n := len(m)
		if n > most {
			most = n
		}
	}

	return most
}

func meeplesPerOtherPlayer(meeples []*Meeple, currentPlayer *Player) map[*Player][]*Meeple {
	meeplesPerPlayer := make(map[*Player][]*Meeple)

	for _, m := range meeples {
		//not the current player
		if m.ParentPlayer == currentPlayer {
			continue
		}

		if _, exists := meeplesPerPlayer[m.ParentPlayer]; exists {
			meeplesPerPlayer[m.ParentPlayer] = append(meeplesPerPlayer[m.ParentPlayer], m)
		} else {
			playerMeeples := make([]*Meeple, 0, 4)
			playerMeeples = append(playerMeeples, m)
			meeplesPerPlayer[m.ParentPlayer] = playerMeeples
		}
	}

	return meeplesPerPlayer
}

func calculateChainScore(f *tile.Feature, chain FeatureChain) int {
	chainLenTiles := len(chain.TilesVisited)

	switch f.Type {
	case tile.Road:
		return f.Type.Score() * chainLenTiles
	case tile.Castle:
		//add base castle value
		score := f.Type.Score() * chainLenTiles

		//add shields
		for vf := range chain.FeaturesVisited {
			if vf.Type == tile.Shield {
				score += f.Type.Score()
			}
		}

		return score
	}

	return 0
}
