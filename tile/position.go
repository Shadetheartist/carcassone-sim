package tile

import "beeb/carcassonne/directions"

type Position struct {
	X int
	Y int
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
