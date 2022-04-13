package deck

import (
	"beeb/carcassonne/db"
	"beeb/carcassonne/tile"
)

type TileInfoDeckFactory struct {
	tileInfoFile *db.GameConfig
	tileFactory  *tile.Factory
}

func CreateDeckFactory(gameConfig *db.GameConfig, tileFactory *tile.Factory) *TileInfoDeckFactory {
	df := TileInfoDeckFactory{}

	df.tileInfoFile = gameConfig
	df.tileFactory = tileFactory

	return &df
}

func (df *TileInfoDeckFactory) BuildRiverDeck() Deck {
	deck := Deck{}

	deckSize := deckDataSize(df.tileInfoFile.RiverDeck.Deck)

	deck.Tiles = make([]*tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range df.tileInfoFile.RiverDeck.Deck {
		for i := 0; i < quantity; i++ {
			deck.Tiles[c] = df.tileFactory.BuildTile(tileName)
			c++
		}
	}

	deck.Shuffle()

	deck.Prepend(df.tileFactory.BuildTile(df.tileInfoFile.RiverDeck.Begin))
	deck.Append(df.tileFactory.BuildTile(df.tileInfoFile.RiverDeck.End))

	return deck
}

func (df *TileInfoDeckFactory) BuildDeck() Deck {
	deck := Deck{}

	deckSize := deckDataSize(df.tileInfoFile.Deck)

	deck.Tiles = make([]*tile.Tile, deckSize)

	var c int = 0
	for tileName, quantity := range df.tileInfoFile.Deck {
		for i := 0; i < quantity; i++ {
			deck.Tiles[c] = df.tileFactory.BuildTile(tileName)
			c++
		}
	}

	deck.Shuffle()

	return deck
}

func deckDataSize(deckData map[string]int) int {
	var deckSize int = 0

	for _, v := range deckData {
		deckSize += v
	}

	return deckSize
}
