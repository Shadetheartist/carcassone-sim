package board

import (
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/tile"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
)

type Board struct {
	TileMatrix      *matrix.Matrix[*tile.Tile]
	PlacedTileCount int
	OpenPositions   map[util.Point[int]]*tile.EdgeSignature
}

func NewBoard(size int) *Board {
	board := &Board{}

	board.TileMatrix = matrix.NewMatrix[*tile.Tile](size)
	board.OpenPositions = make(map[util.Point[int]]*tile.EdgeSignature, 8)

	return board
}

func (b *Board) PlaceTile(pos util.Point[int], t *tile.Tile) {
	b.TileMatrix.Set(pos.X, pos.Y, t)
	b.PlacedTileCount++

	t.Position = pos

	//setup neighbours
	if tl, err := b.TileMatrix.GetPt(pos.North()); err == nil {
		t.Neighbours.SetNorth(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.East()); err == nil {
		t.Neighbours.SetEast(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.South()); err == nil {
		t.Neighbours.SetSouth(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.West()); err == nil {
		t.Neighbours.SetWest(tl)
	}

	for d, n := range t.Neighbours {
		//if neighbour is not there add it as an open position
		if n == nil {
			edgePos := pos.EdgePos(directions.Direction(d))

			if _, err := b.TileMatrix.GetPt(edgePos); err != nil {
				b.OpenPositions[edgePos] = nil
				continue
			}

			b.OpenPositions[edgePos] = b.createOpenPositonSignature(edgePos)
		}
	}

	delete(b.OpenPositions, pos)
}

func (b *Board) createOpenPositonSignature(pos util.Point[int]) *tile.EdgeSignature {
	sig := &tile.EdgeSignature{}

	for i := 0; i < 4; i++ {
		dir := directions.Direction(i)
		complimentDirection := directions.Compliment[dir]
		t, err := b.TileMatrix.GetPt(pos.EdgePos(dir))

		if err != nil || t == nil {
			continue
		}

		sig[dir] = t.Reference.EdgeSignature[complimentDirection]
	}

	return sig
}
