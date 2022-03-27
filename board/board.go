package board

import (
	"beeb/carcassonne/board/road"
	dir "beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"errors"
	"image"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

const imageW = 4000
const imageH = 4000

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
			td := p.GridToTileDir(d)
			tileFeature := t.Feature(td)

			ntd := n.Placement.GridToTileDir(dir.Compliment[d])
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

func (b *Board) DirectionalNeighbours(position tile.Position) []*tile.Tile {
	neighbours := make([]*tile.Tile, 0, 4)

	for _, d := range dir.List {
		neighbourPosition := position.EdgePos(d)

		if neighbourTile, exists := b.Tiles[neighbourPosition]; exists {
			neighbours = append(neighbours, neighbourTile)
		} else {
			neighbours = append(neighbours, nil)
		}

	}

	return neighbours
}

func (b *Board) AddTile(t *tile.Tile, p tile.Placement) error {
	neighbours := b.DirectionalNeighbours(p.Position)

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
			td := p.GridToTileDir(d)
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

func (b *Board) ptpGoRoutine(positions []tile.Position, t *tile.Tile, c chan []tile.Placement, wg *sync.WaitGroup) {
	defer wg.Done()
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

	var wg sync.WaitGroup

	sectionSize := int(math.Floor(float64(len(keys)) / float64(numRoutines)))

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)

		start := i * sectionSize
		end := min((i+1)*sectionSize, len(keys))
		go b.ptpGoRoutine(keys[start:end], t, c, &wg)
	}

	wg.Wait()

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
