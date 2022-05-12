package engine

import (
	"beeb/carcassonne/tile"
)

type Meeple struct {
	Power        int
	ParentPlayer *Player
	Feature      *tile.Feature
}
