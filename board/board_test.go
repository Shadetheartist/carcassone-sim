package board_test

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/db"
	"beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"math/rand"
	"sort"
	"testing"
	"time"
)

func buildTileFactory() *tile.Factory {

	tileInfoLoader := &db.ConfigFileDataLoader{}
	tileInfoLoader.LoadData("../data/tiles.yml")

	bitmapLoader := &db.DirectoryBitmapLoader{}
	bitmapLoader.LoadBitmapsFromDirectory("../data/bitmaps")

	tf := tile.CreateTileFactory(tileInfoLoader.GetAllTileNames(), tileInfoLoader, bitmapLoader)

	return tf
}

func setupBoard(tf *tile.Factory, placedTiles int) board.Board {
	board := board.CreateBoard(tf.ReferenceTiles(), 1000, 1000)

	rand.Seed(0)

	board.AddTile(selectRandomTile(tf), tile.Placement{
		Position: tile.Position{},
	})

	for i := 0; i < placedTiles; i++ {
		//keep looping though tiles until a possible placement for a tile is found
		for {
			t := selectRandomTile(tf)

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

func selectRandomTile(tileFactory *tile.Factory) *tile.Tile {

	tiles := tileFactory.ReferenceTiles()
	keys := make([]string, 0, len(tiles))
	for k := range tiles {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	n := rand.Intn(len(keys))

	randKey := keys[n]

	randTile := tiles[randKey]

	return randTile
}

func TestConnectedFeatures(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	tileA := tileFactory.BuildTile("RiverCurve")
	board.AddTile(tileA, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 90,
	})

	tileB := tileFactory.BuildTile("RiverCurve")
	tileBPlacement := tile.Placement{
		Position: tile.Position{
			X: -1,
			Y: 0,
		},
		Orientation: 270,
	}

	cf := board.ConnectedFeatures(tileB, tileBPlacement)

	if cf[directions.East] != tile.River {
		t.Error("East Feature should be river")
	}
}

func TestConnectedFeatures2(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	tileA := tileFactory.BuildTile("RiverCurve")
	board.AddTile(tileA, tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 0,
		},
		Orientation: 270,
	})

	tileB := tileFactory.BuildTile("RiverCurve")
	tileBPlacement := tile.Placement{
		Position: tile.Position{
			X: 0,
			Y: 1,
		},
		Orientation: 90,
	}

	cf := board.ConnectedFeatures(tileB, tileBPlacement)

	if cf[directions.North] != tile.River {
		t.Error("East Feature should be river")
	}
}

func TestAddTile(t *testing.T) {
	rand.Seed(time.Now().Unix())

	tileFactory := buildTileFactory()
	board := setupBoard(tileFactory, 32)

	_tile := selectRandomTile(tileFactory)

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
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileSouth := tileFactory.BuildTile("CastleEndCap")
	board.AddTile(tileSouth, tile.Placement{
		Position: plc.Position.South(),
	})

	op := board.OpenPositions[plc.Position]

	northFeatureType := op[int(directions.North)]
	if northFeatureType != tile.FeatureTypeRoad {
		t.Error("North Neighbour should be Road")
	}

	southFeatureType := op[int(directions.South)]
	if southFeatureType != tile.Castle {
		t.Error("South Neighbour should be Castle")
	}

}

func TestIsTilePlacable(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	testTile := tileFactory.BuildTile("RoadTerminal3")

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tileFactory.BuildTile("Cloister")
	board.AddTile(tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	orientation, err := board.IsTilePlaceable(testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}

	if orientation != 90 {
		t.Error("Wrong Orientation")
	}
}

func TestIsTilePlacable2(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	testTile := tileFactory.BuildTile("RoadTerminal3")

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(testTile, plc.Position)

	if err == nil {
		t.Error("Tile should not be placable, but is allowed")
	}
}

func TestIsTilePlacable3(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	testTile := tileFactory.BuildTile("RoadTerminal4")

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 90,
	})

	tileEast := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 0,
	})

	tileSouth := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileSouth, tile.Placement{
		Position:    plc.Position.South(),
		Orientation: 90,
	})

	tileWest := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}
}

func TestIsTilePlacable4(t *testing.T) {
	tileFactory := buildTileFactory()
	board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

	testTile := tileFactory.BuildTile("CastleFill3Road")

	plc := tile.Placement{
		Position:    tile.Position{},
		Orientation: 0,
	}

	tileNorth := tileFactory.BuildTile("CastleEndCap")
	board.AddTile(tileNorth, tile.Placement{
		Position:    plc.Position.North(),
		Orientation: 180,
	})

	tileEast := tileFactory.BuildTile("CastleEndCap")
	board.AddTile(tileEast, tile.Placement{
		Position:    plc.Position.East(),
		Orientation: 270,
	})

	tileWest := tileFactory.BuildTile("RoadStraight")
	board.AddTile(tileWest, tile.Placement{
		Position:    plc.Position.West(),
		Orientation: 0,
	})

	_, err := board.IsTilePlaceable(testTile, plc.Position)

	if err != nil {
		t.Error("Tile should be placable, but is not allowed")
	}
}

func TestPossibleTilePlacements(t *testing.T) {
	tileFactory := buildTileFactory()
	board := setupBoard(tileFactory, 512)

	testTile := tileFactory.BuildTile("CastleFill3Road")

	positions := board.PossibleTilePlacements(testTile)

	if len(positions) == 0 {
		t.Error("No positions possible to place tile (incredibly unlikely)")
	}

}

func BenchmarkIsTilePlaceable(b *testing.B) {
	rand.Seed(0)

	tileFactory := buildTileFactory()
	board := setupBoard(tileFactory, 32)

	_tile := selectRandomTile(tileFactory)
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

	tileFactory := buildTileFactory()
	board := setupBoard(tileFactory, boardComplexity)
	_tile := selectRandomTile(tileFactory)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		board.PossibleTilePlacements(_tile)
	}
}

func BenchmarkAddTiles(b *testing.B) {
	tileFactory := buildTileFactory()

	for n := 0; n < b.N; n++ {
		board := board.CreateBoard(tileFactory.ReferenceTiles(), 1000, 1000)

		_tile := selectRandomTile(tileFactory)

		pl := tile.Placement{
			Position:    tile.Position{},
			Orientation: 0,
		}

		b.StartTimer()
		board.AddTile(_tile, pl)
		b.StopTimer()
	}
}
