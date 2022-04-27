package tile

import (
	"beeb/carcassonne/matrix"
	"image"
)

type ReferenceTile struct {
	Name          string
	Orientation   int
	FeatureMatrix *matrix.Matrix[*Feature]
	Features      []*Feature
	Image         image.Image
	EdgeFeatures  EdgeArray[*Feature]
}

type Tile struct {
}
