package board

import (
	"beeb/carcassonne/tile"
	"image"
	"math"
)

func (b *Board) FarmSegmentAtPix(pix image.Point) *tile.FarmSegment {
	_tile, offset := b.TileAtPix(pix)

	matrix := tile.OrientedFarmMatrix(_tile, int(_tile.Placement.Orientation))

	farmSegment := matrix[offset.Y][offset.X]

	return farmSegment
}

func (b *Board) pixToPos(n int) int {
	return int(math.Floor(float64(n) / b.RenderScale / float64(b.BaseSize)))
}

func (b *Board) PixToPos(pix image.Point) tile.Position {
	return tile.Position{
		X: b.pixToPos(pix.X),
		Y: b.pixToPos(pix.Y),
	}
}

func (b *Board) posToPix(n int) int {
	return n * b.BaseSize * int(b.RenderScale)
}

func (b *Board) PosToPix(pos tile.Position) image.Point {
	return image.Point{
		X: b.posToPix(pos.X),
		Y: b.posToPix(pos.Y),
	}
}

func (b *Board) TileAtPix(pix image.Point) (*tile.Tile, image.Point) {
	pos := b.PixToPos(pix)

	//reconverted is the top left corner of the tile
	topLeftPix := b.PosToPix(pos)

	offset := image.Point{
		pix.X - topLeftPix.X,
		pix.Y - topLeftPix.Y,
	}

	return b.Tiles[pos], offset
}
