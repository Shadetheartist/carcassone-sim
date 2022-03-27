package main

import (
	"beeb/carcassonne/game"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var size float64 = 100

func main() {
	rand.Seed(time.Now().Unix())

	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")
	ebiten.SetScreenClearedEveryFrame(false)
	game := game.Game{}
	game.Initialize()

	if err := ebiten.RunGame(&game); err != nil {
		panic(err)
	}
}
