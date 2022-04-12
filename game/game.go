package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/db"
	"beeb/carcassonne/game/deck"
	"beeb/carcassonne/tile"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	TileInfo    db.GameConfig
	Tiles       map[string]tile.Tile
	TileFactory *tile.Factory

	RiverDeck    deck.Deck
	Deck         deck.Deck
	Board        board.Board
	CameraOffset image.Point
	CameraZoom   float64

	HoveredPosition  tile.Position
	SelectedPosition tile.Position
	HighlightedRoads []board.Road

	//these are only relevant to creating the river
	//1 is not ever a valid orientation so it will not be a false positive
	lastRiverTurn uint16
	lastRiverTile *tile.Tile
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.RenderBoard(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400
}
