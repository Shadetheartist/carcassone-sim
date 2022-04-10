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

	//these are non-oriented (board reference)
	Features map[int]*Feature
	Edges    map[directions.Direction]int

	//dir mapped
	Neighbours []*Tile

	//this is used in tight loops a LOT so it is cached here
	CachedEdgeFeatureTypes [4]FeatureType

	RoadSegments [4]*RoadSegment
	Edges2       [4]Edge
}

func (t *Tile) Feature(direction directions.Direction) *Feature {
	return t.Edges2[direction].Feature
}

func (t *Tile) FeaturesByType(ft FeatureType) []*Feature {
	features := make([]*Feature, 0, 1)

	for _, f := range t.Features {
		if f.Type == ft {
			features = append(features, f)
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
