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

type MeeplePlacement struct {
	ParentFeature   *tile.Feature
	SelectedMeeple  *Meeple
	ReturnedMeeples []*Meeple
	ScoreGained     int
}

func (p *Player) DeterminePlacement(placementOptions []Placement, e *Engine) (*Placement, *MeeplePlacement) {
	if len(placementOptions) == 0 {
		return nil, nil
	}

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

type MeepleCostEvaluation struct {
	PotentialScore  int
	DirectScore     int
	MeepleCost      int
	MeeplesReturned []*Meeple
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

		meeplesOnChain := meeplesOnFeatureChain(chain)
		myMeeplesOnChain := p.myMeeplesOnChain(meeplesOnChain)
		otherPlayersMeeples := p.meeplesPerOtherPlayer(meeplesOnChain)
		mostMeeplesByOtherPlayers := mostPlayerMeeples(otherPlayersMeeples)

		//if this player controls the most meeples then we have no reason to add another one
		if len(myMeeplesOnChain) > mostMeeplesByOtherPlayers {
			continue
		}

		featureEval := FeatureEvaluation{}
		featureEval.Feature = f
		featureEval.EvaluatedMeepleCosts = make([]MeepleCostEvaluation, 0, 4)

		featureIsCompleted := isFeatureChainComplete(chain)

		var meeplesReturned []*Meeple
		directScoreFactor := 0
		potentialScoreFactor := 1
		if featureIsCompleted {
			directScoreFactor = 1
			potentialScoreFactor = 0
			meeplesReturned = meeplesOnChain
		}

		featureEval.EvaluatedMeepleCosts = append(featureEval.EvaluatedMeepleCosts, MeepleCostEvaluation{
			MeepleCost:      mostMeeplesByOtherPlayers + 1,
			DirectScore:     directScoreFactor * chainScore,
			PotentialScore:  chainScore * potentialScoreFactor,
			MeeplesReturned: meeplesReturned,
		})

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

// if all the tiles' edges that are part of this feature are connected
// to something, the feature must be complete, so for each edge the feature touches,
// there must be a corresponding link
func isFeatureChainComplete(chain FeatureChain) bool {
	for f := range chain.FeaturesVisited {

		featureEdgeCount := 0
		for _, ef := range f.ParentTile.EdgeFeatures {
			if ef == f {
				featureEdgeCount++
			}
		}

		if len(f.Links) != featureEdgeCount {
			return false
		}
	}

	return true
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

func (p *Player) myMeeplesOnChain(chainMeeples []*Meeple) []*Meeple {
	playerMeeples := make([]*Meeple, 0, 2)

	for _, m := range chainMeeples {
		//the current player
		if m.ParentPlayer == p {
			playerMeeples = append(playerMeeples, m)
		}
	}

	return playerMeeples
}

func (p *Player) meeplesPerOtherPlayer(meeples []*Meeple) map[*Player][]*Meeple {
	meeplesPerPlayer := make(map[*Player][]*Meeple)

	for _, m := range meeples {
		//not the current player
		if m.ParentPlayer == p {
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
