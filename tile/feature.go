package tile

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
	Type              FeatureType
	ParentTile        *Tile
	ParentRefenceTile *ReferenceTile
}

func (f *Feature) String() string {
	return f.Type.String()
}
