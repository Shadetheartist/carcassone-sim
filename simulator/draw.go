package simulator

import (
	"beeb/carcassonne/engine/deck"
	"beeb/carcassonne/util"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
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

	riverDeckImage         *ebiten.Image
	deckImage              *ebiten.Image
	currentlyHeldTileImage *ebiten.Image

	scale        float64
	boardScale   float64
	uiScale      float64
	cameraOffset util.Point[int]
}

func (sim *Simulator) initDraw() {
	sim.initFont()
	sim.drawData.scale = 2
	sim.drawData.uiScale = 4
	sim.drawData.boardScale = 1
	sim.drawData.cameraOffset = util.Point[int]{}

	boardPxSize := sim.Engine.GameBoard.TileMatrix.Size() * TILE_SIZE
	sim.drawData.boardImage = ebiten.NewImage(boardPxSize, boardPxSize)
	sim.drawData.possibleTilePlacementsImage = ebiten.NewImage(boardPxSize, boardPxSize)

	sim.drawData.deckImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	sim.drawData.riverDeckImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)
	sim.drawData.currentlyHeldTileImage = ebiten.NewImage(TILE_SIZE, TILE_SIZE)

	sim.drawData.emptyTileImage = loadImage("./simulator/images/empty_tile.bmp")
	sim.drawData.darkTileBackImage = loadImage("./simulator/images/tile_back_dark.bmp")
	sim.drawData.lightTileBackImage = loadImage("./simulator/images/tile_back_light.bmp")

}

func (sim *Simulator) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Grey200)

	sim.drawBoard()
	sim.drawPossibleTilePlacements()

	op := ebiten.DrawImageOptions{}

	op.GeoM.Scale(sim.drawData.scale, sim.drawData.scale)
	op.GeoM.Scale(sim.drawData.boardScale, sim.drawData.boardScale)
	op.GeoM.Translate(float64(sim.drawData.cameraOffset.X), float64(sim.drawData.cameraOffset.Y))

	screen.DrawImage(sim.drawData.boardImage, &op)
	screen.DrawImage(sim.drawData.possibleTilePlacementsImage, &op)

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
		),
		sim.drawData.font,
		2,
		12,
		color.Black,
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
