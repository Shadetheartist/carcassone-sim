package main

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"beeb/carcassonne/simulator"
	"math/rand"
)

func main() {
	runSimulator()
	//runExplorer()
}

func runSimulator() {
	rand.Seed(11)
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")
	engineInstance := engine.NewEngine(gameData, 16, 4)
	sim := simulator.NewSimulator(engineInstance)
	sim.Simulate()
}

func runExplorer() {
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")

	gameData.Explore()
}
