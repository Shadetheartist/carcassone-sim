package tile

import "beeb/carcassonne/directions"

type Feature struct {
	Id     int
	Type   FeatureType
	Shield bool
	Curve  bool
	Edges  []directions.Direction
}

var DefaultFeature = Feature{
	Type:   Grass,
	Shield: false,
}
