package tile

type FeatureType int

const (
	None FeatureType = iota
	Grass
	River
	Castle
	Road
	Cloister
)

var FeatureTypeStrMap = map[string]FeatureType{
	"none":     None,
	"grass":    Grass,
	"river":    River,
	"castle":   Castle,
	"road":     Road,
	"cloister": Cloister,
}
