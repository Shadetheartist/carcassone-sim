package engine_test

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"testing"
)

func BenchmarkPossibleTilePlacements(b *testing.B) {
	gameData := data.LoadGameData("../data/bitmaps", "../data/mega_deck.yml")
	engine := engine.NewEngine(gameData, 64, 4)

	for i := 0; i < 32*5+1; i++ {
		engine.Step()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.TilePlacementManager.PossibleTilePlacements(engine.HeldRefTileGroup)
	}
}
