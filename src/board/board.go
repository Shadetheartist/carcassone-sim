package board

import (
	"beeb/carcassonne/board/road"
	dir "beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"errors"
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

const imageW = 1000
const imageH = 1000

type Board struct {
	Tiles     map[tile.Position]*tile.Tile
	Roads     map[int]*road.Road
	RoadCount int
	//dir mapped
	OpenPositions map[tile.Position][]tile.FeatureType

	BoardImage *ebiten.Image
	RoadsImage *ebiten.Image
}

func New() Board {
	b := Board{}
	b.Tiles = make(map[tile.Position]*tile.Tile)
	b.Roads = make(map[int]*road.Road)
	b.OpenPositions = make(map[tile.Position][]tile.FeatureType, 4)

	b.BoardImage = ebiten.NewImage(imageW, imageH)
	b.RoadsImage = ebiten.NewImage(imageW, imageH)

	return b
}

func (b *Board) CompareEdge(t *tile.Tile, p tile.Placement, d dir.Direction) (*tile.Tile, error) {
	tileDir := p.TileDirection(d)
	neighbourPos := p.Position.EdgePos(d)

	if neighbour, exists := b.Tiles[neighbourPos]; exists {
		neighbourComplimentDir := dir.Compliment[d]
		neighbourTileDir := neighbour.Placement.TileDirection(neighbourComplimentDir)

		if t.Feature(tileDir).Type != neighbour.Feature(neighbourTileDir).Type {
			return neighbour, errors.New(fmt.Sprint(d, " Feature Type Mismatch"))
		}

		return neighbour, nil
	}

	return nil, nil
}

//returns a dir map (use dir as index) of all valid neighbours
func (b *Board) VerifyPlacement(t *tile.Tile, position tile.Placement) ([]*tile.Tile, error) {

	neighbours := make([]*tile.Tile, 4)

	//for each orthoganal direction, compare the edges of the tile arg and its neighbours
	//if there is a mismatch error, return err
	for _, d := range dir.List {
		n, err := b.CompareEdge(t, position, d)

		if err != nil {
			return neighbours, err
		}

		neighbours[d] = n
	}

	return neighbours, nil
}

func (b *Board) PlaceableOrientations(t *tile.Tile, position tile.Position) []uint16 {
	orientations := make([]uint16, 0, 4)

	op := b.OpenPositions[position]
	eft := t.EdgeFeatureTypes

	for i := 0; i < 4; i++ {

		var match bool = true
		//compare each element in the array
		for j := range op {

			//we dont care about matching 0's
			if op[j] == 0 {
				continue
			}

			//if a non-blank mismatch occurs then this can't be the right orientation
			if op[j] != eft[j] {
				match = false
				break
			}

		}

		if match {
			orientations = append(orientations, uint16((i*90)%360))
		}

		//try shifting everything over one (i.e. rotate the tile 90 degrees)
		tmp := eft[3]
		eft[3] = eft[2]
		eft[2] = eft[1]
		eft[1] = eft[0]
		eft[0] = tmp
	}

	return orientations
}

func (b *Board) IsTilePlaceable(t *tile.Tile, position tile.Position) (uint16, error) {

	op := b.OpenPositions[position]
	eft := t.EdgeFeatureTypes

	for i := 0; i < 4; i++ {

		var match bool = true
		//compare each element in the array
		for j := range op {

			//we dont care about matching 0's
			if op[j] == 0 {
				continue
			}

			//if a non-blank mismatch occurs then this can't be the right orientation
			if op[j] != eft[j] {
				match = false
				break
			}

		}

		if match {
			//if matched then its all good, return the orientation
			return uint16((i * 90) % 360), nil
		}

		//if not matched, try shifting everything over one (i.e. rotate the tile 90 degrees)
		tmp := eft[3]
		eft[3] = eft[2]
		eft[2] = eft[1]
		eft[1] = eft[0]
		eft[0] = tmp
	}

	return 0, errors.New("Not Placable")
}

func (b *Board) ConnectedFeatures(t *tile.Tile, p tile.Placement) map[dir.Direction]tile.FeatureType {
	connectedFeatures := make(map[dir.Direction]tile.FeatureType)

	neighbours := b.OrthoganalNeighbours(p.Position)

	for _, n := range neighbours {

		if n == nil {
			continue
		}

		for _, d := range dir.List {
			td := p.TileDirection(d)
			tileFeature := t.Feature(td)

			ntd := n.Placement.TileDirection(dir.Compliment[d])
			neighbourTileFeature := n.Feature(ntd)

			if tileFeature.Type == neighbourTileFeature.Type {
				connectedFeatures[d] = tileFeature.Type
			}
		}
	}

	return connectedFeatures
}

func (b *Board) newRoad() int {
	b.RoadCount++
	b.Roads[b.RoadCount] = &road.Road{}
	return b.RoadCount
}

func (b *Board) AddTile(t *tile.Tile, p tile.Placement) error {
	neighbours, err := b.VerifyPlacement(t, p)

	if err != nil {
		return err
	}

	//if there are no mismatched tiles,
	for _, d := range dir.List {
		if neighbours[d] != nil {
			neighbour := neighbours[d]

			//set the neighbour relationships between the tiles for later use
			t.Neighbours[d] = neighbour
			neighbour.Neighbours[dir.Compliment[d]] = t

		}
	}

	//place the tile at the position for itself and the board
	t.Placement = p
	t.Neighbours = neighbours
	b.Tiles[p.Position] = t

	for _, d := range dir.List {
		if neighbours[d] == nil {

			//open position
			neighbourPosition := p.Position.EdgePos(d)
			td := p.TileDirection(d)
			featureType := t.Feature(td).Type
			dc := dir.Compliment[d]

			if op, exists := b.OpenPositions[neighbourPosition]; exists {
				op[dc] = featureType
			} else {
				op := make([]tile.FeatureType, 4)
				op[dc] = featureType
				b.OpenPositions[neighbourPosition] = op
			}
		}
	}

	delete(b.OpenPositions, p.Position)

	return nil
}

func (b *Board) ptp(positions []tile.Position, t *tile.Tile) []tile.Placement {

	placements := make([]tile.Placement, 0, len(positions))

	for _, openPos := range positions {
		orientations := b.PlaceableOrientations(t, openPos)

		for _, o := range orientations {
			placements = append(placements, tile.Placement{
				Position:    openPos,
				Orientation: o,
			})
		}
	}

	return placements
}

func (b *Board) ptpGoRoutine(positions []tile.Position, t *tile.Tile, c chan []tile.Placement) {
	c <- b.ptp(positions, t)
}

func (b *Board) PossibleTilePlacements(t *tile.Tile) []tile.Placement {

	placements := make([]tile.Placement, 0, len(b.OpenPositions))

	numRoutines := len(b.OpenPositions) / 64

	keys := make([]tile.Position, 0, len(b.OpenPositions))
	for p := range b.OpenPositions {
		keys = append(keys, p)
	}

	if numRoutines <= 1 {
		return b.ptp(keys, t)
	}

	c := make(chan []tile.Placement, numRoutines)

	sectionSize := len(keys) / numRoutines

	for i := 0; i < numRoutines; i++ {
		start := i * sectionSize
		end := min((i+1)*sectionSize, len(keys)) - 1
		go b.ptpGoRoutine(keys[start:end], t, c)
	}

	for i := 0; i < numRoutines; i++ {
		placements = append(placements, <-c...)
	}

	return placements
}

func (b *Board) OrthoganalNeighbours(p tile.Position) map[tile.Position]*tile.Tile {
	neighbours := make(map[tile.Position]*tile.Tile)

	for _, d := range dir.List {
		edgePos := p.EdgePos(d)
		if neighbour, exists := b.Tiles[edgePos]; exists {
			neighbours[edgePos] = neighbour
		} else {
			neighbours[edgePos] = nil
		}
	}

	return neighbours
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b *Board) BoundingBox() image.Rectangle {
	rect := image.Rectangle{}

	for pos := range b.Tiles {
		rect.Max.X = max(pos.X, rect.Max.X)
		rect.Max.Y = max(pos.Y, rect.Max.Y)
		rect.Min.X = min(pos.X, rect.Min.X)
		rect.Min.Y = min(pos.X, rect.Min.X)
	}

	return rect
}
