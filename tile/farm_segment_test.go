package tile_test

import (
	"beeb/carcassonne/tile"
	"image"
	"testing"
)

func TestMatrixTransposition(t *testing.T) {
	tiles := loadTiles()

	_tile := tiles["CastleCorner"]

	matrix := tile.OrientedFarmMatrix(&_tile, 90)
	if matrix[0][5] == nil || matrix[0][6] != nil {
		t.Error("90 Degree Sus")
	}

	matrix = tile.OrientedFarmMatrix(&_tile, 180)
	if matrix[5][6] == nil || matrix[6][6] != nil {
		t.Error("180 Degree Sus")
	}

	matrix = tile.OrientedFarmMatrix(&_tile, 270)
	if matrix[6][0] == nil || matrix[6][1] != nil {
		t.Error("270 Degree Sus")
	}
}

func TestCloisterRiverRoadSegment(t *testing.T) {
	tiles := loadTiles()

	tile.ComputeFarmMatrix(tiles["CloisterRiverRoad"])

}

func TestEdgePositions(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 7, 7))

	pos := tile.EdgePositions(img)

	if len(pos) != 24 {
		t.Error("Should be 24 edge pixels")
	}
}

func BenchmarkFarmSegment(b *testing.B) {

	tiles := loadTiles()

	_tile := tiles["CloisterRiverRoad"]

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tile.ComputeFarmMatrix(_tile)
	}
}
