package util

import (
	"beeb/carcassonne/util/directions"
	"fmt"
	"math"
)

type numeric interface {
	int | float64
}

type Point[T numeric] struct {
	X T
	Y T
}

func (p Point[T]) OrthogonalNeighbours() [4]Point[T] {
	return [4]Point[T]{
		p.North(),
		p.South(),
		p.East(),
		p.West(),
	}
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

func (p Point[T]) Add(pt Point[T]) Point[T] {
	return Point[T]{
		X: p.X + pt.X,
		Y: p.Y + pt.Y,
	}
}

func (p Point[T]) Subtract(pt Point[T]) Point[T] {
	return Point[T]{
		X: p.X - pt.X,
		Y: p.Y - pt.Y,
	}
}

func (p Point[T]) Magnitude() float64 {
	return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

func (p Point[T]) String() string {
	return fmt.Sprintf("(%d, %d)", p.X, p.Y)
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
