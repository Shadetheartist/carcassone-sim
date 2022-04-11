package tile

import (
	"errors"
	"image"
	"image/color"
)

type FarmSegment struct {
	Parent     *Tile
	EdgePixels []Position
}

type pix struct {
	neighbours [4]*pix
}

func ComputeFarmSegments(t Tile) []FarmSegment {

	img := t.Image
	edgePositions := EdgePositions(img)
	farmSegments := make([]FarmSegment, 0, 4)

	imgLen := img.Bounds().Max.X * img.Bounds().Max.Y

	visited := make([]bool, imgLen)

	var positions []Position

	for !allVisited(visited) {

		nextIndex, _, err := getUnvisitedFarmPos(img, visited)

		if err != nil {
			break
		}

		positions = fillSegment(img, edgePositions, nextIndex, visited)

		farmSegment := FarmSegment{}
		farmSegment.EdgePixels = positions
		farmSegment.Parent = &t

		farmSegments = append(farmSegments, farmSegment)
	}

	return farmSegments
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

func fillSegment(img image.Image, edgePositions []Position, initialIndex int, visited []bool) []Position {
	segmentEdgePositions := make([]Position, 0, 8)

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

		if isFarmColor(color) && isEdgePosition(edgePositions, pos) {
			segmentEdgePositions = append(segmentEdgePositions, pos)
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

	return segmentEdgePositions

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
