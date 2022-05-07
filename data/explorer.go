package data

import (
	"beeb/carcassonne/tile"
	"fmt"
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type GameDataExplorer struct {
	font                 font.Face
	cursorPosition       image.Point
	cameraOffset         image.Point
	scale                float64
	tileScale            float64
	mouseDown            bool
	mouseInitialPosition image.Point
	gameData             *GameData
	tilePositions        map[*tile.ReferenceTile]image.Point
	tileImages           map[*tile.ReferenceTile]*ebiten.Image
	selectedFeature      *tile.Feature
	redrawFeatures       bool
	tilesImage           *ebiten.Image
	overlayImage         *ebiten.Image
	overlayCtx           *gg.Context
}

func newGameDataExplorer(gameData *GameData) *GameDataExplorer {
	gde := GameDataExplorer{}

	gde.gameData = gameData
	gde.setupFont()

	gde.scale = 1
	gde.tileScale = 10
	gde.cursorPosition = image.Point{}

	gde.tilePositions = make(map[*tile.ReferenceTile]image.Point, len(gameData.ReferenceTiles))

	gde.setupTileImages()
	gde.setupTilesImage()
	gde.setupOverlayImage()

	return &gde
}

func (gde *GameDataExplorer) setupTileImages() {
	gde.tileImages = make(map[*tile.ReferenceTile]*ebiten.Image, len(gde.gameData.ReferenceTiles))

	row := 0
	rowMargin := 20

	for _, tileName := range gde.gameData.TileNames {
		rts := gde.gameData.ReferenceTiles[tileName]
		for i := 0; i < 4; i++ {
			rt := rts[i]
			gde.tilePositions[rt] = image.Pt(i*8, row*rowMargin)
			gde.tileImages[rt] = ebiten.NewImageFromImage(rt.Image)
		}

		row++
	}

}

func (gde *GameDataExplorer) setupTilesImage() {
	row := 0
	rowMargin := 20

	rtImageRect := image.Rect(0, 0, 4*8*int(gde.tileScale), len(gde.gameData.ReferenceTiles)*rowMargin*int(gde.tileScale))
	gde.tilesImage = ebiten.NewImage(rtImageRect.Dx(), rtImageRect.Dy())

	for _, tileName := range gde.gameData.TileNames {
		rts := gde.gameData.ReferenceTiles[tileName]
		for i := 0; i < 4; i++ {
			rt := rts[i]
			pos := gde.tilePositions[rt]
			eImg := gde.tileImages[rt]

			rtOp := ebiten.DrawImageOptions{}
			rtOp.GeoM.Translate(float64(pos.X), float64(pos.Y))
			rtOp.GeoM.Scale(float64(gde.tileScale), float64(gde.tileScale))

			gde.tilesImage.DrawImage(eImg, &rtOp)

			text.Draw(
				gde.tilesImage,
				fmt.Sprint("O:", rt.Orientation),
				gde.font,
				pos.X*int(gde.tileScale),
				pos.Y*int(gde.tileScale)+(11*int(gde.tileScale)),
				color.White,
			)
		}

		text.Draw(
			gde.tilesImage,
			tileName,
			gde.font,
			0,
			row*rowMargin*int(gde.tileScale)+(9*int(gde.tileScale)),
			color.White,
		)

		row++
	}
}

func (gde *GameDataExplorer) setupOverlayImage() {
	rowMargin := 20
	rtImageRect := image.Rect(0, 0, 4*8*int(gde.tileScale), len(gde.gameData.ReferenceTiles)*rowMargin*int(gde.tileScale))
	gde.overlayImage = ebiten.NewImage(rtImageRect.Dx(), rtImageRect.Dy())
	gde.overlayCtx = gg.NewContext(rtImageRect.Dx()/int(gde.tileScale), rtImageRect.Dy()/int(gde.tileScale))
}

func (gde *GameDataExplorer) setupFont() {
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

	gde.font = font
}

func (gde *GameDataExplorer) tilePixelSpace(x int, y int) (int, int) {
	x = int(float64(x-gde.cameraOffset.X) / (gde.scale * gde.tileScale))
	y = int(float64(y-gde.cameraOffset.Y) / (gde.scale * gde.tileScale))

	return x, y
}

func (gde *GameDataExplorer) drawFeatures() {
	row := 0

	gde.overlayImage.Clear()

	gde.overlayCtx.SetColor(color.Transparent)
	gde.overlayCtx.Clear()

	gde.overlayCtx.SetColor(colornames.Red300)

	for _, tileName := range gde.gameData.TileNames {
		rts := gde.gameData.ReferenceTiles[tileName]
		for i := 0; i < 4; i++ {
			rt := rts[i]
			pos := gde.tilePositions[rt]

			for _, f := range rt.Features {
				if f == gde.selectedFeature {
					rt.FeatureMatrix.Iterate(func(f *tile.Feature, x int, y int, idx int) {
						if f == gde.selectedFeature {
							gde.overlayCtx.SetPixel(pos.X+x, pos.Y+y)
						}
					})
					break
				}
			}
		}

		row++
	}

	ggImg := gde.overlayCtx.Image()
	ebiImg := ebiten.NewImageFromImage(ggImg)

	op := ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(gde.tileScale), float64(gde.tileScale))

	gde.overlayImage.DrawImage(ebiImg, &op)
}

