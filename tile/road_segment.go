package tile

import "beeb/carcassonne/directions"

type RoadSegment struct {
	ParentTile    *Tile
	ParentFeature *Feature

	//direction mapped
	//which edges of a tile are mapped to this road
	EdgeSegments [4]*RoadSegment
}

func (t *Tile) IntegrateRoads() {

	uniqueRs := t.UniqueRoadSegements()
	for _, rs := range uniqueRs {

		if rs == nil {
			continue
		}

		for _, d := range rs.ParentFeature.Edges {
			nd := t.Placement.TileToGridDir(d)
			neighbour := t.Neighbours[nd]

			if neighbour != nil {
				cd := directions.Compliment[directions.Direction(nd)]
				tcd := neighbour.Placement.GridToTileDir(cd)
				connectedRoadSegment := neighbour.RoadSegments[tcd]
				if connectedRoadSegment != nil {
					connectedRoadSegment.EdgeSegments[tcd] = rs
					rs.EdgeSegments[d] = connectedRoadSegment
				}
			}
		}
	}

}

func (t *Tile) CreateRoadSegment(f *Feature) *RoadSegment {

	if f.Type != FeatureTypeRoad {
		panic("Create road segments only for road features")
	}

	var es [4]*RoadSegment

	rs := RoadSegment{
		ParentTile:    t,
		ParentFeature: f,
		EdgeSegments:  es,
	}

	return &rs
}

func (t *Tile) UniqueRoadSegements() []*RoadSegment {

	ursMap := make(map[*RoadSegment]struct{})

	//using a pointer to the same road segment intentionally here
	for _, rs := range t.RoadSegments {
		if rs != nil {
			ursMap[rs] = struct{}{}
		}
	}

	keys := make([]*RoadSegment, 0, len(ursMap))
	for rs := range ursMap {
		keys = append(keys, rs)
	}

	return keys
}

//must be called after the edge features are set up as it relies on them
func ComputeRoadSegments(t *Tile) [4]*RoadSegment {

	var rs [4]*RoadSegment

	//using a pointer to the same road segment intentionally here
	for _, f := range t.FeaturesByType(FeatureTypeRoad) {
		s := t.CreateRoadSegment(f)
		for _, d := range f.Edges {
			rs[d] = s
		}
	}

	return rs
}
