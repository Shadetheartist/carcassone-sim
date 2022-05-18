package tile

import (
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/util"
	"image"
	"strings"

	"github.com/google/uuid"
)

type ReferenceTileGroup struct {
	Name         string
	Orientations []*ReferenceTile
}

type ReferenceTile struct {
	Name          string
	Orientation   int
	FeatureMatrix *matrix.Matrix[*Feature]
	Features      []*Feature
	Image         image.Image
	EdgeFeatures  *EdgeArray[*Feature]
	EdgeSignature *EdgeSignature
}

type Tile struct {
	Id            uuid.UUID
	Position      util.Point[int]
	Reference     *ReferenceTile
	FeatureMatrix *matrix.Matrix[*Feature]
	Features      []*Feature
	EdgeFeatures  *EdgeArray[*Feature]
	Neighbours    *EdgeArray[*Tile]
}

func (rtg *ReferenceTileGroup) IsRiverTile() bool {
	return strings.Contains(rtg.Name, "River")
}

func (rtg *ReferenceTileGroup) IsRiverTerminus() bool {
	return rtg.Name == "RiverTerminus"
}
