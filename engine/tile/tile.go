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
	Features     []*Feature
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
	Id           uuid.UUID
	Position     util.Point[int]
	Reference    *ReferenceTile
	Features     []*Feature
	EdgeFeatures *EdgeArray[*Feature]
	Neighbours   *EdgeArray[*Tile]
	//this one points from the reference tile feature to this tiles' feature
	ReferenceFeatureMap map[*Feature]*Feature
}

func (rtg *ReferenceTileGroup) IsRiverTile() bool {
	return strings.Contains(rtg.Name, "River")
}

func (rtg *ReferenceTileGroup) IsRiverTerminus() bool {
	return rtg.Name == "RiverTerminus"
}

func (t *Tile) HasFeature(f *Feature) bool {
	for _, tf := range t.Features {
		if tf == f {
			return true
		}
	}

	return false
}
