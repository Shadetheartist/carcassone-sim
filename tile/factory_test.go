package tile_test

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/tile"
	"fmt"
	"testing"
)

func TestTileFactory_NewTileFromReference(t *testing.T) {
	gameData := data.LoadGameData("../data/bitmaps")

	referenceTile := gameData.ReferenceTiles["CloisterRiverRoad"]
	orientedReferenceTile := referenceTile[0]

	tf := &tile.TileFactory{}

	tl := tf.NewTileFromReference(orientedReferenceTile)

	fmt.Println(tl)
}

func BenchmarkTileFactory_NewTileFromReference(b *testing.B) {
	gameData := data.LoadGameData("../data/bitmaps")

	referenceTile := gameData.ReferenceTiles["CloisterRiverRoad"]
	orientedReferenceTile := referenceTile[0]

	tf := &tile.TileFactory{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		tf.NewTileFromReference(orientedReferenceTile)
	}
}
