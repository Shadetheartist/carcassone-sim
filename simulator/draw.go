package simulator

import (
	"beeb/carcassonne/engine/deck"
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/util"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const TILE_SIZE int = 7

type DrawData struct {
	font font.Face

	emptyTileImage     *ebiten.Image
	lightTileBackImage *ebiten.Image
	darkTileBackImage  *ebiten.Image

	boardImage                  *ebiten.Image
	possibleTilePlacementsImage *ebiten.Image
	overlayImg                  *ebiten.Image

	riverDeckImage         *ebiten.Image
	deckImage              *ebiten.Image
	currentlyHeldTileImage *ebiten.Image

	scale        float64
	boardScale   float64
	uiScale      float64
	hdScale      int
	cameraOffset util.Point[int]

	redrawBoard bool

	blackShader          *ebiten.Shader
	transparentRedShader *ebiten.Shader
	colorShader          *ebiten.Shader
}

func (sim *Simulator) initDraw() {
	sim.initFont()
	sim.drawData.scale = 2
	sim.drawData.uiScale = 4
	sim.drawData.boardScale = 1
	sim.drawData.hdScale = 8
	sim.drawData.cameraOffset = util.Point[int]{}

	boardPxSize := sim.Engine.GameBoard.TileMatrix.Size() * TILE_SIZE
	sim.drawData.boardImage = ebiten.NewImage(boardPxSize, boardPxSize)
	sim.drawData.possibleTilePlacementsImage = ebiten.NewImage(boardPxSize, boardPxSize)
	sim.drawData.overlayImg = ebiten.NewImage(boardPxSize*sim.drawData.hdScale, boardPxSize*sim.drawData.hdScale)

	sim.drawData.deckImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	sim.drawData.riverDeckImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	sim.drawData.currentlyHeldTileImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)

	sim.drawData.emptyTileImage = loadImage("./simulator/images/empty_tile.bmp")
	sim.drawData.darkTileBackImage = loadImage("./simulator/images/tile_back_dark.bmp")
	sim.drawData.lightTileBackImage = loadImage("./simulator/images/tile_back_light.bmp")

	sim.drawData.redrawBoard = true

	sim.drawData.blackShader = loadShader("./simulator/shaders/black.kage")
	sim.drawData.transparentRedShader = loadShader("./simulator/shaders/transparent_red.kage")
	sim.drawData.colorShader = loadShader("./simulator/shaders/color.kage")
}

func (sim *Simulator) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Grey200)

	if sim.drawData.redrawBoard {
		sim.drawBoard()
		sim.drawOverlay()

		sim.drawData.redrawBoard = false
	}

	sim.drawPossibleTilePlacements()

	op := ebiten.DrawImageOptions{}

	op.GeoM.Scale(sim.drawData.scale, sim.drawData.scale)
	op.GeoM.Scale(sim.drawData.boardScale, sim.drawData.boardScale)
	op.GeoM.Translate(float64(sim.drawData.cameraOffset.X), float64(sim.drawData.cameraOffset.Y))

	screen.DrawImage(sim.drawData.boardImage, &op)
	screen.DrawImage(sim.drawData.possibleTilePlacementsImage, &op)

	hdInv := float64(1) / float64(sim.drawData.hdScale)
	opHD := ebiten.DrawImageOptions{}
	opHD.GeoM.Scale(sim.drawData.scale, sim.drawData.scale)
	opHD.GeoM.Scale(sim.drawData.boardScale, sim.drawData.boardScale)
	opHD.GeoM.Scale(hdInv, hdInv)
	opHD.GeoM.Translate(float64(sim.drawData.cameraOffset.X), float64(sim.drawData.cameraOffset.Y))
	screen.DrawImage(sim.drawData.overlayImg, &opHD)

	sim.drawUI(screen)
}

func (sim *Simulator) toWorldSpace(screenSpacePoint util.Point[int]) util.Point[float64] {
	return util.Point[float64]{
		X: float64(screenSpacePoint.X) - float64(sim.drawData.cameraOffset.X),
		Y: float64(screenSpacePoint.Y) - float64(sim.drawData.cameraOffset.Y),
	}
}

func (sim *Simulator) toBoardSpace(screenSpacePoint util.Point[int]) util.Point[int] {
	worldSpacePoint := sim.toWorldSpace(screenSpacePoint)
	boardSpacePoint := util.Point[int]{
		X: int(worldSpacePoint.X / (sim.drawData.boardScale * sim.drawData.scale * float64(TILE_SIZE))),
		Y: int(worldSpacePoint.Y / (sim.drawData.boardScale * sim.drawData.scale * float64(TILE_SIZE))),
	}
	return boardSpacePoint
}

func loadImage(fileName string) *ebiten.Image {
	reader, err := os.Open(fileName)

	if err != nil {
		panic(err)
	}

	image, err := bmp.Decode(reader)

	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(image)
}

func loadShader(fileName string) *ebiten.Shader {
	shaderCode, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	shader, err := ebiten.NewShader(shaderCode)

	if err != nil {
		panic(err)
	}

	return shader
}

