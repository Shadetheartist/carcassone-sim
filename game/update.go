package game

import (
	"beeb/carcassonne/board"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Update() error {

	g.handleMouseInput()

	return nil
}

//input state
var mouseDown bool = false
var mouseInitialX, mouseInitialY int

func (g *Game) handleMouseInput() {

	cX, cY := ebiten.CursorPosition()

	g.HoveredPosition = g.Board.PixToPos(image.Point{
		X: cX - g.CameraOffset.X,
		Y: cY - g.CameraOffset.Y,
	})

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) && g.HoveredPosition != g.SelectedPosition {
		g.SelectedPosition = g.HoveredPosition

		if _, exists := g.Board.Tiles[g.HoveredPosition]; exists {
			clickedTile := g.Board.Tiles[g.HoveredPosition]

			roads := make([]board.Road, 0)
			for _, rs := range clickedTile.UniqueRoadSegements() {
				rd := board.CompileRoadFromSegment(rs)
				roads = append(roads, rd)
			}

			g.HighlightedRoads = roads

		}

		if _, exists := g.Board.OpenPositions[g.HoveredPosition]; exists {
			t, err := g.findTileForPos(g.SelectedPosition)

			if err == nil {
				g.Board.AddTile(t, t.Placement)
			}
		}
	}

	//mouse state change tracking
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseDown = false
	}

	//panning support
	if !mouseDown && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseInitialX, mouseInitialY = ebiten.CursorPosition()
		mouseDown = true
	}

	//panning
	if mouseDown {

		g.CameraOffset.X += cX - mouseInitialX
		g.CameraOffset.Y += cY - mouseInitialY

		mouseInitialX = cX
		mouseInitialY = cY
	}

}
