package tile

type FeatureType int

var featureTypeMap []string = []string{
	"Farm",
	"Road",
	"Castle",
	"Cloister",
	"River",
	"Shield",
}

const (
	Farm FeatureType = iota
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
	Type              FeatureType
	ParentTile        *Tile
	ParentRefenceTile *ReferenceTile
}

func (f *Feature) String() string {
	return f.Type.String()
}
