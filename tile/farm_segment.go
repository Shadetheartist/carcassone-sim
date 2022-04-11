package tile

import (
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

	forwardCounter := 0
	visited := make([]bool, len(edgePositions))

	var positions []Position

	for !allVisited(visited) {
		forwardCounter, positions = segmentEdgePositions(img, edgePositions, forwardCounter, visited)

		farmSegment := FarmSegment{}
		farmSegment.EdgePixels = positions
		farmSegment.Parent = &t

		farmSegments = append(farmSegments, farmSegment)
	}

	return farmSegments
}

func allVisited(seen []bool) bool {
	for _, s := range seen {
		if !s {
			return false
		}
	}

	return true
}

func segmentEdgePositions(img image.Image, edgePositions []Position, startIndex int, visited []bool) (int, []Position) {

	forwardCounter := startIndex
	backwardCounter := backCounterIndex(startIndex-1, len(edgePositions)-1)
	farmEdges := make([]Position, 0, len(edgePositions))

	for forwardCounter < len(edgePositions) {

		if visited[forwardCounter] {
			break
		}

		visited[forwardCounter] = true

		pos := edgePositions[forwardCounter]

		color := img.At(pos.X, pos.Y)

		if isFarmColor(color) {
			//farm edge is continuing
			farmEdges = append(farmEdges, pos)
			forwardCounter++
		} else {
			break
		}
	}

	for backwardCounter > 0 {

		if visited[backwardCounter] {
			break
		}

		visited[backwardCounter] = true

		pos := edgePositions[backwardCounter]

		color := img.At(pos.X, pos.Y)

		if isFarmColor(color) {
			farmEdges = append(farmEdges, pos)
			//farm edge is continuing
			backwardCounter = backCounterIndex(backwardCounter-1, len(edgePositions)-1)
		} else {
			break
		}
	}

	return forwardCounter + 1, farmEdges
}

func backCounterIndex(i int, max int) int {

	if i < 0 {
		i += max + 1
	}

	return i
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
