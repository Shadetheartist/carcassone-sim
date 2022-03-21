package tile

import "image"

type Orientation uint16

const (
	OrientationZero       Orientation = 0
	OrientationNinety                 = 90
	OrientationOneEight               = 180
	OrientationTwoSeventy             = 270
)

type Tile struct {
	Name        string
	Image       image.Image
	Orientation Orientation
}

func (t Tile) String() string {
	return t.Name
}
