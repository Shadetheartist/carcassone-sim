package tile_test

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/tile"
	"image"
	"testing"
)

func setupFarmTestScenario(tileFactory *tile.Factory, b *board.Board) {

	_tile1 := tileFactory.BuildTile("CastleRoadStraight")
	b.AddTile(_tile1, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	_tile2 := tileFactory.BuildTile("RoadTerminal3")
	b.AddTile(_tile2, tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 0,
		},
		Orientation: 0,
	})

	_tile3 := tileFactory.BuildTile("RoadTerminal3")
	b.AddTile(_tile3, tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 1,
		},
		Orientation: 180,
	})

	_tile4 := tileFactory.BuildTile("CastleCornerRoadCurve")
	b.AddTile(_tile4, tile.Placement{
		Position: tile.Position{
			X: -2,
			Y: 1,
		},
		Orientation: 180,
	})

	_tile5 := tileFactory.BuildTile("CastleFill3Road")
	b.AddTile(_tile5, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 90,
	})

	_tile6 := tileFactory.BuildTile("CastleCornerRoadCurve")
	b.AddTile(_tile6, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 0,
	})

}

func TestFarmContinuity(t *testing.T) {
	tileFactory := buildTileFactory()
	b := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	setupFarmTestScenario(tileFactory, &b)

	b.FarmSegmentAtPix(image.Point{0, 1})
}

func TestMatrixTransposition(t *testing.T) {
	tileFactory := buildTileFactory()

	_tile := tileFactory.BuildTile("CastleCorner")

	matrix := tile.OrientedFarmMatrix(_tile, 90)
	if matrix[0][5] == nil || matrix[0][6] != nil {
		t.Error("90 Degree Sus")
	}

	matrix = tile.OrientedFarmMatrix(_tile, 180)
	if matrix[5][6] == nil || matrix[6][6] != nil {
		t.Error("180 Degree Sus")
	}

	matrix = tile.OrientedFarmMatrix(_tile, 270)
	if matrix[6][1] == nil || matrix[6][0] != nil {
		t.Error("270 Degree Sus")
	}
}

func TestCloisterRiverRoadSegment(t *testing.T) {
	tileFactory := buildTileFactory()

	_tile := tileFactory.BuildTile("CloisterRiverRoad")
	tile.ComputeFarmMatrix(_tile)

}

func BenchmarkFarmSegment(b *testing.B) {

	tileFactory := buildTileFactory()

	_tile := tileFactory.BuildTile("CloisterRiverRoad")

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		tile.ComputeFarmMatrix(_tile)
	}
}
