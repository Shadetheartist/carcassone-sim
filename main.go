package main

import (
	"beeb/carcassonne/aiLink"
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"beeb/carcassonne/simulator"
	"math/rand"
)

func main() {
	runAILink()
	runSimulator()
	// runExplorer()
}

func runAILink() {
	rand.Seed(1)
	gameData := data.LoadGameData("./data/bitmaps", "./data/custom_deck.yml")
	engineInstance := engine.NewEngine(gameData, 16, 4)
	ai := aiLink.NewAILink(engineInstance)
	steps := (engineInstance.RiverDeck.Remaining() + engineInstance.Deck.Remaining()) * 5
	for i := 0; i < steps; i++ {
		engineInstance.Step()
	}

	ai.Inputs()
}

func runSimulator() {
	rand.Seed(1)
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")
	engineInstance := engine.NewEngine(gameData, 16, 4)
	sim := simulator.NewSimulator(engineInstance)
	sim.Simulate()
}

func runExplorer() {
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")

	gameData.Explore()
}
