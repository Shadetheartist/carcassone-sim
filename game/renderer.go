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
	g.Board.RoadsImage.Clear()

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

	g.Board.FarmsImage.Clear()
	g.renderFarms()
	screen.DrawImage(g.Board.FarmsImage, &op)

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

func (g *Game) toScreenSpace(p image.Point) image.Point {
	return image.Point{
		X: p.X * int(g.Board.RenderScale),
		Y: p.Y * int(g.Board.RenderScale),
	}
}

func sumPoint(p1 image.Point, p2 image.Point) image.Point {
	return image.Point{
		X: p1.X + p2.X,
		Y: p1.Y + p2.Y,
	}
}

func (g *Game) renderFarms() {
	ctx := gg.NewContext(g.Board.FarmsImage.Size())
	ctx.SetRGBA(0, 0, 0, 1)

	for pos, t := range g.Board.Tiles {
		pix := g.Board.PosToPix(pos)

		for _, fs := range t.FarmSegments {
			matrix := tile.OrientedFarmMatrix(t, int(t.Placement.Orientation))
			avgPoint := tile.AvgFarmSegmentPos(fs, matrix)
			ssAvgPoint := sumPoint(pix, g.toScreenSpace(avgPoint))

			ctx.SetPixel(ssAvgPoint.X, ssAvgPoint.Y)

			for _, neighbour := range fs.Neighbours {
				nPos := neighbour.Parent.Placement.Position
				nPix := g.Board.PosToPix(nPos)
				nMatrix := tile.OrientedFarmMatrix(neighbour.Parent, int(neighbour.Parent.Placement.Orientation))
				nAvgPoint := tile.AvgFarmSegmentPos(neighbour, nMatrix)
				nSsAvgPoint := sumPoint(nPix, g.toScreenSpace(nAvgPoint))
				ctx.MoveTo(float64(ssAvgPoint.X), float64(ssAvgPoint.Y))
				ctx.LineTo(float64(nSsAvgPoint.X), float64(nSsAvgPoint.Y))
				ctx.Stroke()
			}
		}

	}

	ebiImg := ebiten.NewImageFromImage(ctx.Image())

	op := ebiten.DrawImageOptions{}
	g.Board.FarmsImage.DrawImage(ebiImg, &op)
}

func (g *Game) renderFarmSegment(fs *tile.FarmSegment) {

	//image size needs to encompass the neighbours of this tile
	ctx := gg.NewContext(g.Board.BaseSize*int(g.Board.RenderScale)*3, g.Board.BaseSize*int(g.Board.RenderScale)*3)
	ctx.SetRGBA(0, 0, 0, 1)

	offset := g.Board.BaseSize * int(g.Board.RenderScale)

	matrix := tile.OrientedFarmMatrix(fs.Parent, int(fs.Parent.Placement.Orientation))
	avgPoint := tile.AvgFarmSegmentPos(fs, matrix)

	tilePosition := fs.Parent.Placement.Position
	ctx.SetPixel(int(float64(avgPoint.X)*g.Board.RenderScale)+offset, int(float64(avgPoint.Y)*g.Board.RenderScale)+offset)

	ctx.SetRGBA(1, 0, 0, 1)

	for _, neighbour := range fs.Neighbours {
		nMatrix := tile.OrientedFarmMatrix(neighbour.Parent, int(neighbour.Parent.Placement.Orientation))
		nAvgPoint := tile.AvgFarmSegmentPos(neighbour, nMatrix)

		if nAvgPoint.X == -1 || nAvgPoint.Y == -1 {
			continue
		}

		neighbourPosition := neighbour.Parent.Placement.Position
		offsetPosition := tile.Position{
			X: neighbourPosition.X - tilePosition.X,
			Y: neighbourPosition.Y - tilePosition.Y,
		}

		neighbourPix := image.Point{
			X: (nAvgPoint.X + int(float64(offsetPosition.X*g.Board.BaseSize))*int(g.Board.RenderScale)),
			Y: (nAvgPoint.Y + int(float64(offsetPosition.Y*g.Board.BaseSize))*int(g.Board.RenderScale)),
		}

		ctx.MoveTo(float64(avgPoint.X)*g.Board.RenderScale+float64(offset), float64(avgPoint.Y)*g.Board.RenderScale+float64(offset))
		ctx.LineTo(float64(neighbourPix.X)+float64(offset), float64(neighbourPix.Y)+float64(offset))
		ctx.Stroke()
	}

	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		float64(tilePosition.X*g.Board.BaseSize*int(g.Board.RenderScale)-offset),
		float64(tilePosition.Y*g.Board.BaseSize*int(g.Board.RenderScale)-offset),
	)

	ebiImg := ebiten.NewImageFromImage(ctx.Image())
	g.Board.FarmsImage.DrawImage(ebiImg, &op)
}

