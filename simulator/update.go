package simulator

import (
	"beeb/carcassonne/util"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

var mouseDown bool = false
var mouseInitialX int = 0
var mouseInitialY int = 0
var rMouseDown = false

func (sim *Simulator) Update() error {

	cX, cY := ebiten.CursorPosition()

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseDown = false
	}

	if !mouseDown && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseInitialX, mouseInitialY = ebiten.CursorPosition()
		mouseDown = true
		e := sim.GetObjectUnderCursor()
		fmt.Println(e)
	}

	if mouseDown {

		sim.drawData.cameraOffset.X += cX - mouseInitialX
		sim.drawData.cameraOffset.Y += cY - mouseInitialY

		mouseInitialX = cX
		mouseInitialY = cY
	}

	_, wheel := ebiten.Wheel()
	if wheel != 0 {
		if wheel > 0 {
			sim.drawData.boardScale *= 2
		} else {
			sim.drawData.boardScale /= 2
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {

		sim.Engine.InitGame()
		steps := (sim.Engine.RiverDeck.Remaining() + sim.Engine.Deck.Remaining()) * 5
		for i := 0; i < steps; i++ {
			sim.Engine.Step()
		}

		rMouseDown = true
	} else {
		rMouseDown = false
	}

	return nil
}

func (sim *Simulator) GetObjectUnderCursor() *interface{} {
	cursorX, cursorY := ebiten.CursorPosition()
	ssPoint := util.Point[int]{X: cursorX, Y: cursorY}
	wPoint := sim.toWorldSpace(ssPoint)
	bPoint := sim.toBoardSpace(ssPoint)

	fmt.Println(ssPoint, wPoint, bPoint)

	return nil
}
