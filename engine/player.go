package engine

import (
	"image/color"
	"math/rand"

	"github.com/google/uuid"
)

type Player struct {
	Id      uuid.UUID
	Name    string
	Color   color.Color
	Score   int
	Meeples []*Meeple
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

	//for _, pl := range options {
	//	p.EvaluatePlacement(pl, e)
	//}

	randN := rand.Int() % len(options)
	return &options[randN]
}

type Evaluation struct {
	Score  int
	Meeple int
}

func (p *Player) EvaluatePlacement(placement Placement, e *Engine) Evaluation {

	return Evaluation{}
}