func (g *Game) renderRoadSegment(op ebiten.DrawImageOptions, rs *tile.RoadSegment) {

	if rs == nil {
		return
	}

	if g.isRoadSegmentHighlighted(rs) == false {
		return
	}

	ctx := gg.NewContext(g.Board.BaseSize, g.Board.BaseSize)

	ctx.SetRGBA(1, 0, 0, 1)

	ep := tile.Position{}

	//middle to edge
	if len(rs.ParentFeature.Edges) == 1 {
		edge := rs.ParentFeature.Edges[0]

		edgePos := ep.EdgePos(edge)

		g.setRoadSegmentColor(ctx, rs)

		ctx.MoveTo(float64(g.Board.BaseSize)/2, float64(g.Board.BaseSize)/2)
		ctx.LineTo(
			float64(g.Board.BaseSize)/2+float64(edgePos.X)*float64(g.Board.BaseSize)/2,
			float64(g.Board.BaseSize)/2+float64(edgePos.Y)*float64(g.Board.BaseSize)/2,
		)

		ctx.Stroke()

	} else if len(rs.ParentFeature.Edges) == 2 {
		//edge to edge

		g.setRoadSegmentColor(ctx, rs)

		edgeA := rs.ParentFeature.Edges[0]
		edgeB := rs.ParentFeature.Edges[1]
		edgeAPos := ep.EdgePos(edgeA)
		edgeBPos := ep.EdgePos(edgeB)

		ctx.MoveTo(
			float64(g.Board.BaseSize)/2+float64(edgeAPos.X)*float64(g.Board.BaseSize)/2,
			float64(g.Board.BaseSize)/2+float64(edgeAPos.Y)*float64(g.Board.BaseSize)/2,
		)

		ctx.LineTo(
			float64(g.Board.BaseSize)/2+float64(edgeBPos.X)*float64(g.Board.BaseSize)/2,
			float64(g.Board.BaseSize)/2+float64(edgeBPos.Y)*float64(g.Board.BaseSize)/2,
		)

		ctx.Stroke()
	}

	ebiImg := ebiten.NewImageFromImage(ctx.Image())

	g.Board.RoadsImage.DrawImage(ebiImg, &op)

}

var hoverImage *ebiten.Image

func (g *Game) initializeHoverImage() {
	ctx := gg.NewContext(g.Board.BaseSize, g.Board.BaseSize)

	ctx.DrawRectangle(0, 0, float64(g.Board.BaseSize), float64(g.Board.BaseSize))
	ctx.SetRGBA(0, 0, 0, 0.25)
	ctx.Fill()

	img := ctx.Image()
	hoverImage = ebiten.NewImageFromImage(img)
}

var openPosImage *ebiten.Image

func (g *Game) initializeOpenPositionImage() {
	ctx := gg.NewContext(g.Board.BaseSize, g.Board.BaseSize)

	ctx.DrawRectangle(0, 0, float64(g.Board.BaseSize), float64(g.Board.BaseSize))
	ctx.SetRGBA(0, 1, 0, 0.25)
	ctx.Fill()

	img := ctx.Image()
	openPosImage = ebiten.NewImageFromImage(img)
}

func (g *Game) renderTileToBoard(t *tile.Tile) {

	pos := t.Placement.Position

	op := ebiten.DrawImageOptions{}

	d := float64(g.Board.BaseSize) / 2

	// Move the image's center to the screen's upper-left corner.
	// This is a preparation for rotating. When geometry matrices are applied,
	// the origin point is the upper-left corner.
	op.GeoM.Translate(-d, -d)

	// Rotate the image. As a result, the anchor point of this rotate is
	// the center of the image.
	op.GeoM.Rotate(float64(t.Placement.Orientation) * math.Pi / 180)

	op.GeoM.Translate(d+float64(pos.X*g.Board.BaseSize), d+float64(pos.Y*g.Board.BaseSize))

	op.GeoM.Scale(g.Board.RenderScale, g.Board.RenderScale)

	for _, rs := range t.UniqueRoadSegements() {
		g.renderRoadSegment(op, rs)
	}

	if t.Rendered {
		return
	}

	img := ebiten.NewImageFromImage(t.Image)
	g.Board.BoardImage.DrawImage(img, &op)

	t.Rendered = true
}

func (g *Game) renderOpenPosition(p tile.Position, fts []tile.FeatureType) {
	op := ebiten.DrawImageOptions{}

	// Move the image to the screen's center.
	op.GeoM.Translate(float64(p.X*g.Board.BaseSize), float64(p.Y*g.Board.BaseSize))

	op.GeoM.Scale(g.Board.RenderScale, g.Board.RenderScale)

	g.Board.OpenPositionsImage.DrawImage(hoverImage, &op)
}

func (g *Game) renderHoveredPosition(screen *ebiten.Image) {

	op := ebiten.DrawImageOptions{}

	op.GeoM.Translate(
		float64(g.HoveredPosition.X)*float64(g.Board.BaseSize),
		float64(g.HoveredPosition.Y)*float64(g.Board.BaseSize),
	)

	op.GeoM.Scale(
		g.Board.RenderScale,
		g.Board.RenderScale,
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
