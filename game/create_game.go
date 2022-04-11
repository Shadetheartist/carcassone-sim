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

	g.ImageW = 1000
	g.ImageH = 1000

	g.baseSize = 7
	g.renderScale = 2

	g.CameraOffset.X = g.ImageW / g.baseSize
	g.CameraOffset.Y = g.ImageH / g.baseSize

	g.HighlightedRoads = make([]board.Road, 0)

	g.TileFactory = tile.CreateTileFactory(gcl.GetAllTileNames(), til, bmpl)

	gameConfig := gcl.GetGameConfig()
	deckFactory := deck.CreateDeckFactory(&gameConfig, g.TileFactory)

	g.Board = board.CreateBoard(g.Tiles, g.ImageW, g.ImageH)
	g.RiverDeck = deckFactory.BuildRiverDeck()
	g.Deck = deckFactory.BuildDeck()

	g.initializeRenderer()

	return &g
}
