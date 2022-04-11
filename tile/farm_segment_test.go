package tile_test

import (
	"beeb/carcassonne/tile"
	"image"
	"testing"
)

func TestCloisterRiverRoadSegment(t *testing.T) {
	tiles := loadTiles()

	segments := tile.ComputeFarmSegments(tiles["CloisterRiverRoad"])

	if len(segments) != 3 {
		t.Error("Expected 3 segments, got ", len(segments))
	}

	if len(segments[0].EdgePixels) != 11 {
		t.Error("Segment 0 should have 11 edge pixels", len(segments))
	}

	if len(segments[1].EdgePixels) != 5 {
		t.Error("Segment 1 should have 5 edge pixels", len(segments))
	}

	if len(segments[2].EdgePixels) != 5 {
		t.Error("Segment 2 should have 5 edge pixels", len(segments))
	}
}

func TestRoadCurveSegment(t *testing.T) {
	tiles := loadTiles()

	segments := tile.ComputeFarmSegments(tiles["RoadCurve"])

	if len(segments) != 2 {
		t.Error("Expected 2 segments, got ", len(segments))
	}

	if len(segments[0].EdgePixels) != 17 {
		t.Error("Segment 0 should have 17 edge pixels", len(segments))
	}

	if len(segments[1].EdgePixels) != 5 {
		t.Error("Segment 1 should have 5 edge pixels", len(segments))
	}
}

func TestRiverRoadCurveSegment(t *testing.T) {
	tiles := loadTiles()

	segments := tile.ComputeFarmSegments(tiles["RiverRoadCurve"])

	if len(segments) != 3 {
		t.Error("Expected 3 segments, got ", len(segments))
	}

	if len(segments[0].EdgePixels) != 10 {
		t.Error("Segment 0 should have 10 edge pixels", len(segments))
	}

	if len(segments[1].EdgePixels) != 5 {
		t.Error("Segment 1 should have 5 edge pixels", len(segments))
	}

	if len(segments[2].EdgePixels) != 5 {
		t.Error("Segment 2 should have 5 edge pixels", len(segments))
	}
}

func TestCastleFill4ShieldSegment(t *testing.T) {
	tiles := loadTiles()

	segments := tile.ComputeFarmSegments(tiles["CastleFill4Shield"])

	if len(segments) != 0 {
		t.Error("Expected 0 segments, got ", len(segments))
	}
}

func TestCloisterSegment(t *testing.T) {
	tiles := loadTiles()

	segments := tile.ComputeFarmSegments(tiles["Cloister"])

	if len(segments) != 1 {
		t.Error("Expected 1 segments, got ", len(segments))
	}

	if len(segments[0].EdgePixels) != 24 {
		t.Error("Segment 0 should have 24 edge pixels", len(segments))
	}
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
		tile.ComputeFarmSegments(_tile)
	}
}
