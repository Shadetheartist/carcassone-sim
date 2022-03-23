package game

import (
	"beeb/carcassonne/tile"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

const baseSize int = 8
const scale float64 = 1

func (g *Game) RenderBoard(screen *ebiten.Image) {

	screen.Fill(colornames.Grey800)

	for _, t := range g.Board.Tiles {
		g.renderTileToScreen(t, screen)
	}
}

func (g *Game) renderTileToScreen(t *tile.Tile, screen *ebiten.Image) {

	pos := t.Placement.Position
	op := ebiten.DrawImageOptions{}
	img := ebiten.NewImageFromImage(t.Image)

	d := float64(img.Bounds().Max.X) / 2

	// Move the image's center to the screen's upper-left corner.
	// This is a preparation for rotating. When geometry matrices are applied,
	// the origin point is the upper-left corner.
	op.GeoM.Translate(-d, -d)

	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	op.GeoM.Rotate(float64(t.Placement.Orientation) * math.Pi / 180)

	// Move the image to the screen's center.
	op.GeoM.Translate(d+float64(pos.X*baseSize), d+float64(pos.Y*baseSize))
	op.GeoM.Translate(float64(g.CameraOffset.X), float64(g.CameraOffset.Y))

	op.GeoM.Scale(scale, scale)

	screen.DrawImage(img, &op)
}
