package tile

import (
	"beeb/carcassonne/directions"
	"fmt"
	"image"
)

const (
	OrientationZero       uint16 = 0
	OrientationNinety            = 90
	OrientationOneEight          = 180
	OrientationTwoSeventy        = 270
)

var OrientationList = [...]uint16{0, 90, 180, 270}

type Position struct {
	X int
	Y int
}

type Placement struct {
	Position    Position
	Orientation uint16
}

func (p Position) North() Position {
	return Position{
		X: p.X,
		Y: p.Y - 1,
	}
}

func (p Position) South() Position {
	return Position{
		X: p.X,
		Y: p.Y + 1,
	}
}

func (p Position) East() Position {
	return Position{
		X: p.X + 1,
		Y: p.Y,
	}
}

func (p Position) West() Position {
	return Position{
		X: p.X - 1,
		Y: p.Y,
	}
}

type Feature struct {
	Type   FeatureType
	Shield bool
	Curve  bool
}

var DefaultFeature = Feature{
	Type:   "grass",
	Shield: false,
}

type Tile struct {
	Name      string
	Image     image.Image
	Placement Placement
	IsPlaced  bool

	//these are non-oriented (board reference)
	Features   map[int]Feature
	Edges      map[directions.Direction]int
	Neighbours map[directions.Direction]*Tile
}

func (t *Tile) EdgeFeatures() map[directions.Direction]Feature {
	edgeFeatures := make(map[directions.Direction]Feature)

	for edge, featureId := range t.Edges {
		if feature, exists := t.Features[featureId]; exists {
			edgeFeatures[edge] = feature
		}

		edgeFeatures[edge] = DefaultFeature
	}

	return edgeFeatures
}

func (t *Tile) Feature(direction directions.Direction) Feature {
	if edge, exists := t.Edges[direction]; exists {
		if feature, exists := t.Features[edge]; exists {
			return feature
		}

		panic(fmt.Sprint("Edge does not have a corresponding feature mapped. ", edge))
	}

	return DefaultFeature
}

func (p Placement) TileDirection(direction directions.Direction) directions.Direction {
	if p.Orientation == 0 {
		return direction
	}

	shift := int(p.Orientation) / 90

	//this will get us an int 0-3, which we can add to our
	//shift subtraction from our edge (plus 4, mod 4)
	//to basically rotate the edge clockwise
	dirInt := directions.StrMap[direction]
	shiftedDirInt := (4 + dirInt - shift) % 4

	return directions.IntMap[shiftedDirInt]
}

func (p Position) EdgePos(dir directions.Direction) Position {
	switch dir {
	case directions.North:
		return p.North()
	case directions.South:
		return p.South()
	case directions.West:
		return p.West()
	case directions.East:
		return p.East()
	default:
		panic("That's not a direction")
	}
}

func (t *Tile) FreeFeatureEdge(featureType string) directions.Direction {
	//for dir, feature := range t.Features {}
	return directions.North
}

func (t *Tile) String() string {
	return t.Name
}
func (pl *Placement) String() string {
	return fmt.Sprint("X: ", pl.Position.X, " Y: ", pl.Position.Y, " O: ", pl.Orientation)
}
