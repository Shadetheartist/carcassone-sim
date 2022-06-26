package engine

import "beeb/carcassonne/engine/tile"

//a feature chain is an interlinked group of features,
//like a big castle, or long road, or expansive farm
type FeatureChain struct {
	Feature          *tile.Feature
	TilesVisited     map[*tile.Tile]struct{}
	FeaturesVisited  map[*tile.Feature]struct{}
	isComplete       bool
	meeples          []*Meeple
	playerMeeplesMap map[*Player][]*Meeple
	owners           []*Player
	score            int
}

func newFeatureChain(feature *tile.Feature) FeatureChain {
	featureChain := FeatureChain{}
	featureChain.Feature = feature
	featureChain.FeaturesVisited = make(map[*tile.Feature]struct{})
	featureChain.TilesVisited = make(map[*tile.Tile]struct{})

	featureChain.traverseFeatureLinks(feature)
	featureChain.isComplete = featureChain.determineCompleteness()

	return featureChain
}

//iteratively go through each linked feature and add them to the featurechain's maps
func (fc *FeatureChain) traverseFeatureLinks(feature *tile.Feature) {
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
}

// if all the tiles' edges that are part of this feature are connected
// to something, the feature must be complete, so for each edge the feature touches,
// there must be a corresponding link
func (fc *FeatureChain) determineCompleteness() bool {
	for f := range fc.FeaturesVisited {

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

func (featureChain *FeatureChain) computeScore() {
	f := featureChain.Feature
	chainLenTiles := len(featureChain.TilesVisited)

	switch f.Type {
	case tile.Road:
		featureChain.score = f.Type.Score() * chainLenTiles
		return
	case tile.Castle:
		//add base castle value
		score := f.Type.Score() * chainLenTiles

		//add shields
		for vf := range featureChain.FeaturesVisited {
			if vf.Type == tile.Shield {
				score += f.Type.Score()
			}
		}

		featureChain.score = score
		return
	}
}

func (featureChain *FeatureChain) computeMeeples() {
	featureChain.meeples = make([]*Meeple, 0, 4)
	for f := range featureChain.FeaturesVisited {
		for _, mi := range f.AttachedMeeples {
			m := mi.(*Meeple)
			featureChain.meeples = append(featureChain.meeples, m)
		}
	}
}

func (featureChain *FeatureChain) computePlayerMeeplesMap() {
	featureChain.playerMeeplesMap = make(map[*Player][]*Meeple)
	for _, m := range featureChain.meeples {
		if _, exists := featureChain.playerMeeplesMap[m.ParentPlayer]; exists {
			featureChain.playerMeeplesMap[m.ParentPlayer] = append(featureChain.playerMeeplesMap[m.ParentPlayer], m)
		} else {
			playerMeeples := make([]*Meeple, 0, 4)
			playerMeeples = append(playerMeeples, m)
			featureChain.playerMeeplesMap[m.ParentPlayer] = playerMeeples
		}
	}
}

func (featureChain *FeatureChain) computeOwners() {
	playerMeeplesMap := featureChain.playerMeeplesMap

	//get the most meeples on the feature, for any player
	mostMeeplesOnFeature := 0
	for _, meeples := range playerMeeplesMap {
		numMeeples := len(meeples)

		if numMeeples > mostMeeplesOnFeature {
			mostMeeplesOnFeature = numMeeples
		}
	}

	//the owners of the feature are the players with the most meeples
	featureChain.owners = make([]*Player, 2)
	for p, meeples := range playerMeeplesMap {
		numMeeples := len(meeples)

		if numMeeples == mostMeeplesOnFeature {
			featureChain.owners = append(featureChain.owners, p)
		}
	}
}

func (featureChain *FeatureChain) direct() int {
	if !featureChain.isComplete {
		return 0
	}

	return featureChain.score
}

func (featureChain *FeatureChain) potential() int {
	if featureChain.isComplete {
		return 0
	}

	return featureChain.score
}

func (featureChain *FeatureChain) distanceFromOwner(p *Player) int {
	//get the most meeples on the feature, for any player
	mostMeeplesOnFeature := 0
	for _, meeples := range featureChain.playerMeeplesMap {
		numMeeples := len(meeples)

		if numMeeples > mostMeeplesOnFeature {
			mostMeeplesOnFeature = numMeeples
		}
	}

	return mostMeeplesOnFeature - len(featureChain.playerMeeplesMap[p])
}

func (featureChain *FeatureChain) hasOwner() bool {
	return len(featureChain.meeples) > 0
}

func (featureChain *FeatureChain) isOwner(p *Player) bool {
	for _, owner := range featureChain.owners {
		if p == owner {
			return true
		}
	}

	return false
}
