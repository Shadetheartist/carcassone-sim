package imageHelpers

import (
	"image"
)

//return true if orthogonal neighbour behaviour is desired
type SegmentCallback func(image.Image, image.Point, int) bool

func FillAll(img image.Image, segmentCallback SegmentCallback, processor FillProcessor) {
	imageSize := img.Bounds().Dx() * img.Bounds().Dy()
	visited := make([]bool, imageSize)

	for nextIdx := getNextUnvisited(visited); nextIdx != -1; nextIdx = getNextUnvisited(visited) {
		nextPoint := IndexToPoint(img, nextIdx)

		orthogonal := segmentCallback(img, nextPoint, nextIdx)

		Fill(img, nextPoint, orthogonal, func(img image.Image, p image.Point, idx int) bool {

			continueProcessing := processor(img, p, idx)

			if continueProcessing {
				visited[idx] = true
			}

			return continueProcessing
		})
	}
}

func getNextUnvisited(visited []bool) int {
	for i, b := range visited {
		if !b {
			return i
		}
	}

	return -1
}