func (sim *Simulator) drawBoard() {
	sim.drawData.boardImage.Clear()

	boardSize := sim.Engine.GameBoard.TileMatrix.Size()
	op := ebiten.DrawImageOptions{}

	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			tile := sim.Engine.GameBoard.TileMatrix.Get(x, y)
			tx, ty := float64(x*TILE_SIZE), float64(y*TILE_SIZE)
			op.GeoM.Translate(tx, ty)

			if tile != nil {
				//draw tile
				ebiImg := ebiten.NewImageFromImage(tile.Reference.Image)
				sim.drawData.boardImage.DrawImage(ebiImg, &op)

			} else {
				sim.drawData.boardImage.DrawImage(sim.drawData.emptyTileImage, &op)
			}

			op.GeoM.Translate(-tx, -ty)
		}
	}
}

func (sim *Simulator) drawOpenPositions() {
	var path vector.Path

	s := float64(sim.drawData.hdScale)

	for _, pt := range sim.Engine.GameBoard.OpenPositionsList() {

		x := pt.X
		y := pt.Y

		tSize := float32(float64(TILE_SIZE) * s)
		tx, ty := float32(float64(x*TILE_SIZE)*s), float32(float64(y*TILE_SIZE)*s)

		path.MoveTo(tx, ty)
		path.LineTo(tx+tSize, ty)
		path.LineTo(tx+tSize, ty+tSize)
		path.LineTo(tx, ty+tSize)
		path.LineTo(tx, ty)
	}

	op := &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.FillAll,
	}

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	sim.drawData.overlayImg.DrawTrianglesShader(vs, is, sim.drawData.transparentRedShader, op)

}

func (sim *Simulator) drawFeatureLinks() {
	var path vector.Path

	drawnFeatures := make(map[*tile.Feature]*tile.Feature)

	boardSize := sim.Engine.GameBoard.TileMatrix.Size()
	s := float64(sim.drawData.hdScale)
	var m float32 = 2
	for y := 0; y < boardSize; y++ {
		for x := 0; x < boardSize; x++ {
			t := sim.Engine.GameBoard.TileMatrix.Get(x, y)
			tx, ty := float32(float64(x*TILE_SIZE)*s), float32(float64(y*TILE_SIZE)*s)

			if t != nil {
				for _, rf := range t.Reference.Features {

					tileFeature := t.ReferenceFeatureMap[rf]
					avgPos := t.Reference.AvgFeaturePos[rf]
					aX := tx + float32(avgPos.X*s)
					aY := ty + float32(avgPos.Y*s)
					path.MoveTo(aX, aY)

					for _, lf := range tileFeature.Links {

						//dont redraw
						if _, exists := drawnFeatures[lf]; exists {
							continue
						}

						lfRefTile := lf.ParentTile.Reference
						lfAvgPos := lfRefTile.AvgFeaturePos[lf.ParentFeature]

						lfPos := lf.ParentTile.Position
						lfTx, lfTy := float32(float64(lfPos.X*TILE_SIZE)*s), float32(float64(lfPos.Y*TILE_SIZE)*s)
						lfaX := lfTx + float32(lfAvgPos.X*s)
						lfaY := lfTy + float32(lfAvgPos.Y*s)
						path.LineTo(lfaX-m, lfaY-m)
						path.LineTo(lfaX+m, lfaY+m)
						path.LineTo(aX+m, aY+m)
						path.LineTo(aX-m, aY-m)
						path.MoveTo(aX, aY)
					}

					drawnFeatures[tileFeature] = tileFeature
				}
			}
		}
	}

	op := &ebiten.DrawTrianglesShaderOptions{
		FillRule: ebiten.FillAll,
	}

	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)

	sim.drawData.overlayImg.DrawTrianglesShader(vs, is, sim.drawData.blackShader, op)
}

func (sim *Simulator) drawMeeples() {
	players := sim.Engine.Players

	s := float32(sim.drawData.hdScale)

	for _, p := range players {
		var path vector.Path

		for _, m := range p.Meeples {
			if m.Feature == nil {
				continue
			}

			var meepleScale float32 = float32(s) * 2

			f := m.Feature
			t := f.ParentTile

			avg := t.Reference.AvgFeaturePos[f.ParentFeature]

			x := float32(t.Position.X*TILE_SIZE) * s
			y := float32(t.Position.Y*TILE_SIZE) * s

			x += float32(avg.X) * s
			y += float32(avg.Y) * s

			path.MoveTo(x, y)
			path.LineTo(x+(0.5*meepleScale), y+(0.5*meepleScale))
			path.LineTo(x+(0*meepleScale), y+(1*meepleScale))
			path.LineTo(x-(0.5*meepleScale), y+(0.5*meepleScale))
			path.LineTo(x-(0*meepleScale), y+(0*meepleScale))
		}

		op := &ebiten.DrawTrianglesShaderOptions{
			FillRule: ebiten.FillAll,
		}

		op.Uniforms = make(map[string]interface{})
		r, g, b, a := p.Color.RGBA()
		op.Uniforms["RGBA"] = []float32{float32(r) / 65535, float32(g) / 65535, float32(b) / 65535, float32(a) / 65535}

		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		sim.drawData.overlayImg.DrawTrianglesShader(vs, is, sim.drawData.colorShader, op)
	}

}

