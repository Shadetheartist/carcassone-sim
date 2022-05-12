package util

import "beeb/carcassonne/util/directions"

type numeric interface {
	int | float64
}

type Point[T numeric] struct {
	X T
	Y T
}

func (pt Point[T]) OrthogonalNeighbours() {
}

func (p Point[T]) North() Point[T] {
	return Point[T]{
		X: p.X,
		Y: p.Y - 1,
	}
}

func (p Point[T]) South() Point[T] {
	return Point[T]{
		X: p.X,
		Y: p.Y + 1,
	}
}

func (p Point[T]) East() Point[T] {
	return Point[T]{
		X: p.X + 1,
		Y: p.Y,
	}
}

func (p Point[T]) West() Point[T] {
	return Point[T]{
		X: p.X - 1,
		Y: p.Y,
	}
}

func (p Point[T]) EdgePos(dir directions.Direction) Point[T] {
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
