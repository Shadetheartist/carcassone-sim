package tile

import "github.com/google/uuid"

type FeatureType int

var featureTypeStrMap []string = []string{
	"None",
	"Farm",
	"Road",
	"Castle",
	"Cloister",
	"River",
	"Shield",
}

var featureTypeScoreMap []int = []int{
	0, //"None",
	0, //"Farm",
	1, //"Road",
	2, //"Castle",
	0, //"Cloister",
	0, //"River",
	1, //"Shield",
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
	return featureTypeStrMap[ft]
}

func (ft FeatureType) Score() int {
	return featureTypeScoreMap[ft]
}

type Feature struct {
	Id                     uuid.UUID
	Type                   FeatureType
	ParentTile             *Tile
	ParentRefenceTileGroup *ReferenceTileGroup
	ParentFeature          *Feature
	//both key and value are the same
	Links map[*Feature]*Feature
}

func (f *Feature) String() string {
	return f.Type.String()
}