func (sim *Simulator) drawOverlay() {
	sim.drawData.overlayImg.Clear()
	//sim.drawOpenPositions()
	//sim.drawFeatureLinks()
	sim.drawMeeples()
}

func (sim *Simulator) drawPossibleTilePlacements() {
	sim.drawData.possibleTilePlacementsImage.Clear()

	op := ebiten.DrawImageOptions{}

	for _, placement := range sim.Engine.CurrentPossibleTilePlacements {
		tx, ty := float64(placement.Position.X*TILE_SIZE), float64(placement.Position.Y*TILE_SIZE)
		op.GeoM.Translate(tx, ty)

		sim.drawData.possibleTilePlacementsImage.DrawImage(sim.drawData.lightTileBackImage, &op)

		op.GeoM.Translate(-tx, -ty)
	}
}

func (sim *Simulator) drawCurrentlyHeldTile() {
	img := sim.drawData.currentlyHeldTileImage

	if sim.Engine.HeldRefTileGroup == nil {
		img.Fill(color.Black)
		return
	}

	tileImg := sim.Engine.HeldRefTileGroup.Orientations[0].Image

	ebiTileImg := ebiten.NewImageFromImage(tileImg)

	op := ebiten.DrawImageOptions{}
	img.DrawImage(ebiTileImg, &op)
}

func (sim *Simulator) drawDeck(deck deck.Deck) {
	img := sim.drawData.deckImage

	op := &ebiten.DrawImageOptions{}

	topTileRefGroup, err := deck.Scry()

	//if no tiles left, then just use the default img (tile back probably), or else black
	if err != nil {
		img.Fill(color.Black)
		return
	}

	tileImg := topTileRefGroup.Orientations[0].Image

	ebiTileImg := ebiten.NewImageFromImage(tileImg)
	img.DrawImage(ebiTileImg, op)
}

func (sim *Simulator) drawUI(screen *ebiten.Image) {
	cursorX, cursorY := ebiten.CursorPosition()

	text.Draw(
		screen,
		fmt.Sprint(
			"FPS:", int(ebiten.CurrentFPS()),
			" TPS: ", int(ebiten.CurrentTPS()),
			" CAM_POS: (X: ", int(sim.drawData.cameraOffset.X), ", Y: ", int(sim.drawData.cameraOffset.Y), ")",
			" CUR_POS: (X: ", cursorX, ", Y: ", cursorY, ")",
			" STEP: ", sim.Engine.TurnCounter, "-", sim.Engine.TurnStage,
		),
		sim.drawData.font,
		2,
		12,
		color.Black,
	)

	text.Draw(
		screen,
		fmt.Sprint(
			"Player:", sim.Engine.CurrentPlayer().Name,
		),
		sim.drawData.font,
		2,
		24,
		sim.Engine.CurrentPlayer().Color,
	)

	screenMaxX := float64(screen.Bounds().Dx())
	screenMaxY := float64(screen.Bounds().Dy())
	totalUIScale := sim.drawData.uiScale * sim.drawData.scale
	scaledTileSize := float64(TILE_SIZE) * totalUIScale

	//draw deck onto screen
	op := &ebiten.DrawImageOptions{}
	sim.drawDeck(*sim.Engine.Deck)
	op.GeoM.Scale(totalUIScale, totalUIScale)
	op.GeoM.Translate(screenMaxX-scaledTileSize-4, screenMaxY-scaledTileSize-4)
	screen.DrawImage(sim.drawData.deckImage, op)

	r := sim.Engine.RiverDeck.Remaining()
	if r > 0 {
		screen.DrawImage(sim.drawData.lightTileBackImage, op)
	} else {
		screen.DrawImage(sim.drawData.deckImage, op)
	}

	//draw river deck onto screen
	op = &ebiten.DrawImageOptions{}
	sim.drawDeck(*sim.Engine.RiverDeck)
	op.GeoM.Scale(totalUIScale, totalUIScale)
	op.GeoM.Translate(screenMaxX-(scaledTileSize+4)*2, screenMaxY-scaledTileSize-4)
	screen.DrawImage(sim.drawData.deckImage, op)

	op = &ebiten.DrawImageOptions{}
	sim.drawCurrentlyHeldTile()
	op.GeoM.Scale(totalUIScale, totalUIScale)
	op.GeoM.Translate(screenMaxX-(scaledTileSize+4)*3, screenMaxY-scaledTileSize-4)
	screen.DrawImage(sim.drawData.currentlyHeldTileImage, op)
}

func (sim *Simulator) initFont() {
	const dpi = 72

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		panic(err)
	}

	font, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})

	if err != nil {
		panic(err)
	}

	sim.drawData.font = font
}
