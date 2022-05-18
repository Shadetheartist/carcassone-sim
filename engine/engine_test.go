package engine_test

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"math/rand"
	"testing"
)

func BenchmarkEngine(b *testing.B) {

	gameData := data.LoadGameData("../data/bitmaps", "../data/standard_deck.yml")

	rand.Seed(0)
	e1 := engine.NewEngine(gameData, 32, 4, true)

	steps := (e1.RiverDeck.Remaining() + e1.Deck.Remaining()) * 5

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		e1.InitGame()
		for i := 0; i < steps; i++ {
			e1.Step()
		}
	}
}
