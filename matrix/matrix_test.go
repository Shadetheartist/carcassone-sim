package matrix

import (
	"testing"
)

func TestMatrix_Index(t *testing.T) {

	type fields struct {
		Size int
		data []interface{}
	}

	type args struct {
		x int
		y int
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"0", fields{Size: 7}, args{0, 0}, 0},
		{"1", fields{Size: 7}, args{1, 0}, 1},
		{"6", fields{Size: 7}, args{6, 0}, 6},
		{"7", fields{Size: 7}, args{0, 1}, 7},
		{"8", fields{Size: 7}, args{1, 1}, 8},
		{"13", fields{Size: 7}, args{6, 1}, 13},
		{"14", fields{Size: 7}, args{0, 2}, 14},
		{"15", fields{Size: 7}, args{1, 2}, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrix[int](tt.fields.Size)
			if got := m.Index(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("Matrix.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatrix_Transpose(t *testing.T) {
	s := 4
	m := NewMatrix[int](s)

	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			d := m.Index(x, y)
			m.Set(x, y, d)
		}
	}

	preTransposedMatrix := NewMatrix[int](s)

	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			//swap index x and y
			d := preTransposedMatrix.Index(y, x)
			preTransposedMatrix.Set(x, y, d)
		}
	}

	m.Transpose()

	if !m.Equal(preTransposedMatrix) {
		t.Errorf("Transposed matricies are not equal")
	}

}

func TestMatrix_ReverseRows(t *testing.T) {
	s := 3
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}

	preReversedMatrix := NewMatrix[int](s)

	c = 0
	for y := 0; y < s; y++ {
		for x := s - 1; x >= 0; x-- {
			preReversedMatrix.Set(x, y, c)
			c++
		}
	}

	m.ReverseRows()

	if !m.Equal(preReversedMatrix) {
		t.Errorf("Reversed (Rows) matricies are not equal")
	}

}

func TestMatrix_ReverseColumns(t *testing.T) {
	s := 4
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}

	preReversedMatrix := NewMatrix[int](s)

	c = 0
	for y := s - 1; y >= 0; y-- {
		for x := 0; x < s; x++ {
			preReversedMatrix.Set(x, y, c)
			c++
		}
	}

	m.ReverseColumns()

	if !m.Equal(preReversedMatrix) {
		t.Errorf("Reversed (Columns) matricies are not equal")
	}

}

func TestMatrix_Rotation90(t *testing.T) {
	s := 3
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}
	mCopy := m.Copy()

	m.Print()

	m.Rotate90()
	m.Print()

	m.Rotate90()
	m.Print()

	m.Rotate90()
	m.Print()

	m.Rotate90()
	m.Print()

	if !m.Equal(mCopy) {
		t.Errorf("4 90 degree Rotations did not reset the matrix")
	}

}

func TestMatrix_Rotation180(t *testing.T) {
	s := 3
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}
	mCopy := m.Copy()

	m.Print()

	m.Rotate180()
	m.Print()

	m.Rotate180()
	m.Print()

	if !m.Equal(mCopy) {
		t.Errorf("2 180 degree Rotations did not reset the matrix")
	}
}

func TestMatrix_Rotation270(t *testing.T) {
	s := 3
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}
	mCopy := m.Copy()

	m.Print()

	m.Rotate270()
	m.Print()

	m.Rotate90()
	m.Print()

	if !m.Equal(mCopy) {
		t.Errorf("a 270 degree Rotation, followed by a 90 degree Rotation did not reset the matrix")
	}
}

func TestMatrix_Rotation(t *testing.T) {
	s := 3
	m := NewMatrix[int](s)

	c := 0
	for y := 0; y < s; y++ {
		for x := 0; x < s; x++ {
			m.Set(x, y, c)
			c++
		}
	}
	mCopy := m.Copy()

	m.Print()

	m.Rotate(270)
	m.Print()

	m.Rotate(90)
	m.Print()

	if !m.Equal(mCopy) {
		t.Errorf("a 270 degree Rotation, followed by a 90 degree Rotation did not reset the matrix")
	}

	m.Rotate(180)
	m.Print()

	m.Rotate(180)
	m.Print()

	if !m.Equal(mCopy) {
		t.Errorf("a 270 degree Rotation, followed by a 90 degree Rotation did not reset the matrix")
	}
}

func BenchmarkTranspose(b *testing.B) {
	s := 7
	m := NewMatrix[int](s)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Transpose()
	}
}

func BenchmarkRotate90(b *testing.B) {
	s := 7
	m := NewMatrix[int](s)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Rotate90()
	}
}

func BenchmarkIndex(b *testing.B) {
	s := 12
	m := NewMatrix[int](s)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Index(4, 2)
	}
}

func BenchmarkIterate(b *testing.B) {
	s := 7
	m := NewMatrix[int](s)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		m.Iterate(func(rt int, x int, y int, idx int) {

		})
	}
}
