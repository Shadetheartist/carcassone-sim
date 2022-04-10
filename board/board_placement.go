package board

import (
	"beeb/carcassonne/tile"
	"errors"
	"math"
	"sync"
)

//determines all the possible placements (positions & orientation combination) a tile could be put in on the board
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

	sectionSize := int(math.Ceil(float64(len(keys)) / float64(numRoutines)))

	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		start := min(i*sectionSize, len(keys))
		end := min((i+1)*sectionSize, len(keys))
		go b.ptpGoRoutine(keys[start:end], t, c, &wg)
	}

	wg.Wait()

	for i := 0; i < numRoutines; i++ {
		placements = append(placements, <-c...)
	}

	return placements
}

//determines if a tile has an orientation that allows it to be placed at a certain position
func (b *Board) IsTilePlaceable(t *tile.Tile, position tile.Position) (uint16, error) {

	op := b.OpenPositions[position]
	eft := t.CachedEdgeFeatureTypes

	for i := 0; i < 4; i++ {
		if compareEft(op, eft) {
			return tile.LimitToOrientation(i), nil
		}

		//if not matched, try shifting everything over one (i.e. rotate the tile 90 degrees)
		shiftEft(&eft)
	}

	return 0, errors.New("Not Placable")
}

//builds a slice of placable orientatations for a position
func (b *Board) PlaceableOrientations(t *tile.Tile, position tile.Position) []uint16 {
	orientations := make([]uint16, 0, 4)

	op := b.OpenPositions[position]
	eft := t.CachedEdgeFeatureTypes

	for i := 0; i < 4; i++ {
		if compareEft(op, eft) {
			orientations = append(orientations, tile.LimitToOrientation(i))
		}

		//try shifting everything over one (i.e. rotate the tile 90 degrees)
		shiftEft(&eft)
	}

	return orientations
}

func (b *Board) ptpGoRoutine(positions []tile.Position, t *tile.Tile, c chan []tile.Placement, wg *sync.WaitGroup) {
	defer wg.Done()
	c <- b.ptp(positions, t)
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

//rotates the array by 90 degrees essentially
func shiftEft(eft *[4]tile.FeatureType) {
	tmp := eft[3]
	eft[3] = eft[2]
	eft[2] = eft[1]
	eft[1] = eft[0]
	eft[0] = tmp
}

func compareEft(op []tile.FeatureType, eft [4]tile.FeatureType) bool {
	//compare each element in the array
	for j := range op {

		//we dont care about matching 0's
		if op[j] == 0 {
			continue
		}

		//if a non-blank mismatch occurs then this can't be the right orientation
		if op[j] != eft[j] {
			return false
		}
	}

	return true
}
