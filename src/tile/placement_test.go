package tile

import (
	"beeb/carcassonne/directions"
	"fmt"
	"math/rand"
	"testing"
)

func testDir(t *testing.T, orientation uint16, inputDir directions.Direction, expectedDir directions.Direction) {
	pl := Placement{
		Position:    Position{},
		Orientation: orientation,
	}

	tileDir := pl.TileDirection(inputDir)
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

func TestTileDirection(t *testing.T) {

	testDir(t, 0, directions.North, directions.North)
	testDir(t, 0, directions.East, directions.East)
	testDir(t, 0, directions.South, directions.South)
	testDir(t, 0, directions.West, directions.West)

	testDir(t, 90, directions.North, directions.East)
	testDir(t, 90, directions.East, directions.South)
	testDir(t, 90, directions.South, directions.West)
	testDir(t, 90, directions.West, directions.North)

	testDir(t, 180, directions.North, directions.South)
	testDir(t, 180, directions.East, directions.West)
	testDir(t, 180, directions.South, directions.North)
	testDir(t, 180, directions.West, directions.East)

	testDir(t, 270, directions.North, directions.West)
	testDir(t, 270, directions.East, directions.North)
	testDir(t, 270, directions.South, directions.East)
	testDir(t, 270, directions.West, directions.South)
}

func benchmarkTileDir(placement Placement, dir directions.Direction, b *testing.B) {
	placement.TileDirection(dir)
}

func BenchmarkTileDirection(b *testing.B) {
	for n := 0; n < b.N; n++ {
		randOri := uint16(rand.Intn(3) * 90)
		pl := Placement{
			Position:    Position{},
			Orientation: randOri,
		}

		randDir := directions.Direction(rand.Intn(3))

		benchmarkTileDir(pl, randDir, b)
	}
}
