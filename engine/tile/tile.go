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
	AvgFeaturePos map[*Feature]util.Point[float64]
}

type Tile struct {
	Id                  uuid.UUID
	Position            util.Point[int]
	Reference           *ReferenceTile
	Features            []*Feature
	EdgeFeatures        *EdgeArray[*Feature]
	Neighbours          *EdgeArray[*Tile]
	ReferenceFeatureMap map[*Feature]*Feature
}

func (rtg *ReferenceTileGroup) IsRiverTile() bool {
	return strings.Contains(rtg.Name, "River")
}

func (rtg *ReferenceTileGroup) IsRiverTerminus() bool {
	return rtg.Name == "RiverTerminus"
}
