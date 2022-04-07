package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/loader"
	"beeb/carcassonne/tile"
	"os"
	"path/filepath"
)

func (g *Game) Initialize() {

	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)

	ymlPath := filepath.Join(exeDir, "data/tiles.yml")
	bitmapDirectory := filepath.Join(exeDir, "data/bitmaps")

	g.ImageW = 1000
	g.ImageH = 1000

	g.Tiles, g.TileInfo = loader.LoadTiles(ymlPath, bitmapDirectory)
	g.Board = board.New(g.Tiles, g.ImageW, g.ImageH)
	g.RiverDeck = g.buildRiverDeck()
	g.Deck = g.buildDeck()
	g.baseSize = 7
	g.renderScale = 2

	g.lastRiverTurn = 1
	g.lastRiverTile = nil

	g.CameraOffset.X = g.ImageW / g.baseSize
	g.CameraOffset.Y = g.ImageH / g.baseSize

	g.initializeRenderer()

	g.HighlightedRoads = make([]board.Road, 0)
}

//how many tiles are in the deck data loaded from the yml file
func deckDataSize(deckData map[string]int) int {
	var deckSize int = 0

	for _, v := range deckData {
		deckSize += v
	}

	return deckSize
}

func (g *Game) buildRiverDeck() Deck {

	deck := Deck{}

	deckSize := deckDataSize(g.TileInfo.RiverDeck.Deck)

	deck.Tiles = make([]tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range g.TileInfo.RiverDeck.Deck {
		for i := 0; i < quantity; i++ {
			tile := g.Tiles[tileName]

			//compute the road segements for the copy of the tile that is in the deck, not the tile from the loader
			//otherwise we will have pointer issues
			tile.RoadSegments = tile.ComputeRoadSegments()

			deck.Tiles[c] = tile

			c++
		}
	}

	deck.Shuffle()

	deck.Prepend(g.Tiles[g.TileInfo.RiverDeck.Begin])

	deck.Append(g.Tiles[g.TileInfo.RiverDeck.End])

	return deck
}

func (g *Game) buildDeck() Deck {
	deck := Deck{}

	deckSize := deckDataSize(g.TileInfo.Deck)

	deck.Tiles = make([]tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range g.TileInfo.Deck {
		for i := 0; i < quantity; i++ {
			tile := g.Tiles[tileName]

			//compute the road segements for the copy of the tile that is in the deck, not the tile from the loader
			//otherwise we will have pointer issues
			tile.RoadSegments = tile.ComputeRoadSegments()

			deck.Tiles[c] = tile

			c++
		}
	}

	deck.Shuffle()

	return deck
}
