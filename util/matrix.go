package util

type Matrix struct {
	values []interface{}
}

func NewMatrix(size int, capacity int) *Matrix {
	matrix := Matrix{}

	matrix.values = make([]interface{}, size, capacity)

	return &matrix
}

func (m Matrix) Copy() *Matrix {

	newMatrix := NewMatrix(len(m.values), cap(m.values))

	copy(newMatrix.values, m.values)

	return newMatrix
}

func (m Matrix) Rotate() {

}

func (m Matrix) rotate90() {

}

func (m Matrix) rotate180() {

}

func (m Matrix) rotate270() {

}

func (m Matrix) Index(x int, y int) interface {
}

func (m Matrix) transpose() {
	newMatrix := make([][]*FarmSegment, len(matrix))

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {

			newMatrix[j] = append(newMatrix[j], matrix[i][j])
		}
	}

	return newMatrix
}
