package simulator

import (
	"beeb/carcassonne/board"
	"beeb/carcassonne/data"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/exp/shiny/materialdesign/colornames"
)

var PLAYER_COLOR_LIST = [...]color.RGBA{
	colornames.White,
	colornames.Red500,
	colornames.Blue500,
	colornames.Green500,
	colornames.Black,
}

type Simulator struct {
	GameBoard *board.Board
	GameData  *data.GameData
	DrawData  *DrawData

	Players       []*Player
	CurrentPlayer int
}

func NewSimulator(gameData *data.GameData, boardSize int, numPlayers int) *Simulator {

	if numPlayers > len(PLAYER_COLOR_LIST) {
		panic(fmt.Sprint("too many players, max ", len(PLAYER_COLOR_LIST)))
	}

	sim := &Simulator{}

	sim.GameBoard = board.NewBoard(boardSize)
	sim.GameData = gameData
	sim.Players = make([]*Player, numPlayers)

	for i := 0; i < numPlayers; i++ {
		playerName := fmt.Sprint("Player ", i)
		sim.Players[i] = NewPlayer(playerName, PLAYER_COLOR_LIST[i])
	}

	sim.DrawData = &DrawData{}
	sim.setupFont()

	return sim
}

func (sim *Simulator) Simulate() {
	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")
	ebiten.SetScreenClearedEveryFrame(false)

	if err := ebiten.RunGame(sim); err != nil {
		panic(err)
	}
}

func (sim *Simulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 700, 700
}
