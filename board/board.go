package board

import (
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/tile"
)

type Board struct {
	TileMatrix *matrix.Matrix[*tile.Tile]
}

func NewBoard(size int) *Board {
	board := &Board{}

	board.TileMatrix = matrix.NewMatrix[*tile.Tile](size)

	return board
}
