package imageHelpers

import "image"

func PointToIndex(img image.Image, p image.Point) int {
	return (p.Y * img.Bounds().Dx()) + p.X
}

func IndexToPoint(img image.Image, idx int) image.Point {
	return image.Point{
		X: idx % img.Bounds().Dx(),
		Y: idx / img.Bounds().Dx(),
	}
}
