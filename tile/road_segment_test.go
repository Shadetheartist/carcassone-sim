package tile_test

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"fmt"
	"testing"
)

func setupBoard(tf *tile.Factory, placedTiles int) board.Board {
	board := board.CreateBoard(tf.ReferenceTiles(), 1000, 1000)
	return board
}

func setupTestScenario(tileFactory *tile.Factory, b *board.Board) {

	straightRoad := tileFactory.BuildTile("RoadStraight")
	straightRoad.RoadSegments = tile.ComputeRoadSegments(straightRoad)
	b.AddTile(straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve0 := tileFactory.BuildTile("RoadCurve")
	roadCurve0.RoadSegments = tile.ComputeRoadSegments(roadCurve0)
	b.AddTile(roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tileFactory.BuildTile("RoadCurve")
	roadCurve1.RoadSegments = tile.ComputeRoadSegments(roadCurve1)
	b.AddTile(roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tileFactory.BuildTile("RoadStraight")
	straightRoad2.RoadSegments = tile.ComputeRoadSegments(straightRoad2)
	b.AddTile(straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 0,
	})

	roadCurve2 := tileFactory.BuildTile("RoadCurve")
	roadCurve2.RoadSegments = tile.ComputeRoadSegments(roadCurve2)
	b.AddTile(roadCurve2, tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 1,
		},
		Orientation: 180,
	})

	roadCurve3 := tileFactory.BuildTile("RoadCurve")
	roadCurve3.RoadSegments = tile.ComputeRoadSegments(roadCurve3)
	b.AddTile(roadCurve3, tile.Placement{
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
	tileFactory := buildTileFactory()
	b := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	setupTestScenario(tileFactory, &b)

	testRoad(t, b, tile.Position{X: 0, Y: 0}, 1, [4]bool{false, true, false, true})
	testRoad(t, b, tile.Position{X: 1, Y: 0}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: 1, Y: 1}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: 0, Y: 1}, 1, [4]bool{false, true, false, true})
	testRoad(t, b, tile.Position{X: -1, Y: 1}, 1, [4]bool{false, false, true, true})
	testRoad(t, b, tile.Position{X: -1, Y: 0}, 1, [4]bool{false, false, true, true})
}

func setupRoadTestOpen(tileFactory *tile.Factory, b *board.Board) {

	straightRoad := tileFactory.BuildTile("RoadStraight")
	straightRoad.RoadSegments = tile.ComputeRoadSegments(straightRoad)
	b.AddTile(straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve0 := tileFactory.BuildTile("RoadCurve")
	roadCurve0.RoadSegments = tile.ComputeRoadSegments(roadCurve0)
	b.AddTile(roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tileFactory.BuildTile("RoadCurve")
	roadCurve1.RoadSegments = tile.ComputeRoadSegments(roadCurve1)
	b.AddTile(roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tileFactory.BuildTile("RoadStraight")
	straightRoad2.RoadSegments = tile.ComputeRoadSegments(straightRoad2)
	b.AddTile(straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 0,
	})

}

func setupRoadTestLoop(tileFactory *tile.Factory, b *board.Board) {

	straightRoad := tileFactory.BuildTile("RoadCurve")
	straightRoad.RoadSegments = tile.ComputeRoadSegments(straightRoad)
	b.AddTile(straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 270,
	})

	roadCurve0 := tileFactory.BuildTile("RoadCurve")
	roadCurve0.RoadSegments = tile.ComputeRoadSegments(roadCurve0)
	b.AddTile(roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tileFactory.BuildTile("RoadCurve")
	roadCurve1.RoadSegments = tile.ComputeRoadSegments(roadCurve1)
	b.AddTile(roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tileFactory.BuildTile("RoadCurve")
	straightRoad2.RoadSegments = tile.ComputeRoadSegments(roadCurve1)
	b.AddTile(straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 180,
	})
}

func setupRoadTestTerminals(tileFactory *tile.Factory, b *board.Board) {

	straightRoad := tileFactory.BuildTile("RoadTerminal4")
	straightRoad.RoadSegments = tile.ComputeRoadSegments(straightRoad)
	b.AddTile(straightRoad, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 0,
	})

	terminal2 := tileFactory.BuildTile("RoadTerminal4")
	terminal2.RoadSegments = tile.ComputeRoadSegments(terminal2)
	b.AddTile(terminal2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: -1,
		},
		Orientation: 0,
	})

	roadCurve0 := tileFactory.BuildTile("RoadCurve")
	roadCurve0.RoadSegments = tile.ComputeRoadSegments(roadCurve0)
	b.AddTile(roadCurve0, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 0,
		},
		Orientation: 0,
	})

	roadCurve1 := tileFactory.BuildTile("RoadTerminal4")
	roadCurve1.RoadSegments = tile.ComputeRoadSegments(roadCurve1)
	b.AddTile(roadCurve1, tile.Placement{
		Position: tile.Position{
			X: 1,
			Y: 1,
		},
		Orientation: 90,
	})

	straightRoad2 := tileFactory.BuildTile("RoadCurve")
	straightRoad2.RoadSegments = tile.ComputeRoadSegments(straightRoad2)
	b.AddTile(straightRoad2, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 180,
	})
}

func TestRoadOpen(t *testing.T) {
	tileFactory := buildTileFactory()
	b := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	setupRoadTestOpen(tileFactory, &b)

	seg := b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[1]

	rd := board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 4 {
		t.Error(fmt.Sprint("Road should be 4 segments long, it was ", len(rd.Segments)))
	}

}

func TestRoadLoop(t *testing.T) {
	tileFactory := buildTileFactory()
	b := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	setupRoadTestLoop(tileFactory, &b)

	seg := b.Tiles[tile.Position{X: 0, Y: 0}].RoadSegments[2]

	rd := board.CompileRoadFromSegment(seg)

	if len(rd.Segments) != 4 {
		t.Error(fmt.Sprint("Road should be 4 segments long, it was ", len(rd.Segments)))
	}
}

func TestRoadTerminals(t *testing.T) {
	tileFactory := buildTileFactory()
	b := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	setupRoadTestTerminals(tileFactory, &b)

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
