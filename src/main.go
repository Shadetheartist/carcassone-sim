package main

import (
	"beeb/carcassonne/game"

	"github.com/hajimehoshi/ebiten/v2"
)

var size float64 = 100

func main() {
	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")

	game := game.CreateGame()

	if err := ebiten.RunGame(&game); err != nil {
		panic(err)
	}

}
