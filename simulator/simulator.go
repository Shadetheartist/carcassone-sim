package simulator

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/engine"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Simulator struct {
	Engine    *engine.Engine
	GameData  *data.GameData
	drawData  *DrawData
	playSpeed time.Duration
}

func NewSimulator(engine *engine.Engine) *Simulator {
	sim := &Simulator{}
	sim.Engine = engine
	sim.drawData = &DrawData{}
	sim.playSpeed = 100 * time.Millisecond
	sim.initDraw()

	return sim
}

func (sim *Simulator) Simulate() {
	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")
	ebiten.SetScreenClearedEveryFrame(false)
	ebiten.SetMaxTPS(ebiten.SyncWithFPS)

	if err := ebiten.RunGame(sim); err != nil {
		panic(err)
	}
}

func (sim *Simulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 700, 525
}
