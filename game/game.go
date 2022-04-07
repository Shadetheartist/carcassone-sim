package game

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/loader"
	"beeb/carcassonne/tile"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	TileInfo    loader.TileInfoFile
	Tiles       map[string]tile.Tile
	TileFactory tile.Factory

	RiverDeck    Deck
	Deck         Deck
	Board        board.Board
	CameraOffset image.Point
	CameraZoom   float64

	HoveredPosition  tile.Position
	SelectedPosition tile.Position
	HighlightedRoads []board.Road

	//the size of the bitmap images, add 1 to show a grid
	baseSize    int
	renderScale float64

	//these are only relevant to creating the river
	//1 is not ever a valid orientation so it will not be a false positive
	lastRiverTurn uint16
	lastRiverTile *tile.Tile

	ImageW int
	ImageH int
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.RenderBoard(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 400, 400
}
