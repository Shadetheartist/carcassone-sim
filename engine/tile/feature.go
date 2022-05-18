package tile

import "github.com/google/uuid"

type FeatureType int

var featureTypeMap []string = []string{
	"None",
	"Farm",
	"Road",
	"Castle",
	"Cloister",
	"River",
	"Shield",
}

const (
	None FeatureType = iota
	Farm
	Road
	Castle
	Cloister
	River
	Shield
)

func (ft FeatureType) String() string {
	return featureTypeMap[ft]
}

type Feature struct {
	Id                uuid.UUID
	Type              FeatureType
	ParentTile        *Tile
	ParentRefenceTile *ReferenceTile
	Links             map[*Feature]*Feature
}

func (f *Feature) String() string {
	return f.Type.String()
}
