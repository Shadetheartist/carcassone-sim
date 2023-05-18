package aiLink

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"testing"
)

func BenchmarkAILink_Inputs(b *testing.B) {
	gameData := data.LoadGameData("../data/bitmaps", "../data/standard_deck.yml")
	engineInstance := engine.NewEngine(gameData, 16, 4)
	ai := NewAILink(engineInstance)

	// steps to complete game
	steps := (engineInstance.RiverDeck.Remaining() + engineInstance.Deck.Remaining()) * 5
	for i := 0; i < steps; i++ {
		engineInstance.Step()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ai.Inputs()
	}

}
