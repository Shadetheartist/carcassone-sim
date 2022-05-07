package tile

import "beeb/carcassonne/matrix"

type TileFactory struct {
}

func (f *TileFactory) NewTileFromReference(rt *ReferenceTile) *Tile {
	tile := &Tile{}

	tile.Reference = rt

	f.rebuildFeaturesFromReference(tile, rt)

	return tile
}

// initializes the tiles feature pointers to different memory for each feature on the reference tile,
// this way, the new tile can modify and reference them as unique memory
// and not affect every tile built from this reference tile
func (f *TileFactory) rebuildFeaturesFromReference(t *Tile, rt *ReferenceTile) {
	t.FeatureMatrix = matrix.NewMatrix[*Feature](rt.FeatureMatrix.Size())
	t.Features = make([]*Feature, len(rt.Features))
	t.EdgeFeatures = EdgeArray[*Feature]{}

	// mapping original features to new features for easy lookup & replacement later
	featureMap := make(map[*Feature]*Feature)

	for i, f := range rt.Features {
		newFeature := &Feature{
			ParentTile:        t,
			ParentRefenceTile: rt,
			Type:              f.Type,
		}

		featureMap[f] = newFeature

		t.Features[i] = newFeature
	}

	//use feature map to easily remap edge features
	for i, f := range rt.EdgeFeatures {
		t.EdgeFeatures[i] = featureMap[f]
	}

	//use feature map to easily remap feature matrix
	rt.FeatureMatrix.Iterate(func(rt *Feature, x int, y int, idx int) {
		t.FeatureMatrix.Set(x, y, featureMap[rt])
	})

}
