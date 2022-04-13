package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/db"
	"beeb/carcassonne/game/deck"
	"beeb/carcassonne/tile"
)

func CreateGame(gcl db.GameConfigLoader, til db.TileInfoLoader, bmpl db.BitmapLoader) *Game {
	g := Game{}

	g.lastRiverTurn = 1
	g.lastRiverTile = nil

	g.CameraOffset.X = 0
	g.CameraOffset.Y = 0

	g.HighlightedRoads = make([]board.Road, 0)

	g.TileFactory = tile.CreateTileFactory(gcl.GetAllTileNames(), til, bmpl)

	gameConfig := gcl.GetGameConfig()
	deckFactory := deck.CreateDeckFactory(&gameConfig, g.TileFactory)

	g.Board = board.CreateBoard(g.TileFactory.ReferenceTiles(), 1000, 1000)
	g.RiverDeck = deckFactory.BuildRiverDeck()
	g.Deck = deckFactory.BuildDeck()

	g.initializeRenderer()

	return &g
}
