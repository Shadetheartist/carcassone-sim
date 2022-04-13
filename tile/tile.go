package tile

import (
	"beeb/carcassonne/directions"
	"image"
)

type Edge struct {
	Direction directions.Direction
	Parent    *Tile
	Feature   *Feature
	Neighbour *Edge
}

type Tile struct {
	Name      string
	Image     image.Image
	Placement Placement
	Rendered  bool

	//dir mapped
	Neighbours []*Tile

	//this is used in tight loops a LOT so it is cached here
	CachedEdgeFeatureTypes [4]FeatureType

	RoadSegments [4]*RoadSegment
	FarmMatrix   [][]*FarmSegment
	FarmSegments []*FarmSegment
	Edges2       [4]Edge
}

func (t *Tile) Feature(direction directions.Direction) *Feature {
	return t.Edges2[direction].Feature
}

func (t *Tile) FeaturesByType(ft FeatureType) []*Feature {
	features := make([]*Feature, 0, 4)

	for _, e := range t.Edges2 {
		if e.Feature != nil && e.Feature.Type == ft {
			features = append(features, e.Feature)
		}
	}

	return features
}

func (t *Tile) CacheEdgeFeatureTypes() [4]FeatureType {

	var ef [4]FeatureType

	for _, d := range directions.List {
		if f := t.Edges2[d].Feature; f != nil {
			ef[d] = f.Type
		}
	}

	return ef
}

func (t *Tile) String() string {
	return t.Name
}
