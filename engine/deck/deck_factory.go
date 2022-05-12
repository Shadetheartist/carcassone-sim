package deck

import (
	"beeb/carcassonne/data"
	"beeb/carcassonne/tile"
)

func BuildRiverDeck(gameData *data.GameData) *Deck {
	riverTiles := make([]*tile.ReferenceTileGroup, 0)

	for _, referenceTileGroup := range gameData.ReferenceTileGroups {
		if referenceTileGroup.IsRiverTerminus() {
			continue
		}

		if referenceTileGroup.IsRiverTile() {
			riverTiles = append(riverTiles, referenceTileGroup)
		}
	}

	deck := &Deck{}

	for _, riverTileGroup := range riverTiles {
		tileName := riverTileGroup.Name
		tileCount := gameData.DeckInfo.Deck[tileName]
		for i := 0; i < tileCount; i++ {
			deck.Append(riverTileGroup)
		}
	}

	terminus := gameData.ReferenceTileGroups["RiverTerminus"]

	deck.Shuffle()

	deck.Prepend(terminus)
	deck.Append(terminus)

	return deck
}

func BuildDeck(gameData *data.GameData) *Deck {
	nonRiverTiles := make([]*tile.ReferenceTileGroup, 0)

	for _, referenceTileGroup := range gameData.ReferenceTileGroups {
		if !referenceTileGroup.IsRiverTile() {
			nonRiverTiles = append(nonRiverTiles, referenceTileGroup)
		}
	}

	deck := &Deck{}

	for _, nonRiverTileGroup := range nonRiverTiles {
		tileName := nonRiverTileGroup.Name
		tileCount := gameData.DeckInfo.Deck[tileName]
		for i := 0; i < tileCount; i++ {
			deck.Append(nonRiverTileGroup)
		}
	}

	deck.Shuffle()

	return deck
}
