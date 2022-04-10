package tile_test

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"fmt"
	"testing"
)

func setupBoard(tiles map[string]tile.Tile, placedTiles int) board.Board {
	board := board.CreateBoard(tiles, 1000, 1000)
	return board
}

func setupTestScenario(tiles map[string]tile.Tile, b *board.Board) {

	straightRoad := tiles["RoadStraight"]
	straightRoad.RoadSegments = straightRoad.ComputeRoadSegments()
	b.AddTile(&straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve0 := tiles["RoadCurve"]
	roadCurve0.RoadSegments = roadCurve0.ComputeRoadSegments()
	b.AddTile(&roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tiles["RoadCurve"]
	roadCurve1.RoadSegments = roadCurve1.ComputeRoadSegments()
	b.AddTile(&roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tiles["RoadStraight"]
	straightRoad2.RoadSegments = straightRoad2.ComputeRoadSegments()
	b.AddTile(&straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 0,
	})

	roadCurve2 := tiles["RoadCurve"]
	roadCurve2.RoadSegments = roadCurve2.ComputeRoadSegments()
	b.AddTile(&roadCurve2, tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 1,
		},
		Orientation: 180,
	})

	roadCurve3 := tiles["RoadCurve"]
	roadCurve3.RoadSegments = roadCurve3.ComputeRoadSegments()
	b.AddTile(&roadCurve3, tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 0,
		},
		Orientation: 270,
	})
}

func testRoad(
	t *testing.T,
	b board.Board,
	pos tile.Position,
	numUniqueSegments int,
	edgeConnectivity [4]bool,
) {

	roadTile := b.Tiles[pos]

	rslen := len(roadTile.UniqueRoadSegements())
	if rslen != numUniqueSegments {
		t.Error(fmt.Sprint("Expected ", numUniqueSegments, " Unique Road Segments, got ", rslen))
	} else {
		rs := roadTile.UniqueRoadSegements()[0]
		if rs.ParentTile != roadTile {
			t.Error("Wrong parent tile")
		}

		for _, d := range directions.List {
			if edgeConnectivity[d] && rs.EdgeSegments[d] == nil {
				t.Error(fmt.Sprint(directions.IntMap[d]), " Edge Should be connected")
			}

			if !edgeConnectivity[d] && rs.EdgeSegments[d] != nil {
				t.Error(fmt.Sprint(directions.IntMap[d]), " Edge Should NOT be connected")
			}
		}
	}
}

func TestRoadContinuity(t *testing.T) {
	tiles := loadTiles()
	b := board.CreateBoard(tiles, 1000, 1000)

	setupTestScenario(tiles, &b)

	testRoad(t, b, tile.Position{X: 0, Y: 0}, 1, [4]bool{false, true, false, true})
	testRoad(t, b, tile.Position{X: 1, Y: 0}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: 1, Y: 1}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: 0, Y: 1}, 1, [4]bool{false, true, false, true})
	testRoad(t, b, tile.Position{X: -1, Y: 1}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: -1, Y: 0}, 1, [4]bool{false, false, true, true})
}

func setupRoadTestOpen(tiles map[string]tile.Tile, b *board.Board) {

	straightRoad := tiles["RoadStraight"]
	straightRoad.RoadSegments = straightRoad.ComputeRoadSegments()
	b.AddTile(&straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve0 := tiles["RoadCurve"]
	roadCurve0.RoadSegments = roadCurve0.ComputeRoadSegments()
	b.AddTile(&roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tiles["RoadCurve"]
	roadCurve1.RoadSegments = roadCurve1.ComputeRoadSegments()
	b.AddTile(&roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tiles["RoadStraight"]
	straightRoad2.RoadSegments = straightRoad2.ComputeRoadSegments()
	b.AddTile(&straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 0,
	})

}

func setupRoadTestLoop(tiles map[string]tile.Tile, b *board.Board) {

	straightRoad := tiles["RoadCurve"]
	straightRoad.RoadSegments = straightRoad.ComputeRoadSegments()
	b.AddTile(&straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 270,
	})

	roadCurve0 := tiles["RoadCurve"]
	roadCurve0.RoadSegments = roadCurve0.ComputeRoadSegments()
	b.AddTile(&roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tiles["RoadCurve"]
	roadCurve1.RoadSegments = roadCurve1.ComputeRoadSegments()
	b.AddTile(&roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tiles["RoadCurve"]
	straightRoad2.RoadSegments = straightRoad2.ComputeRoadSegments()
	b.AddTile(&straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 180,
	})
}

func setupRoadTestTerminals(tiles map[string]tile.Tile, b *board.Board) {

	straightRoad := tiles["RoadTerminal4"]
	straightRoad.RoadSegments = straightRoad.ComputeRoadSegments()
	b.AddTile(&straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	terminal2 := tiles["RoadTerminal4"]
	terminal2.RoadSegments = terminal2.ComputeRoadSegments()
	b.AddTile(&terminal2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: -1,
		},
		Orientation: 0,
	})

	roadCurve0 := tiles["RoadCurve"]
	roadCurve0.RoadSegments = roadCurve0.ComputeRoadSegments()
	b.AddTile(&roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tiles["RoadTerminal4"]
	roadCurve1.RoadSegments = roadCurve1.ComputeRoadSegments()
	b.AddTile(&roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tiles["RoadCurve"]
	straightRoad2.RoadSegments = straightRoad2.ComputeRoadSegments()
	b.AddTile(&straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 180,
	})
}

func TestRoadOpen(t *testing.T) {
	tiles := loadTiles()
	b := board.CreateBoard(tiles, 1000, 1000)

	setupRoadTestOpen(tiles, &b)

	seg := b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[1]

	rd := board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 4 {
		t.Error(fmt.Sprint("Road should be 4 segments long, it was ", len(rd.Segments)))
	}

}

func TestRoadLoop(t *testing.T) {
	tiles := loadTiles()
	b := board.CreateBoard(tiles, 1000, 1000)

	setupRoadTestLoop(tiles, &b)

	seg := b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[2]

	rd := board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 4 {
		t.Error(fmt.Sprint("Road should be 4 segments long, it was ", len(rd.Segments)))
	}
}

func TestRoadTerminals(t *testing.T) {
	tiles := loadTiles()
	b := board.CreateBoard(tiles, 1000, 1000)

	setupRoadTestTerminals(tiles, &b)

	seg := b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[0]

	rd := board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 2 {
		t.Error(fmt.Sprint("Road should be 2 segments long, it was ", len(rd.Segments)))
	}

	seg = b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[1]

	rd = board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 3 {
		t.Error(fmt.Sprint("Road should be 3 segments long, it was ", len(rd.Segments)))
	}

	seg = b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[2]

	rd = board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 3 {
		t.Error(fmt.Sprint("Road should be 3 segments long, it was ", len(rd.Segments)))
	}

	seg = b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[3]

	rd = board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 1 {
		t.Error(fmt.Sprint("Road should be 1 segments long, it was ", len(rd.Segments)))
	}
}
