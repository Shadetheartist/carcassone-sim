package game

import (
	"beeb/carcassonne/board/road"
	"beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func (g *Game) initializeRenderer() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.initializeHoverImage()

}

var hoverImage *ebiten.Image

func (g *Game) initializeHoverImage() {
	ctx := gg.NewContext(g.baseSize, g.baseSize)

	ctx.DrawRectangle(0, 0, float64(g.baseSize), float64(g.baseSize))
	ctx.SetRGBA(0, 0, 0, 0.25)
	ctx.Fill()

	ctx.SetRGBA(1, 1, 1, 1)
	ctx.DrawString(fmt.Sprint("X:", g.HoveredPosition.X), 0, 0)
	ctx.DrawString(fmt.Sprint("Y:", g.HoveredPosition.Y), 8, 0)
	ctx.Stroke()

	img := ctx.Image()
	hoverImage = ebiten.NewImageFromImage(img)
}

func (g *Game) RenderBoard(screen *ebiten.Image) {

	screen.Fill(colornames.Grey600)

	for _, t := range g.Board.Tiles {
		g.renderTileToBoard(t)
	}

	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		float64(g.CameraOffset.X),
		float64(g.CameraOffset.Y),
	)

	g.renderRoads()

	screen.DrawImage(g.Board.BoardImage, &op)
	screen.DrawImage(g.Board.RoadsImage, &op)
	g.renderHoveredPosition(screen)

}

func (g *Game) renderTileToBoard(t *tile.Tile) {

	if t.Rendered {
		return
	}

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
	op.GeoM.Translate(d+float64(pos.X*g.baseSize), d+float64(pos.Y*g.baseSize))

	op.GeoM.Scale(g.renderScale, g.renderScale)

	g.Board.BoardImage.DrawImage(img, &op)

	t.Rendered = true
}

func (g *Game) renderHoveredPosition(screen *ebiten.Image) {

	text.Draw(screen, fmt.Sprint("X:", g.HoveredPosition.X), mplusNormalFont, 20, 20, color.White)
	text.Draw(screen, fmt.Sprint("Y:", g.HoveredPosition.Y), mplusNormalFont, 20, 60, color.White)

	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		float64(g.HoveredPosition.X)*float64(g.baseSize),
		float64(g.HoveredPosition.Y)*float64(g.baseSize),
	)

	op.GeoM.Scale(
		g.renderScale,
		g.renderScale,
	)

	op.GeoM.Translate(
		float64(g.CameraOffset.X),
		float64(g.CameraOffset.Y),
	)

	screen.DrawImage(hoverImage, &op)
}

func (g *Game) renderRoadsByTile() {
	ctx := gg.NewContext(1000, 1000)

	for _, t := range g.Board.Tiles {

		x := float64(t.Placement.Position.X*g.baseSize) + float64(g.baseSize)/2
		y := float64(t.Placement.Position.Y*g.baseSize) + float64(g.baseSize)/2

		roads := t.FeaturesByType(tile.Road)

		for _, r := range roads {
			dirs := t.EdgeDirsFromFeature(r)

			for _, d := range dirs {
				tileDir := t.Placement.TileToGridDir(d)
				offset := dirOffset(tileDir)

				offsetX := float64(offset.X*g.baseSize) / 2
				offsetY := float64(offset.Y*g.baseSize) / 2

				//to the center of the edge of the feature
				ctx.LineTo(x+offsetX, y+offsetY)
			}
		}

		ctx.SetLineWidth(1)
		ctx.SetRGBA(0, 0, 1, 1)
		ctx.Stroke()
	}

	roadsImage := ctx.Image()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)

	ebitenImage := ebiten.NewImageFromImage(roadsImage)
	g.Board.RoadsImage.Clear()
	g.Board.RoadsImage.DrawImage(ebitenImage, &op)
}

func (g *Game) renderRoads() {
	ctx := gg.NewContext(1000, 1000)

	for id, r := range g.Board.Roads {

		for next := r.First(); next != nil; next = next.Next() {
			g.renderRoad(ctx, id, next)
		}

		ctx.SetLineWidth(1)
		ctx.SetRGBA(0, 0, 1, 1)
		ctx.Stroke()
	}

	roadsImage := ctx.Image()

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(g.renderScale, g.renderScale)

	ebitenImage := ebiten.NewImageFromImage(roadsImage)
	g.Board.RoadsImage.Clear()
	g.Board.RoadsImage.DrawImage(ebitenImage, &op)
}

func dirOffset(dir directions.Direction) image.Point {
	switch dir {
	case directions.North:
		return image.Point{X: 0, Y: -1}
	case directions.East:
		return image.Point{X: 1, Y: 0}
	case directions.South:
		return image.Point{X: 0, Y: 1}
	case directions.West:
		return image.Point{X: -1, Y: 0}
	}

	panic("Invalid Direction Supplied")
}

func (g *Game) renderRoad(ctx *gg.Context, roadId int, node *road.Node) {
	tile := node.Value

	x := float64(tile.Placement.Position.X*g.baseSize) + float64(g.baseSize)/2
	y := float64(tile.Placement.Position.Y*g.baseSize) + float64(g.baseSize)/2

	feature := tile.FeatureById(roadId)

	if feature == nil {
		return
		//panic("Edge Feature Not found by Id")
	}

	dirs := tile.EdgeDirsFromFeature(feature)

	for _, d := range dirs {
		tileDir := tile.Placement.TileToGridDir(d)
		offset := dirOffset(tileDir)

		offsetX := float64(offset.X*g.baseSize) / 2
		offsetY := float64(offset.Y*g.baseSize) / 2

		//to the center of the edge of the feature
		ctx.LineTo(x+offsetX, y+offsetY)
	}
}
