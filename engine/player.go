package engine

import (
	"github.com/google/uuid"
	"image/color"
)

type Player struct {
	Id      uuid.UUID
	Name    string
	Color   color.Color
	Score   int
	Meeples []*Meeple

	AI *PlayerAI
}

var MaxMeeples int = 7

func NewPlayer(name string, color color.Color) *Player {
	player := &Player{}

	player.Id = uuid.New()
	player.Name = name
	player.Meeples = make([]*Meeple, MaxMeeples)
	player.Color = color

	for i := 0; i < MaxMeeples; i++ {
		player.Meeples[i] = &Meeple{
			Id:           uuid.New(),
			ParentPlayer: player,
			Power:        1,
			Feature:      nil,
		}
	}

	player.AI = &PlayerAI{
		Player:     player,
		Evaluation: Evaluation{},
	}

	return player
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

func (p *Player) DeterminePlacement(e *Engine, placementOptions []Placement) (*Placement, *MeeplePlacement) {
	return p.AI.DeterminePlacement(e, placementOptions)
}
