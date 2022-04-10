package tile

import (
	"beeb/carcassonne/db"
	"beeb/carcassonne/directions"
)

type Factory struct {
	tileInfoLoader db.TileInfoLoader
	bitmapLoader   db.BitmapLoader
	tileNames      []string

	referenceTiles map[string]Tile
}

func (factory *Factory) Initialize(tileNames []string, tileInfoLoader db.TileInfoLoader, bitmapLoader db.BitmapLoader) {
	factory.tileNames = tileNames
	factory.tileInfoLoader = tileInfoLoader
	factory.bitmapLoader = bitmapLoader

	factory.referenceTiles = make(map[string]Tile)
	for _, tileName := range factory.tileNames {
		factory.referenceTiles[tileName] = factory.BuildTile(tileName)
	}
}

func (factory *Factory) ReferenceTiles() map[string]Tile {
	return factory.referenceTiles
}

func (factory *Factory) GetReferenceTile(tileName string) Tile {
	return factory.referenceTiles[tileName]
}

func (factory *Factory) BuildTile(tileName string) Tile {

	//get tile info and bitmap from loaders
	tileInfo, err := factory.tileInfoLoader.GetTileInfo(tileName)

	if err != nil {
		panic(err)
	}

	tileBitmap, err := factory.bitmapLoader.GetTileBitmap(tileInfo.Image)

	if err != nil {
		panic(err)
	}

	edges := make(map[directions.Direction]int)

	for edgeStr, featureId := range tileInfo.Edges {
		edges[directions.StrMap[edgeStr]] = featureId
	}

	features := make(map[int]*Feature)

	for featureId, featureInfo := range tileInfo.Features {

		edgesForFeature := make([]directions.Direction, 0)

		for dir, fid := range edges {
			if fid == featureId {
				edgesForFeature = append(edgesForFeature, dir)
			}
		}

		features[featureId] = &Feature{
			Type:   FeatureTypeStrMap[featureInfo.Type],
			Shield: featureInfo.Shield,
			Curve:  featureInfo.Curve,
			Edges:  edgesForFeature,
		}
	}

	t := Tile{
		Name:     tileName,
		Image:    tileBitmap,
		Features: features,
		Edges:    edges,
		Placement: Placement{
			Position:    Position{},
			Orientation: 0,
		},
	}

	//recompute pointers to new memory
	t.Neighbours = make([]*Tile, 4)
	t.EdgeFeatures = t.ComputeEdgeFeatures()
	t.EdgeFeatureTypes = t.ComputeEdgeFeatureTypes()
	t.RoadSegments = t.ComputeRoadSegments()

	return t
}
