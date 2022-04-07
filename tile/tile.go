package tile

import (
	"beeb/carcassonne/directions"
	"fmt"
	"image"
)

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
	EdgeFeatureTypes [4]FeatureType
	EdgeFeatures     [4]*Feature

	RoadSegments [4]*RoadSegment
}

func (t *Tile) Feature(direction directions.Direction) *Feature {
	return t.EdgeFeatures[direction]
}

func (t *Tile) ComputeFeature(direction directions.Direction) *Feature {
	if edge, exists := t.Edges[direction]; exists {
		if feature, exists := t.Features[edge]; exists {
			return feature
		}

		panic(fmt.Sprint("Edge does not have a corresponding feature mapped. ", edge))
	}

	//should really return nil
	return &DefaultFeature
}

func (t *Tile) FeatureById(id int) *Feature {
	for _, f := range t.Features {
		if f.Id == id {
			return f
		}
	}

	return nil
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

func (t *Tile) EdgeDirsFromFeature(feature *Feature) []directions.Direction {

	dirs := make([]directions.Direction, 0)

	for i, f := range t.Features {
		if f == feature {
			for dir, e := range t.Edges {
				if e == i {
					dirs = append(dirs, dir)
				}
			}
		}
	}

	return dirs
}

func (t *Tile) ComputeEdgeFeatureTypes() [4]FeatureType {

	var ef [4]FeatureType

	for _, d := range directions.List {
		if f := t.ComputeFeature(d); f != nil {
			ef[d] = f.Type
		}
	}

	return ef
}

func (t *Tile) ComputeEdgeFeatures() [4]*Feature {

	var ef [4]*Feature

	for _, d := range directions.List {
		if f := t.ComputeFeature(d); f != nil {
			ef[d] = f
		} else {
			ef[d] = &DefaultFeature
		}
	}

	return ef
}

func (t *Tile) String() string {
	return t.Name
}
