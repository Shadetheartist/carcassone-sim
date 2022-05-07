package tile

import (
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/util"
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
	Position      util.Position
	Reference     *ReferenceTile
	FeatureMatrix *matrix.Matrix[*Feature]
	Features      []*Feature
	EdgeFeatures  EdgeArray[*Feature]
}
