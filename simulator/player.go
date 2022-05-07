package simulator

import "image/color"

type Player struct {
	Name    string
	Color   color.Color
	Score   int
	Meeples []*Meeple
}

func NewPlayer(name string, color color.Color) *Player {
	player := &Player{}

	meepleCount := 7

	player.Name = name
	player.Meeples = make([]*Meeple, meepleCount)

	for i := 0; i < meepleCount; i++ {
		player.Meeples[i] = &Meeple{
			ParentPlayer: player,
			Power:        1,
			Feature:      nil,
		}
	}

	return player
}
