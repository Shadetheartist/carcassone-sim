package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/db"
	"beeb/carcassonne/game/deck"
	"beeb/carcassonne/tile"
	"os"
	"path/filepath"
)

func (g *Game) Initialize(gcl db.GameConfigLoader, til db.TileInfoLoader, bmpl db.BitmapLoader) {

	g.TileFactory = &tile.Factory{}
	g.TileFactory.Initialize(gcl.GetAllTileNames(), til, bmpl)

	deckFactory := deck.TileInfoDeckFactory{}
	tileInfoFile := gcl.GetGameConfig()
	deckFactory.Initialize(&tileInfoFile, g.TileFactory)

	g.ImageW = 1000
	g.ImageH = 1000

	g.Board = board.New(g.Tiles, g.ImageW, g.ImageH)
	g.RiverDeck = deckFactory.BuildRiverDeck()
	g.Deck = deckFactory.BuildDeck()
	g.baseSize = 7
	g.renderScale = 2

	g.lastRiverTurn = 1
	g.lastRiverTile = nil

	g.CameraOffset.X = g.ImageW / g.baseSize
	g.CameraOffset.Y = g.ImageH / g.baseSize

	g.initializeRenderer()

	g.HighlightedRoads = make([]board.Road, 0)
}

func ymlConfigFilePath() string {
	return filepath.Join(exeDir(), "data/tiles.yml")
}

func bitmapDirectory() string {
	return filepath.Join(exeDir(), "data/bitmaps")
}

func exeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exePath)
}
