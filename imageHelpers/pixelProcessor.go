package imageHelpers

import "image"

//used in fill algorithm, return true if you want to visit the neighbours of the pixel in context
type FillProcessor func(image.Image, image.Point, int) bool
