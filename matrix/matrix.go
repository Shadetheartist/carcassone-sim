package matrix

import (
	"beeb/carcassonne/util"
	"errors"
	"fmt"
)

type Matrix[T comparable] struct {
	size int
	data []T
}

func NewMatrix[T comparable](size int) *Matrix[T] {
	matrix := Matrix[T]{}

	matrix.size = size
	matrix.data = make([]T, size*size)

	return &matrix
}

func (m *Matrix[T]) Copy() *Matrix[T] {
	newMatrix := NewMatrix[T](m.size)

	copy(newMatrix.data, m.data)

	return newMatrix
}

func (m *Matrix[T]) Rotate(n int) {

	n = n % 360

	switch n {
	case 0:
		return
	case -270:
	case 90:
		m.Rotate90()
		return
	case -180:
	case 180:
		m.Rotate180()
		return
	case -90:
	case 270:
		m.Rotate270()
		return
	default:
		panic(fmt.Sprintln("Invalid rotation ", n, ". All rotations must be multiples of 90"))
	}
}

func (m *Matrix[T]) Rotate90() {
	m.Transpose()
	m.ReverseRows()
}

func (m *Matrix[T]) Rotate180() {
	m.ReverseRows()
	m.ReverseColumns()
}

func (m *Matrix[T]) Rotate270() {
	m.ReverseRows()
	m.Transpose()
}

func (m *Matrix[T]) Index(x int, y int) int {
	return (y * m.size) + x
}

func (m *Matrix[T]) Get(x int, y int) T {
	i := m.Index(x, y)

	return m.data[i]
}

func (m *Matrix[T]) GetPt(pt util.Point[int]) (T, error) {

	if !m.IsInBounds(pt.X, pt.Y) {
		var noop T
		return noop, errors.New("point not within bounds")
	}

	i := m.Index(pt.X, pt.Y)

	return m.data[i], nil
}

func (m *Matrix[T]) GetI(index int) T {
	return m.data[index]
}

func (m *Matrix[T]) Set(x int, y int, d T) {
	i := m.Index(x, y)
	m.data[i] = d
}

func (m *Matrix[T]) SetPt(pt util.Point[int], d T) error {
	if !m.IsInBounds(pt.X, pt.Y) {
		return errors.New("point not within bounds")
	}

	i := m.Index(pt.X, pt.Y)

	m.data[i] = d

	return nil
}

func (m *Matrix[T]) SetI(i int, d T) {
	m.data[i] = d
}

func (m *Matrix[T]) Transpose() {
	transposedData := make([]T, len(m.data))

	for y := 0; y < m.size; y++ {
		for x := 0; x < m.size; x++ {
			index := m.Index(x, y)
			indexTransposed := m.Index(y, x)
			transposedData[indexTransposed] = m.data[index]
		}
	}

	m.data = transposedData
}

func (m *Matrix[T]) ReverseRows() {
	for y := 0; y < m.size; y++ {
		rowIndex := y * m.size
		x1 := rowIndex
		x2 := rowIndex + m.size - 1

		for x1 < x2 {
			m.data[x1], m.data[x2] = m.data[x2], m.data[x1]
			x1 = x1 + 1
			x2 = x2 - 1
		}
	}
}

func (m *Matrix[T]) ReverseColumns() {
	for x := 0; x < m.size; x++ {
		colIndex := x
		y1 := colIndex
		y2 := colIndex + (m.size-1)*(m.size)

		for y1 < y2 {
			m.data[y1], m.data[y2] = m.data[y2], m.data[y1]
			y1 = y1 + m.size
			y2 = y2 - m.size
		}
	}
}

func (m *Matrix[T]) Equal(m2 *Matrix[T]) bool {
	for y := 0; y < m.size; y++ {
		for x := 0; x < m.size; x++ {
			index := m.Index(x, y)
			if m.data[index] != m2.data[index] {
				return false
			}
		}
	}

	return true
}

func (m *Matrix[T]) Print() {
	fmt.Println()

	for y := 0; y < m.size; y++ {
		for x := 0; x < m.size; x++ {
			d := m.Get(x, y)
			fmt.Print(d, "\t")
		}
		fmt.Println()
	}

	fmt.Println()
}

type MatrixIterator[T comparable] func(rt T, x int, y int, idx int)

func (m *Matrix[T]) Iterate(iter MatrixIterator[T]) {
	for y := 0; y < m.size; y++ {
		for x := 0; x < m.size; x++ {
			d := m.Get(x, y)
			idx := m.Index(x, y)
			iter(d, x, y, idx)
		}
	}
}

func (m *Matrix[T]) Size() int {
	return m.size
}

func (m *Matrix[T]) Len() int {
	return len(m.data)
}

func (m *Matrix[T]) IsInBounds(x int, y int) bool {
	return x >= 0 && x < m.size && y >= 0 && y < m.size
}

func (m *Matrix[T]) IsIndexInBounds(i int) bool {
	return i >= 0 && i < len(m.data)
}
