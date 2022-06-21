package board

import (
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
	"fmt"
	"image"
	"time"
)

type Board struct {
	TileMatrix        *matrix.Matrix[*tile.Tile]
	PlacedTileCount   int
	OpenPositions     map[util.Point[int]]*tile.EdgeSignature
	OpenPositionsList []util.Point[int]
	EdgePixReference  [][]util.Point[int]
}

func NewBoard(size int) *Board {
	board := &Board{}

	board.TileMatrix = matrix.NewMatrix[*tile.Tile](size)
	board.OpenPositions = make(map[util.Point[int]]*tile.EdgeSignature, 128)
	board.OpenPositionsList = make([]util.Point[int], 0, 128)
	board.EdgePixReference = edgePix(image.Rect(0, 0, 7, 7))

	return board
}

func (b *Board) RemoveTileAt(pos util.Point[int]) {
	tStart := time.Now()

	t := b.TileMatrix.Get(pos.X, pos.Y)

	//no tile had yet been placed there
	if t == nil {
		return
	}

	// clear out neighbour relationships
	for i := 0; i < 4; i++ {
		dir := directions.Direction(i)
		complimentDir := directions.Compliment[dir]
		if tl := t.Neighbours[dir]; tl != nil {
			tl.Neighbours[complimentDir] = nil
		}
		t.Neighbours[dir] = nil
	}

	//remove the tile from the matrix
	b.TileMatrix.Set(pos.X, pos.Y, nil)

	// add back the vacant position
	b.OpenPositionsList = append(b.OpenPositionsList, pos)
	b.OpenPositions[pos] = b.createOpenPositonSignature(pos)

	b.PlacedTileCount--
	t.Position = util.Point[int]{}

	// remove the vacancies created by this tile, if any,
	// by looking at the positions around the tile we removed,
	// if any have no real neighbours, we must remove them, as they are floating
	hasNeighbours := false
	for _, pn := range pos.OrthogonalNeighbours() {

		// look though the adjacent tile's neighbours,
		// if any are set, that open position can remain open
		for _, p2 := range pn.OrthogonalNeighbours() {
			matrixTile, err := b.TileMatrix.GetPt(p2)

			if err != nil {
				continue
			}

			if matrixTile != nil {
				hasNeighbours = true
				break
			}
		}

		// if the position has no real neigbours, it must be removed
		if !hasNeighbours {
			//remove position from list
			for i, p := range b.OpenPositionsList {
				if p != pn {
					continue
				}
				b.OpenPositionsList[i] = b.OpenPositionsList[len(b.OpenPositionsList)-1]
				b.OpenPositionsList = b.OpenPositionsList[:len(b.OpenPositionsList)-1]
				break
			}

			//remove position from map
			delete(b.OpenPositions, pn)
		}
	}

	//clear out all feature links
	for _, f := range t.Features {
		for l := range f.Links {
			delete(l.Links, f)
			delete(f.Links, l)
		}
	}

	fmt.Println("t us:", time.Since(tStart).Microseconds())
}

func (b *Board) linkNeighbours(t *tile.Tile) {
	//link neighbours
	for i := 0; i < 4; i++ {
		dir := directions.Direction(i)
		complimentDir := directions.Compliment[dir]

		pt := t.Position.EdgePos(dir)

		if tl, err := b.TileMatrix.GetPt(pt); err == nil {
			t.Neighbours[dir] = tl
			if tl != nil {
				tl.Neighbours[complimentDir] = t
			}
		}
	}
}

func (b *Board) PlaceTile(pos util.Point[int], t *tile.Tile) {
	b.TileMatrix.Set(pos.X, pos.Y, t)
	b.PlacedTileCount++
	t.Position = pos

	b.linkNeighbours(t)

	for d, n := range t.Neighbours {
		//if neighbour is not there add it as an open position
		if n == nil {
			edgePos := pos.EdgePos(directions.Direction(d))

			//don't add neighbours which are out of bounds
			if !b.TileMatrix.IsInBounds(edgePos.X, edgePos.Y) {
				continue
			}

			//add position to list, if it's not already in the map
			if _, exists := b.OpenPositions[edgePos]; !exists {
				b.OpenPositionsList = append(b.OpenPositionsList, edgePos)
			}

			//add position to map
			b.OpenPositions[edgePos] = b.createOpenPositonSignature(edgePos)
		} else {
			dir := directions.Direction(d)
			complimentDir := directions.Compliment[dir]

			pix := b.EdgePixReference[dir]
			complimentPix := b.EdgePixReference[complimentDir]

			matrix := t.Reference.FeatureMatrix
			neighbourFarmMatrix := n.Reference.FeatureMatrix

			for i := range pix {
				referenceFeature, err := matrix.GetPt(pix[i])

				if err != nil {
					panic(err)
				}

				if referenceFeature == nil {
					continue
				}

				neighbourReferenceFeature, err := neighbourFarmMatrix.GetPt(complimentPix[i])

				if err != nil {
					panic(err)
				}

				if neighbourReferenceFeature == nil {
					continue
				}

				if referenceFeature.Type == neighbourReferenceFeature.Type {
					tileFeature := t.ReferenceFeatureMap[referenceFeature]
					neighbourFeature := n.ReferenceFeatureMap[neighbourReferenceFeature]

					tileFeature.Links[neighbourFeature] = neighbourFeature
					neighbourFeature.Links[tileFeature] = tileFeature
				}
			}

		}
	}

	//remove position from list
	for i, p := range b.OpenPositionsList {
		if p != pos {
			continue
		}

		b.OpenPositionsList[i] = b.OpenPositionsList[len(b.OpenPositionsList)-1]
		b.OpenPositionsList = b.OpenPositionsList[:len(b.OpenPositionsList)-1]
		break
	}

	//remove position from map
	delete(b.OpenPositions, pos)
}

func (b *Board) createOpenPositonSignature(pos util.Point[int]) *tile.EdgeSignature {
	sig := &tile.EdgeSignature{}

	for i := 0; i < 4; i++ {
		dir := directions.Direction(i)
		complimentDirection := directions.Compliment[dir]
		t, err := b.TileMatrix.GetPt(pos.EdgePos(dir))

		if err != nil || t == nil {
			continue
		}

		sig[dir] = t.Reference.EdgeSignature[complimentDirection]
	}

	return sig
}

func edgePix(rect image.Rectangle) [][]util.Point[int] {
	edgePix := make([][]util.Point[int], 4)

	edgePix[0] = northPix(rect)
	edgePix[1] = eastPix(rect)
	edgePix[2] = southPix(rect)
	edgePix[3] = westPix(rect)

	return edgePix
}

func northPix(rect image.Rectangle) []util.Point[int] {
	pix := make([]util.Point[int], rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = util.Point[int]{X: i, Y: 0}
	}

	return pix
}

func southPix(rect image.Rectangle) []util.Point[int] {
	pix := make([]util.Point[int], rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = util.Point[int]{X: i, Y: rect.Max.X - 1}
	}

	return pix
}

func westPix(rect image.Rectangle) []util.Point[int] {
	pix := make([]util.Point[int], rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = util.Point[int]{X: 0, Y: i}
	}

	return pix
}

func eastPix(rect image.Rectangle) []util.Point[int] {
	pix := make([]util.Point[int], rect.Max.X)

	for i := 0; i < rect.Max.X; i++ {
		pix[i] = util.Point[int]{X: rect.Max.X - 1, Y: i}
	}

	return pix
}
