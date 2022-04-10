package tile_test

import (
	"beeb/carcassonne/db"
	"beeb/carcassonne/tile"
	"fmt"
	"testing"
)

func buildTileFactory() *tile.Factory {

	tileInfoLoader := db.ConfigFileDataLoader{}
	tileInfoLoader.LoadData("../data/tiles.yml")

	bitmapLoader := db.DirectoryBitmapLoader{}
	bitmapLoader.LoadBitmapsFromDirectory("../data/bitmaps")

	tileFactory := &tile.Factory{}
	tileFactory.Initialize(tileInfoLoader.GetAllTileNames(), &tileInfoLoader, &bitmapLoader)

	return tileFactory
}

func TestFactoyBuiltTileReferences(t *testing.T) {
	factory := buildTileFactory()

	tileA := factory.BuildTile("RiverStraight")
	tileB := factory.BuildTile("RiverStraight")
	tileC := factory.BuildTile("RiverStraight")

	tileA.Neighbours[0] = &tileC
	if tileA.Neighbours[0] == tileB.Neighbours[0] {
		t.Errorf("Tiles reference the same memory - Neighbours")
	}

	tileA.Features[0] = &tile.Feature{}
	if tileA.Features[0] == tileB.Features[0] {
		t.Errorf("Tiles reference the same memory - Features")
	}

	tileA.RoadSegments[0] = &tile.RoadSegment{}
	if tileA.RoadSegments[0] == tileB.RoadSegments[0] {
		t.Errorf("Tiles reference the same memory - RoadSegments")
	}

	if tileA.Neighbours[0] == tileB.Neighbours[0] {
		t.Errorf("Tiles reference the same memory - Neighbours")
	}

	fmt.Println(tileA, tileB)
}

func BenchmarkBuildTile(b *testing.B) {
	factory := buildTileFactory()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		factory.BuildTile("RiverStraight")
	}
}
