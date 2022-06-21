package tile

import (
	"github.com/google/uuid"
)

type TileFactory struct {
}

func (f *TileFactory) NewTileFromReference(rt *ReferenceTile) *Tile {

	if rt == nil {
		panic("Dont send this a null reference tile")
	}

	tile := &Tile{}
	tile.Id = uuid.Nil
	tile.Reference = rt

	f.rebuildFeaturesFromReference(tile, rt)

	return tile
}

// initializes the tiles feature pointers to different memory for each feature on the reference tile,
// this way, the new tile can modify and reference them as unique memory
// and not affect every tile built from this reference tile
func (f *TileFactory) rebuildFeaturesFromReference(t *Tile, rt *ReferenceTile) {
	t.Features = make([]*Feature, len(rt.Features))
	t.EdgeFeatures = &EdgeArray[*Feature]{}

	// mapping original features to new features for easy lookup & replacement later
	t.ReferenceFeatureMap = make(map[*Feature]*Feature)

	for i, f := range rt.Features {
		newFeature := &Feature{
			Id:            uuid.Nil,
			ParentTile:    t,
			ParentFeature: f,
			Type:          f.Type,
		}

		newFeature.Links = make(map[*Feature]*Feature)

		t.ReferenceFeatureMap[f] = newFeature

		t.Features[i] = newFeature
	}

	//use feature map to easily remap edge features
	for i, f := range rt.EdgeFeatures {
		t.EdgeFeatures[i] = t.ReferenceFeatureMap[f]
	}

	t.Neighbours = &EdgeArray[*Tile]{}

}
