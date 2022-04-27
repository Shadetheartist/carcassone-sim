package main

import (
	"beeb/carcassonne/data"
	"fmt"
)

func main() {
	gameData := data.LoadGameData()
	fmt.Println(gameData)

	gameData.Explore()
}
