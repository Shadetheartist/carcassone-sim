package simulator

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/exp/shiny/materialdesign/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type DrawData struct {
	font font.Face
}

func (sim *Simulator) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Blue300)

	sim.drawUI(screen)
}

func (sim *Simulator) drawUI(screen *ebiten.Image) {
	screen.Fill(colornames.Blue400)

	text.Draw(
		screen,
		fmt.Sprint("FPS:", int(ebiten.CurrentFPS()), " TPS: ", int(ebiten.CurrentTPS())),
		sim.DrawData.font,
		2,
		12,
		color.White,
	)
}

func (sim *Simulator) setupFont() {
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

	sim.DrawData.font = font
}
