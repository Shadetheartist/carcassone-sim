package tile

import (
	"beeb/carcassonne/directions"
	"errors"
	"image"
	"image/color"
	"math"
)

type FarmSegment struct {
	Parent     *Tile
	Neighbours []*FarmSegment
	AvgPoint   image.Point
}

func AvgFarmSegmentPos(fs *FarmSegment, matrix [][]*FarmSegment) image.Point {

	//calculating the center of the segment via averaging

	totalX := 0
	totalY := 0
	n := 0

	for y := range matrix {
		for x := range matrix[y] {
			p := matrix[y][x]
			if p == fs {
				totalX += x
				totalY += y
				n++
			}
		}
	}

	//no dividing by zero!
	if n == 0 {
		return image.Point{-1, -1}
	}

	avgX := totalX / n
	avgY := totalY / n
	return image.Point{avgX, avgY}
}

func (t *Tile) IntegrateFarms() {
	matrix := OrientedFarmMatrix(t, int(t.Placement.Orientation))
	edgePix := edgePix(t.Image.Bounds())

	for dirInt, n := range t.Neighbours {

		if n == nil {
			continue
		}

		dir := directions.Direction(dirInt)
		complimentDir := directions.Compliment[dir]

		pix := edgePix[dir]
		complimentPix := edgePix[complimentDir]

		neighbourFarmMatrix := OrientedFarmMatrix(n, int(n.Placement.Orientation))

		for i := range pix {
			farmSegment := matrix[pix[i].Y][pix[i].X]

			if farmSegment == nil {
				continue
			}

			neighbouringFarmSegment := neighbourFarmMatrix[complimentPix[i].Y][complimentPix[i].X]

			if neighbouringFarmSegment == nil {
				continue
			}

			linkFarmSegments(farmSegment, neighbouringFarmSegment)
		}
	}
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
		farmSegment.Neighbours = make([]*FarmSegment, 0, 4)
		t.FarmSegments = append(t.FarmSegments, &farmSegment)

		fillSegment(img, nextIndex, visited, farmMatrix, &farmSegment)
	}

	return farmMatrix
}

func areFarmSegmentsLinked(fsA *FarmSegment, fsB *FarmSegment) bool {
	for _, n := range fsA.Neighbours {
		if n == fsB {
			return true
		}
	}

	return false
}

func linkFarmSegments(fsA *FarmSegment, fsB *FarmSegment) {
	if !areFarmSegmentsLinked(fsA, fsB) {
		fsA.Neighbours = append(fsA.Neighbours, fsB)
	}

	if !areFarmSegmentsLinked(fsB, fsA) {
		fsB.Neighbours = append(fsB.Neighbours, fsA)
	}

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

func fillSegment(img image.Image, initialIndex int, visited []bool, farmMatrix [][]*FarmSegment, segment *FarmSegment) {
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

func edgePix(rect image.Rectangle) [][]image.Point {
	edgePix := make([][]image.Point, 4)

	edgePix[0] = northPix(rect)
	edgePix[1] = eastPix(rect)
	edgePix[2] = southPix(rect)
	edgePix[3] = westPix(rect)

	return edgePix
}

func northPix(rect image.Rectangle) []image.Point {
	pix := make([]image.Point, rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = image.Point{X: i, Y: 0}
	}

	return pix
}

func southPix(rect image.Rectangle) []image.Point {
	pix := make([]image.Point, rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = image.Point{X: i, Y: rect.Max.X - 1}
	}

	return pix
}

func westPix(rect image.Rectangle) []image.Point {
	pix := make([]image.Point, rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = image.Point{X: 0, Y: i}
	}

	return pix
}

func eastPix(rect image.Rectangle) []image.Point {
	pix := make([]image.Point, rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = image.Point{X: rect.Max.X - 1, Y: i}
	}

	return pix
}
