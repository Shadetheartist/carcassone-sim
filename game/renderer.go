package game

import (
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

func (g *Game) RenderBoard(screen *ebiten.Image) {

	screen.Fill(colornames.Grey600)

	for _, t := range g.Board.Tiles {
		g.renderTileToBoard(t)
	}

	g.Board.OpenPositionsImage.Clear()

	for p, fts := range g.Board.OpenPositions {
		g.renderOpenPosition(p, fts)
	}

	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		float64(g.CameraOffset.X),
		float64(g.CameraOffset.Y),
	)

	screen.DrawImage(g.Board.BoardImage, &op)
	screen.DrawImage(g.Board.OpenPositionsImage, &op)
	screen.DrawImage(g.Board.RoadsImage, &op)

	text.Draw(screen, fmt.Sprint("X:", g.HoveredPosition.X), mplusNormalFont, 20, 20, color.White)
	text.Draw(screen, fmt.Sprint("Y:", g.HoveredPosition.Y), mplusNormalFont, 20, 60, color.White)

	g.renderHoveredPosition(screen)
}

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
		Size:    14,
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
	g.initializeOpenPositionImage()

}

func (g *Game) isRoadSegmentHighlighted(rs *tile.RoadSegment) bool {
	for _, rd := range g.HighlightedRoads {
		if rd.ContainsSegment(rs) {
			return true
		}
	}

	return false
}

func (g *Game) setRoadSegmentColor(ctx *gg.Context, rs *tile.RoadSegment) {
	if g.isRoadSegmentHighlighted(rs) {
		ctx.SetRGBA(0, 1, 1, 1)

	} else {
		ctx.SetRGBA(0, 0, 0, 1)
	}
}

func (g *Game) renderRoadSegment(img *ebiten.Image, rs *tile.RoadSegment) {

	if rs == nil {
		return
	}

	ctx := gg.NewContext(g.baseSize, g.baseSize)

	ep := tile.Position{}

	//middle to edge
	if len(rs.ParentFeature.Edges) == 1 {
		edge := rs.ParentFeature.Edges[0]

		edgePos := ep.EdgePos(edge)

		g.setRoadSegmentColor(ctx, rs)

		ctx.MoveTo(float64(g.baseSize)/2, float64(g.baseSize)/2)
		ctx.LineTo(
			float64(g.baseSize)/2+float64(edgePos.X)*float64(g.baseSize)/2,
			float64(g.baseSize)/2+float64(edgePos.Y)*float64(g.baseSize)/2,
		)
		ctx.SetLineWidth(0.5)
		ctx.Stroke()

		ctx.SetRGBA(1, 1, 1, 1)
		ctx.DrawCircle(float64(g.baseSize)/2, float64(g.baseSize)/2, 1.5)
		ctx.Fill()

	} else if len(rs.ParentFeature.Edges) == 2 {
		//edge to edge

		g.setRoadSegmentColor(ctx, rs)

		edgeA := rs.ParentFeature.Edges[0]
		edgeB := rs.ParentFeature.Edges[1]
		edgeAPos := ep.EdgePos(edgeA)
		edgeBPos := ep.EdgePos(edgeB)

		ctx.MoveTo(
			float64(g.baseSize)/2+float64(edgeAPos.X)*float64(g.baseSize)/2,
			float64(g.baseSize)/2+float64(edgeAPos.Y)*float64(g.baseSize)/2,
		)

		ctx.LineTo(
			float64(g.baseSize)/2+float64(edgeBPos.X)*float64(g.baseSize)/2,
			float64(g.baseSize)/2+float64(edgeBPos.Y)*float64(g.baseSize)/2,
		)

		ctx.Stroke()
	}

	ebiImg := ebiten.NewImageFromImage(ctx.Image())

	op := ebiten.DrawImageOptions{}
	img.DrawImage(ebiImg, &op)

}

var hoverImage *ebiten.Image

func (g *Game) initializeHoverImage() {
	ctx := gg.NewContext(g.baseSize, g.baseSize)

	ctx.DrawRectangle(0, 0, float64(g.baseSize), float64(g.baseSize))
	ctx.SetRGBA(0, 0, 0, 0.25)
	ctx.Fill()

	img := ctx.Image()
	hoverImage = ebiten.NewImageFromImage(img)
}

var openPosImage *ebiten.Image

func (g *Game) initializeOpenPositionImage() {
	ctx := gg.NewContext(g.baseSize, g.baseSize)

	ctx.DrawRectangle(0, 0, float64(g.baseSize), float64(g.baseSize))
	ctx.SetRGBA(0, 1, 0, 0.25)
	ctx.Fill()

	img := ctx.Image()
	openPosImage = ebiten.NewImageFromImage(img)
}

func (g *Game) renderTileToBoard(t *tile.Tile) {

	if t.Rendered {
		//return
	}

	pos := t.Placement.Position
	op := ebiten.DrawImageOptions{}
	img := ebiten.NewImageFromImage(t.Image)

	d := float64(img.Bounds().Max.X) / 2

	for _, rs := range t.UniqueRoadSegements() {
		g.renderRoadSegment(img, rs)
	}

	// Move the image's center to the screen's upper-left corner.
	// This is a preparation for rotating. When geometry matrices are applied,
	// the origin point is the upper-left corner.
	op.GeoM.Translate(-d, -d)

	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	op.GeoM.Rotate(float64(t.Placement.Orientation) * math.Pi / 180)

	op.GeoM.Translate(d+float64(pos.X*g.baseSize), d+float64(pos.Y*g.baseSize))

	op.GeoM.Scale(g.renderScale, g.renderScale)

	g.Board.BoardImage.DrawImage(img, &op)

	t.Rendered = true
}

func (g *Game) renderOpenPosition(p tile.Position, fts []tile.FeatureType) {
	op := ebiten.DrawImageOptions{}

	// Move the image to the screen's center.
	op.GeoM.Translate(float64(p.X*g.baseSize), float64(p.Y*g.baseSize))

	op.GeoM.Scale(g.renderScale, g.renderScale)

	g.Board.OpenPositionsImage.DrawImage(hoverImage, &op)
}

func (g *Game) renderHoveredPosition(screen *ebiten.Image) {

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