func (gde *GameDataExplorer) Draw(screen *ebiten.Image) {

	screen.Fill(colornames.Grey600)

	if gde.redrawFeatures {
		gde.drawFeatures()
		gde.redrawFeatures = false
	}

	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(gde.cameraOffset.X)/gde.scale, float64(gde.cameraOffset.Y)/gde.scale)
	op.GeoM.Scale(gde.scale, gde.scale)

	screen.DrawImage(gde.tilesImage, &op)
	screen.DrawImage(gde.overlayImage, &op)

	tx, ty := gde.tilePixelSpace(gde.cursorPosition.X, gde.cursorPosition.Y)
	text.Draw(
		screen,
		fmt.Sprint("X:", gde.cursorPosition.X, " TX: ", tx),
		gde.font,
		10,
		10,
		color.White,
	)

	text.Draw(
		screen,
		fmt.Sprint("Y:", gde.cursorPosition.Y, " TY: ", ty),
		gde.font,
		10,
		25,
		color.White,
	)

	text.Draw(
		screen,
		fmt.Sprint("FPS:", int(ebiten.CurrentFPS()), " TPS: ", int(ebiten.CurrentTPS())),
		gde.font,
		10,
		40,
		color.White,
	)

}

func (gde *GameDataExplorer) featureUnderCursor() *tile.Feature {
	tx, ty := gde.tilePixelSpace(gde.cursorPosition.X, gde.cursorPosition.Y)

	var tileUnderCursor *tile.ReferenceTile
	var pixelUnderCursor image.Point
	for rt, pos := range gde.tilePositions {
		if tx >= pos.X && tx < pos.X+7 && ty >= pos.Y && ty < pos.Y+7 {
			tileUnderCursor = rt
			pixelUnderCursor = image.Pt(tx-pos.X, ty-pos.Y)
			break
		}
	}

	if tileUnderCursor != nil {
		feature := tileUnderCursor.FeatureMatrix.Get(pixelUnderCursor.X, pixelUnderCursor.Y)
		return feature
	}

	return nil
}

func (gde *GameDataExplorer) Update() error {

	cX, cY := ebiten.CursorPosition()

	gde.cursorPosition.X = cX
	gde.cursorPosition.Y = cY

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		f := gde.featureUnderCursor()
		if f != nil && f != gde.selectedFeature {
			gde.selectedFeature = f
			gde.redrawFeatures = true
			fmt.Printf("%s %s %p\n", gde.selectedFeature.ParentRefenceTile.Name, gde.selectedFeature.Type, gde.selectedFeature)
		}
	}

	//mouse state change tracking
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		gde.mouseDown = false
	}

	//panning support
	if !gde.mouseDown && ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		mouseInitialX, mouseInitialY := ebiten.CursorPosition()
		gde.mouseInitialPosition.X = mouseInitialX
		gde.mouseInitialPosition.Y = mouseInitialY
		gde.mouseDown = true
	}

	//panning
	if gde.mouseDown {

		gde.cameraOffset.X += cX - gde.mouseInitialPosition.X
		gde.cameraOffset.Y += cY - gde.mouseInitialPosition.Y

		gde.mouseInitialPosition.X = cX
		gde.mouseInitialPosition.Y = cY
	}

	return nil
}

func (gde *GameDataExplorer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 700, 700
}

func (gd *GameData) Explore() {
	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")
	ebiten.SetScreenClearedEveryFrame(false)

	explorer := newGameDataExplorer(gd)

	if err := ebiten.RunGame(explorer); err != nil {
		panic(err)
	}
}
