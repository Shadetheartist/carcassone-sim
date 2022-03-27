package board

import (
	"beeb/carcassonne/directions"
	"beeb/carcassonne/loader"
	"beeb/carcassonne/tile"
	"math/rand"
	"testing"
	"time"
)

func loadTiles() map[string]tile.Tile {
	tiles, _ := loader.LoadTiles("../data/tiles.yml", "../data/bitmaps")
	return tiles
}

func setupBoard(tiles map[string]tile.Tile, placedTiles int) Board {
	board := New()

	rand.Seed(0)

	board.AddTile(selectRandomTile(tiles), tile.Placement{
		Position: tile.Position{},
	})

	for i := 0; i < placedTiles; i++ {
		//keep looping though tiles until a possible placement for a tile is found
		for {
			t := selectRandomTile(tiles)

			placements := board.PossibleTilePlacements(t)

			if len(placements) < 1 {
				continue
			}

			randomPlacement := placements[rand.Intn(len(placements))]
			board.AddTile(t, randomPlacement)

			break
		}
	}

	return board
}

func selectRandomTile(tiles map[string]tile.Tile) *tile.Tile {
	keys := make([]string, 0, len(tiles))
	for k := range tiles {
		keys = append(keys, k)
	}

	n := rand.Intn(len(keys))

	randKey := keys[n]

	randTile := tiles[randKey]

	return &randTile
}

func TestConnectedFeatures(t *testing.T) {
	tiles := loadTiles()
	board := New()

	tileA := tiles["RiverCurve"]
	board.AddTile(&tileA, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 90,
	})

	tileB := tiles["RiverCurve"]
	tileBPlacement := tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 0,
		},
		Orientation: 270,
	}

	cf := board.ConnectedFeatures(&tileB, tileBPlacement)

	if cf[directions.East] != tile.River {
		t.Error("East Feature should be river")
	}
}

func TestConnectedFeatures2(t *testing.T) {
	tiles := loadTiles()
	board := New()

	tileA := tiles["RiverCurve"]
	board.AddTile(&tileA, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 270,
	})

	tileB := tiles["RiverCurve"]
	tileBPlacement := tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 90,
	}

	cf := board.ConnectedFeatures(&tileB, tileBPlacement)

	if cf[directions.North] != tile.River {
		t.Error("East Feature should be river")
	}
}

func TestAddTile(t *testing.T) {
	rand.Seed(time.Now().Unix())

	tiles := loadTiles()
	board := setupBoard(tiles, 32)

	_tile := selectRandomTile(tiles)

	pl := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	board.AddTile(_tile, pl)

	addedTile := board.Tiles[pl.Position]
	if addedTile == nil {
		t.Error("Tile was not added")
	}
}

func TestOpenPositionState(t *testing.T) {
	tiles := loadTiles()
	board := New()

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tiles["RoadStraight"]
	board.AddTile(&tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileSouth := tiles["CastleEndCap"]
	board.AddTile(&tileSouth, tile.Placement{
		Position: plc.Position.South(),
	})

	op := board.OpenPositions[plc.Position]

	northFeatureType := op[int(directions.North)]
	if northFeatureType != tile.Road {
		t.Error("North Neighbour should be Road")
	}

	southFeatureType := op[int(directions.South)]
	if southFeatureType != tile.Castle {
		t.Error("South Neighbour should be Castle")
	}

}

func TestIsTilePlacable(t *testing.T) {
	tiles := loadTiles()
	board := New()

	testTile := tiles["RoadTerminal3"]

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tiles["RoadStraight"]
	board.AddTile(&tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tiles["Cloister"]
	board.AddTile(&tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tiles["RoadStraight"]
	board.AddTile(&tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tiles["RoadStraight"]
	board.AddTile(&tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	orientation, err := board.IsTilePlaceable(&testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}

	if orientation != 90 {
		t.Error("Wrong Orientation")
	}
}

func TestIsTilePlacable2(t *testing.T) {
	tiles := loadTiles()
	board := New()

	testTile := tiles["RoadTerminal3"]

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tiles["RoadStraight"]
	board.AddTile(&tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tiles["RoadStraight"]
	board.AddTile(&tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tiles["RoadStraight"]
	board.AddTile(&tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tiles["RoadStraight"]
	board.AddTile(&tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(&testTile, plc.Position)

	if err == nil {
		t.Error("Tile should not be placable, but is allowed")
	}
}

func TestIsTilePlacable3(t *testing.T) {
	tiles := loadTiles()
	board := New()

	testTile := tiles["RoadTerminal4"]

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tiles["RoadStraight"]
	board.AddTile(&tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tiles["RoadStraight"]
	board.AddTile(&tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tiles["RoadStraight"]
	board.AddTile(&tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tiles["RoadStraight"]
	board.AddTile(&tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(&testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}
}

func TestIsTilePlacable4(t *testing.T) {
	tiles := loadTiles()
	board := New()

	testTile := tiles["CastleFill3Road"]

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tiles["CastleEndCap"]
	board.AddTile(&tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 180,
	})

	tileEast := tiles["CastleEndCap"]
	board.AddTile(&tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 270,
	})

	tileSouth := tiles["CastleEndCap"]
	board.AddTile(&tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 0,
	})

	tileWest := tiles["RoadStraight"]
	board.AddTile(&tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(&testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}
}

func BenchmarkIsTilePlaceable(b *testing.B) {
	rand.Seed(0)

	tiles := loadTiles()
	board := setupBoard(tiles, 32)

	_tile := selectRandomTile(tiles)
	keys := make([]tile.Position, 0, len(board.OpenPositions))
	for k := range board.OpenPositions {
		keys = append(keys, k)
	}
	randomPlacement := keys[rand.Intn(len(keys))]

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		board.IsTilePlaceable(_tile, randomPlacement)
	}
}

func BenchmarkPossiblePlacement32(b *testing.B) {
	benchmarkPossiblePlacement(32, b)
}

func BenchmarkPossiblePlacement128(b *testing.B) {
	benchmarkPossiblePlacement(128, b)
}

func BenchmarkPossiblePlacement512(b *testing.B) {
	benchmarkPossiblePlacement(512, b)
}

func BenchmarkPossiblePlacement2048(b *testing.B) {
	benchmarkPossiblePlacement(2048, b)
}

func benchmarkPossiblePlacement(boardComplexity int, b *testing.B) {
	rand.Seed(0)

	tiles := loadTiles()
	board := setupBoard(tiles, boardComplexity)
	_tile := selectRandomTile(tiles)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		board.PossibleTilePlacements(_tile)
	}
}

func BenchmarkAddTiles(b *testing.B) {
	tiles := loadTiles()

	for n := 0; n < b.N; n++ {
		board := New()

		_tile := selectRandomTile(tiles)

		pl := tile.Placement{
			Position:    tile.Position{},
			Orientation: 0,
		}

		b.StartTimer()
		board.AddTile(_tile, pl)
		b.StopTimer()
	}
}
