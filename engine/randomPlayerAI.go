package engine

import "math/rand"

// RandomPlayerAI just literally places pieces randomly
type RandomPlayerAI struct{}

func (p *RandomPlayerAI) DeterminePlacement(_ *Engine, placementOptions []Placement) (*Placement, *MeeplePlacement) {
	if len(placementOptions) == 0 {
		return nil, nil
	}

	r := rand.Intn(len(placementOptions))
	return &placementOptions[r], nil
}
