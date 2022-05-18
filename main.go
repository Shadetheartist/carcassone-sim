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
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")

	rand.Seed(0)
	engine := engine.NewEngine(gameData, 45, 4)
	sim := simulator.NewSimulator(engine)
	sim.Simulate()
}

func runExplorer() {
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")

	gameData.Explore()
}
