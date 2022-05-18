package engine

import (
	"image/color"

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
