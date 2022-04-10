package db

//types required for unmarshalling the tile info file data
type GameConfig struct {
	Tiles     map[string]TileInfo
	Deck      map[string]int
	RiverDeck RiverDeck
}

type RiverDeck struct {
	Begin string
	End   string
	Deck  map[string]int
}

type GameConfigLoader interface {
	GetGameConfig() GameConfig
	GetAllTileNames() []string
}
