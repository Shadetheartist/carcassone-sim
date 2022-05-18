package board

import (
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
	"image"
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

func (b *Board) PlaceTile(pos util.Point[int], t *tile.Tile) {
	b.TileMatrix.Set(pos.X, pos.Y, t)
	b.PlacedTileCount++

	t.Position = pos

	//setup neighbours
	if tl, err := b.TileMatrix.GetPt(pos.North()); err == nil {
		t.Neighbours.SetNorth(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.East()); err == nil {
		t.Neighbours.SetEast(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.South()); err == nil {
		t.Neighbours.SetSouth(tl)
	}

	if tl, err := b.TileMatrix.GetPt(pos.West()); err == nil {
		t.Neighbours.SetWest(tl)
	}

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

		b.OpenPositionsList = append(b.OpenPositionsList[:i], b.OpenPositionsList[i+1:]...)
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
