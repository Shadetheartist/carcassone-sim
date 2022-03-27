package tile

import (
	"beeb/carcassonne/directions"
	"fmt"
)

type Placement struct {
	Position    Position
	Orientation uint16
}

func (p Placement) TileDirection(direction directions.Direction) directions.Direction {

	shift := int(p.Orientation) / 90

	//this will get us an int 0-3, which we can add to our
	//shift subtraction from our edge (plus 4, mod 4)
	//to basically rotate the edge clockwise
	dirInt := int(direction)
	shiftedDirInt := (dirInt + 4 - shift) % 4

	return directions.Direction(shiftedDirInt)
}

func (pl *Placement) String() string {
	return fmt.Sprint("X: ", pl.Position.X, " Y: ", pl.Position.Y, " O: ", pl.Orientation)
}
