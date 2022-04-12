package tile

import (
	"errors"
	"image"
	"image/color"
	"math"
)

type FarmSegment struct {
	Parent *Tile
}

func (t *Tile) IntegrateFarms() {

}

func TransposeFarmMatrix(matrix [][]*FarmSegment) [][]*FarmSegment {

	newMatrix := make([][]*FarmSegment, len(matrix))

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			newMatrix[j] = append(newMatrix[j], matrix[i][j])
		}
	}

	return newMatrix
}

func OrientedFarmMatrix(t *Tile, r int) [][]*FarmSegment {
	n := int(math.Abs(float64(r))) / 90 % 4

	var matrix [][]*FarmSegment

	switch n {
	case 0:
		//absolute dogshit way to return a copy of the original matrix
		matrix = TransposeFarmMatrix(t.FarmMatrix)
		matrix = TransposeFarmMatrix(matrix)

	case 1:
		{
			//transpose & reverse each row = 90 deg
			//this copies the existing matrix
			matrix = TransposeFarmMatrix(t.FarmMatrix)

			//reverse each row
			for _, row := range matrix {
				for i, j := 0, len(row)-1; i < j; i, j = i+1, j-1 {
					row[i], row[j] = row[j], row[i]
				}
			}
		}

	case 2:
		{
			//transpose = 180 deg
			//this copies the existing matrix
			matrix = TransposeFarmMatrix(t.FarmMatrix)
		}

	case 3:
		{
			//transpose & reverse each column = 270 / -90 deg
			//this copies the existing matrix
			matrix = TransposeFarmMatrix(t.FarmMatrix)

			//reverse each column
			for i, j := 0, len(matrix)-1; i < j; i, j = i+1, j-1 {
				matrix[i], matrix[j] = matrix[j], matrix[i]
			}

		}
	}

	return matrix
}

func ComputeFarmMatrix(t *Tile) [][]*FarmSegment {

	img := t.Image
	edgePositions := EdgePositions(img)

	imgLen := img.Bounds().Max.X * img.Bounds().Max.Y

	farmMatrix := make([][]*FarmSegment, img.Bounds().Max.Y)

	for y := 0; y < img.Bounds().Max.Y; y++ {
		farmMatrix[y] = make([]*FarmSegment, img.Bounds().Max.X)
	}

	visited := make([]bool, imgLen)

	for !allVisited(visited) {

		nextIndex, _, err := getUnvisitedFarmPos(img, visited)

		if err != nil {
			break
		}

		farmSegment := FarmSegment{}
		farmSegment.Parent = t

		fillSegment(img, edgePositions, nextIndex, visited, farmMatrix, &farmSegment)
	}

	return farmMatrix
}

func getUnvisitedFarmPos(img image.Image, visited []bool) (int, Position, error) {
	for i, v := range visited {
		if v {
			continue
		}

		pos := indexToPos(img, i)

		color := img.At(pos.X, pos.Y)

		if isFarmColor(color) {
			return i, pos, nil
		}
	}

	return 0, Position{}, errors.New("No more unvisited positions")
}

func fillSegment(img image.Image, edgePositions []Position, initialIndex int, visited []bool, farmMatrix [][]*FarmSegment, segment *FarmSegment) {
	stack := make([]int, 0, 8)
	stack = append(stack, initialIndex)

	var i int

	for len(stack) > 0 {

		// pop off stack
		i, stack = stack[len(stack)-1], stack[:len(stack)-1]

		if visited[i] {
			continue
		}

		visited[i] = true

		pos := indexToPos(img, i)
		color := img.At(pos.X, pos.Y)

		if isFarmColor(color) {
			farmMatrix[pos.Y][pos.X] = segment
		}

		//add neighbours to stack
		neighbours := indexNeighbours(img, i)
		for _, n := range neighbours {
			p := indexToPos(img, n)
			c := img.At(p.X, p.Y)
			if !visited[n] && isFarmColor(c) {
				stack = append(stack, n)
			}
		}
	}
}

func isEdgePosition(edgePositions []Position, pos Position) bool {
	for _, p := range edgePositions {
		if p == pos {
			return true
		}
	}

	return false
}

func indexNeighbours(img image.Image, i int) []int {
	neighbours := make([]int, 0, 4)
	if neighbour, valid := northIndex(img, i); valid {
		neighbours = append(neighbours, neighbour)
	}

	if neighbour, valid := southIndex(img, i); valid {
		neighbours = append(neighbours, neighbour)
	}

	if neighbour, valid := eastIndex(img, i); valid {
		neighbours = append(neighbours, neighbour)
	}

	if neighbour, valid := westIndex(img, i); valid {
		neighbours = append(neighbours, neighbour)
	}

	return neighbours
}

func northIndex(img image.Image, i int) (int, bool) {
	_i := i - img.Bounds().Max.X

	if _i < 0 {
		return 0, false
	}

	return _i, true
}

func southIndex(img image.Image, i int) (int, bool) {
	max := img.Bounds().Max.X * img.Bounds().Max.Y
	_i := i + img.Bounds().Max.X

	if _i >= max {
		return 0, false
	}

	return _i, true
}

func eastIndex(img image.Image, i int) (int, bool) {
	rowIndex := i % img.Bounds().Max.X

	if rowIndex-1 < 0 {
		return 0, false
	}

	return i - 1, true
}

func westIndex(img image.Image, i int) (int, bool) {
	rowIndex := i % img.Bounds().Max.X

	if rowIndex+1 >= img.Bounds().Max.X {
		return 0, false
	}

	return i + 1, true
}

func indexToPos(img image.Image, index int) Position {
	x := index % img.Bounds().Max.X
	y := (index / img.Bounds().Max.X)

	return Position{X: x, Y: y}
}

func allVisited(seen []bool) bool {
	for _, s := range seen {
		if !s {
			return false
		}
	}

	return true
}

var farmColorGreen uint32 = 190 | 190<<8

func isFarmColor(c color.Color) bool {
	_, g, _, _ := c.RGBA()
	return g == farmColorGreen
}

//creates a slice of positions that describe the outer edge of a square input image
//elements of the slice are in clockwise order starting from the top left corner
func EdgePositions(img image.Image) []Position {
	//these images are square, no need for a Y
	maxSize := img.Bounds().Max.X

	//edge lengths with corner pixels removed * num edges, then re add corner pixel count
	numEdgePix := (maxSize-2)*4 + 4
	edgePixels := make([]Position, 0, numEdgePix)

	//we're going clockwise around the edge to create the slice

	//top row
	for i := 0; i < maxSize; i++ {
		edgePixels = append(edgePixels, Position{X: i, Y: 0})
	}

	//right side, without top pos, which has already been added
	for i := 1; i < maxSize; i++ {
		edgePixels = append(edgePixels, Position{X: maxSize - 1, Y: i})
	}

	//bottom row, going backward,
	for i := maxSize - 2; i >= 0; i-- {
		edgePixels = append(edgePixels, Position{X: i, Y: maxSize - 1})
	}

	//left side, without top and bottom pos, which have already been added, going backward
	for i := maxSize - 2; i > 0; i-- {
		edgePixels = append(edgePixels, Position{X: 0, Y: i})
	}

	return edgePixels
}
