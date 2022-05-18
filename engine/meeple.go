package engine

import (
	"beeb/carcassonne/engine/tile"

	"github.com/google/uuid"
)

type Meeple struct {
	Id           uuid.UUID
	Power        int
	ParentPlayer *Player
	Feature      *tile.Feature
}
