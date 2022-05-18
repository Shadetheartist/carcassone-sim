package main

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"beeb/carcassonne/simulator"
)

func main() {
	runSimulator()
	//runExplorer()
}

func runSimulator() {
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")
	engine := engine.NewEngine(gameData, 32, 4)
	sim := simulator.NewSimulator(engine)
	sim.Simulate()
}

func runExplorer() {
	gameData := data.LoadGameData("./data/bitmaps", "./data/standard_deck.yml")

	gameData.Explore()
}
