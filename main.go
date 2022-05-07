package main

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/simulator"
)

func main() {
	runSimulator()
	//runExplorer()
}

func runSimulator() {
	gameData := data.LoadGameData("./data/bitmaps")

	sim := simulator.NewSimulator(gameData, 32, 4)
	sim.Simulate()
}

func runExplorer() {
	gameData := data.LoadGameData("./data/bitmaps")

	gameData.Explore()
}
