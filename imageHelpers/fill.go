package imageHelpers

import (
	"image"
)

func Fill(img image.Image, p image.Point, orthogonal bool, processor FillProcessor) {
	imageSize := img.Bounds().Dx() * img.Bounds().Dy()
	stackCapacity := imageSize / 2

	//bool array to facilitate only visiting each pixel once
	visited := make([]bool, imageSize)

	var pt image.Point
	stack := make([]image.Point, 0, stackCapacity)

	//initial stack element
	stack = append(stack, p)

	for len(stack) > 0 {

		// pop off stack
		pt, stack = stack[len(stack)-1], stack[:len(stack)-1]

		idx := PointToIndex(img, pt)

		//avoid revisiting pixels
		if visited[idx] {
			continue
		}

		//keep track of what pixels have been visited already
		visited[idx] = true

		//processor returens if true if neighbours should be added to the stack
		if processor(img, pt, idx) {
			//add neighbours to stack
			if orthogonal {
				stack = append(stack, OrthogonalNeighbours(img, pt)...)
			} else {
				stack = append(stack, Neighbours(img, pt)...)
			}
		}
	}
}
