package tile

type Factory struct {
	loadedTiles map[string]Tile
}

func (factory *Factory) Initialize(loadedTiles map[string]Tile) {
	factory.loadedTiles = loadedTiles
}

func (factory *Factory) BuildTile(tileName string) Tile {

	//this will get a new copy of the tile from the map
	tile := factory.loadedTiles[tileName]

	//compute the road segements for the new copy of the tile, otherwise there are pointer issues
	tile.RoadSegments = tile.ComputeRoadSegments()

	return tile
}
