package tile_test

import (
	"beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"fmt"
	"testing"
)

func testGridToTileDir(t *testing.T, orientation uint16, inputDir directions.Direction, expectedDir directions.Direction) {
	pl := tile.Placement{
		Position:    tile.Position{},
		Orientation: orientation,
	}

	tileDir := pl.GridToTileDir(inputDir)
	if tileDir != expectedDir {
		t.Error(
			fmt.Sprint("Placement Tile Direction Mismapped.",
				" orientation: ", orientation,
				" input: ", directions.IntMap[inputDir], inputDir,
				" expected: ", directions.IntMap[expectedDir], expectedDir,
				" got: ", directions.IntMap[tileDir], tileDir,
			))
	}
}

func testTileToGridDir(t *testing.T, orientation uint16, inputDir directions.Direction, expectedDir directions.Direction) {
	pl := tile.Placement{
		Position:    tile.Position{},
		Orientation: orientation,
	}

	tileDir := pl.TileToGridDir(inputDir)
	if tileDir != expectedDir {
		t.Error(
			fmt.Sprint("Placement Tile Direction Mismapped.",
				" orientation: ", orientation,
				" input: ", directions.IntMap[inputDir], inputDir,
				" expected: ", directions.IntMap[expectedDir], expectedDir,
				" got: ", directions.IntMap[tileDir], tileDir,
			))
	}
}

func TestTileToGridDir(t *testing.T) {

	testTileToGridDir(t, 0, directions.North, directions.North)
	testTileToGridDir(t, 0, directions.East, directions.East)
	testTileToGridDir(t, 0, directions.South, directions.South)
	testTileToGridDir(t, 0, directions.West, directions.West)

	testTileToGridDir(t, 90, directions.North, directions.East)
	testTileToGridDir(t, 90, directions.East, directions.South)
	testTileToGridDir(t, 90, directions.South, directions.West)
	testTileToGridDir(t, 90, directions.West, directions.North)

	testTileToGridDir(t, 180, directions.North, directions.South)
	testTileToGridDir(t, 180, directions.East, directions.West)
	testTileToGridDir(t, 180, directions.South, directions.North)
	testTileToGridDir(t, 180, directions.West, directions.East)

	testTileToGridDir(t, 270, directions.North, directions.West)
	testTileToGridDir(t, 270, directions.East, directions.North)
	testTileToGridDir(t, 270, directions.South, directions.East)
	testTileToGridDir(t, 270, directions.West, directions.South)
}

func TestGridToTileDir(t *testing.T) {

	testGridToTileDir(t, 0, directions.North, directions.North)
	testGridToTileDir(t, 0, directions.East, directions.East)
	testGridToTileDir(t, 0, directions.South, directions.South)
	testGridToTileDir(t, 0, directions.West, directions.West)

	testGridToTileDir(t, 90, directions.North, directions.West)
	testGridToTileDir(t, 90, directions.East, directions.North)
	testGridToTileDir(t, 90, directions.South, directions.East)
	testGridToTileDir(t, 90, directions.West, directions.South)

	testGridToTileDir(t, 180, directions.North, directions.South)
	testGridToTileDir(t, 180, directions.East, directions.West)
	testGridToTileDir(t, 180, directions.South, directions.North)
	testGridToTileDir(t, 180, directions.West, directions.East)

	testGridToTileDir(t, 270, directions.North, directions.East)
	testGridToTileDir(t, 270, directions.East, directions.South)
	testGridToTileDir(t, 270, directions.South, directions.West)
	testGridToTileDir(t, 270, directions.West, directions.North)
}
