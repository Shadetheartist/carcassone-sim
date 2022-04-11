package tile_test

import (
	"beeb/carcassonne/tile"
	"image"
	"testing"
)

func TestComputeFarmSegments(t *testing.T) {
	tiles := loadTiles()

	crrTile := tiles["CloisterRiverRoad"]

	tile.ComputeFarmSegments(crrTile)
}

func TestEdgePositions(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 7, 7))

	pos := tile.EdgePositions(img)

	if len(pos) != 24 {
		t.Error("Should be 24 edge pixels")
	}
}
