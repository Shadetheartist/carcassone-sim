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

type Evaluation struct {
	EvaluatedFeatures map[*tile.Feature]FeatureEvaluation
}

type FeatureEvaluation struct {
	Feature              *tile.Feature
	EvaluatedMeepleCosts []MeepleCostEvaluation
}

type MeepleCostEvaluation struct {
	PotentialScore    int
	DirectScore       int
	MeepleCost        int
	MeeplesReturned   []*Meeple
	PlayerScoreChange map[*Player]int
}

type MeeplePlacement struct {
	ParentFeature   *tile.Feature
	SelectedMeeple  *Meeple
	ReturnedMeeples []*Meeple
	ScoreGained     int
}

var MAX_MEEPLES int = 7

func NewPlayer(name string, color color.Color) *Player {
	player := &Player{}

	player.Id = uuid.New()
	player.Name = name
	player.Meeples = make([]*Meeple, MAX_MEEPLES)
	player.Color = color

	for i := 0; i < MAX_MEEPLES; i++ {
		player.Meeples[i] = &Meeple{
			Id:           uuid.New(),
			ParentPlayer: player,
			Power:        1,
			Feature:      nil,
		}
	}

	return player
}

func (p *Player) scoreMeepleCostEval(meepleCostEval MeepleCostEvaluation, e *Engine) float32 {

	var playerRiskFactor float32 = 0.75

	var directScoreFactor float32 = 1
	var potentialScoreFactor float32 = 0.35

	numMeeplesRemaining := p.numRemainingMeeples()
	var meeplesRemainingFactor float32 = 2 - (float32(numMeeplesRemaining) / float32(MAX_MEEPLES))

	//see if we even have a meeple with enough power to do this
	selectedMeeple := e.CurrentPlayer().GetAvailableMeepleWithPower(meepleCostEval.MeepleCost)

	//we can't place a meeple, but maybe we can worsen someone else's position?
	//if we can't extend ours...
	if selectedMeeple == nil {
		//spiteFactor = 0.5
	}

	directScore := (directScoreFactor * float32(meepleCostEval.DirectScore))
	potentialScore := (playerRiskFactor * potentialScoreFactor * float32(meepleCostEval.PotentialScore))

	//for each meeple we dont have in our pool, we like the direct score a little more
	directScore *= meeplesRemainingFactor

	return directScore + potentialScore
}

func (p *Player) numRemainingMeeples() int {
	c := 0

	for _, m := range p.Meeples {
		if m.Feature == nil {
			c++
		}
	}

	return c
}

func (p *Player) GetAvailableMeepleWithPower(power int) *Meeple {
	for _, m := range p.Meeples {
		if m.Feature == nil && m.Power == power {
			return m
		}
	}

	return nil
}

func (p *Player) DeterminePlacement(placementOptions []Placement, e *Engine) (*Placement, *MeeplePlacement) {
	if len(placementOptions) == 0 {
		return nil, nil
	}

	// just place randomly
	// r := rand.Intn(len(placementOptions))
	// return &placementOptions[r], nil

	var bestScore float32 = 0
	var bestMeepleCostEval MeepleCostEvaluation
	var bestParentFeature *tile.Feature
	var bestPlacement *Placement

	for i, pl := range placementOptions {
		eval := p.EvaluatePlacement(pl, e)
		for _, featureEval := range eval.EvaluatedFeatures {
			for _, meepleCostEval := range featureEval.EvaluatedMeepleCosts {
				calculatedScore := p.scoreMeepleCostEval(meepleCostEval, e)
				if calculatedScore > bestScore {
					bestPlacement = &placementOptions[i]
					bestParentFeature = featureEval.Feature.ParentFeature
					bestScore = calculatedScore
					bestMeepleCostEval = meepleCostEval
				}
			}
		}
	}

	if bestPlacement == nil {
		randN := rand.Int() % len(placementOptions)
		return &placementOptions[randN], nil
	}

	selectedMeeple := e.CurrentPlayer().GetAvailableMeepleWithPower(bestMeepleCostEval.MeepleCost)

	if selectedMeeple != nil {

		//case for adding then removing the meeple on the same step
		if bestMeepleCostEval.DirectScore != 0 && len(bestMeepleCostEval.MeeplesReturned) < 1 {
			bestMeepleCostEval.MeeplesReturned = append(bestMeepleCostEval.MeeplesReturned, selectedMeeple)
		}
	}

	return bestPlacement, &MeeplePlacement{
		SelectedMeeple:  selectedMeeple,
		ParentFeature:   bestParentFeature,
		ReturnedMeeples: bestMeepleCostEval.MeeplesReturned,
		ScoreGained:     bestMeepleCostEval.DirectScore,
	}
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

		featureChain := newFeatureChain(f)

		//support to avoid re-evaluating features later
		for f := range featureChain.FeaturesVisited {
			if t.HasFeature(f) {
				visitedFeaturesOfTile[f] = struct{}{}
			}
		}

		featureChain.computeScore()

		//quick exit if the score is irrelevent
		if featureChain.score < 1 {
			continue
		}

		featureChain.computeMeeples()
		featureChain.computePlayerMeeplesMap()
		featureChain.computeOwners()

		//just don't add to features that are owned, but not by you
		if featureChain.hasOwner() && !featureChain.isOwner(p) {
			continue
		}

		featureEval := FeatureEvaluation{}
		featureEval.Feature = f
		featureEval.EvaluatedMeepleCosts = make([]MeepleCostEvaluation, 0, 4)

		var meeplesReturned []*Meeple
		if featureChain.isComplete {
			meeplesReturned = featureChain.meeples
		}

		meepleCost := 1

		if featureChain.isOwner(p) {
			meepleCost = 0
		}

		meepleCostEval := MeepleCostEvaluation{
			MeepleCost:        meepleCost,
			DirectScore:       featureChain.direct(),
			PotentialScore:    featureChain.potential(),
			MeeplesReturned:   meeplesReturned,
			PlayerScoreChange: make(map[*Player]int),
		}

		for _, p := range featureChain.owners {
			meepleCostEval.PlayerScoreChange[p] = featureChain.score
		}

		featureEval.EvaluatedMeepleCosts = append(featureEval.EvaluatedMeepleCosts, meepleCostEval)

		eval.EvaluatedFeatures[f] = featureEval

		//estimate chance of meeple returning before the game ends

	}

	e.GameBoard.RemoveTileAt(placement.Position)

	return eval
}
