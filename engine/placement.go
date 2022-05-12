package engine

import (
	"beeb/carcassonne/tile"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
	"math/rand"
)

type Connection struct {
	FeatureA *tile.Feature
	EdgeA    directions.Direction
	FeatureB *tile.Feature
	EdgeB    directions.Direction
}

type Placement struct {
	Position          util.Point[int]
	ReferenceTile     *tile.ReferenceTile
	ConnectedFeatures []Connection
}

func (e *Engine) PossibleTilePlacements(rtg *tile.ReferenceTileGroup) []Placement {

	if e.GameBoard.PlacedTileCount < 1 {
		return e.defaultPlacements(rtg)
	}

	placements := make([]Placement, 0, len(e.GameBoard.OpenPositions))

	//reused when checking placable orientations of each position
	orientationBuffer := make([]*tile.ReferenceTile, 4)

	for openPos := range e.GameBoard.OpenPositions {

		if openPos.X < 0 || openPos.X >= e.GameBoard.TileMatrix.Size() {
			continue
		}

		if openPos.Y < 0 || openPos.Y >= e.GameBoard.TileMatrix.Size() {
			continue
		}

		connections := e.getPlaceableOrientations(orientationBuffer, openPos, rtg)

		for _, rt := range orientationBuffer {
			if rt != nil {
				placements = append(placements, Placement{
					Position:          openPos,
					ReferenceTile:     rt,
					ConnectedFeatures: connections,
				})
			}
		}
	}

	return placements
}

func (e *Engine) getPlaceableOrientations(buffer []*tile.ReferenceTile, openPosKey util.Point[int], rtg *tile.ReferenceTileGroup) []Connection {
	openPositionEdgeSignature := e.GameBoard.OpenPositions[openPosKey]
	connections := make([]Connection, 0, 4)

	for i := 0; i < 4; i++ {
		rt := rtg.Orientations[i]
		tileEdgeSignature := rt.EdgeSignature
		if tileEdgeSignature.Compatible(openPositionEdgeSignature) {
			buffer[i] = rt

			for edge, feature := range rt.EdgeFeatures {
				edgeDir := directions.Direction(edge)
				neighbourPos := openPosKey.EdgePos(edgeDir)
				otherTile, err := e.GameBoard.TileMatrix.GetPt(neighbourPos)

				if err != nil {
					continue
				}

				complimentDir := directions.Compliment[edgeDir]

				if otherTile != nil {
					connections = append(connections, Connection{
						FeatureA: feature,
						EdgeA:    edgeDir,
						FeatureB: otherTile.EdgeFeatures[complimentDir],
						EdgeB:    complimentDir,
					})
				}
			}

		} else {
			buffer[i] = nil
		}
	}

	return connections
}

func (e *Engine) defaultPlacements(rtg *tile.ReferenceTileGroup) []Placement {
	placements := make([]Placement, 0, len(e.GameBoard.OpenPositions))

	middle := e.GameBoard.TileMatrix.Size() / 2
	for i := 0; i < 4; i++ {
		placements = append(placements,
			Placement{
				Position: util.Point[int]{
					X: middle,
					Y: middle,
				},
				ReferenceTile:     rtg.Orientations[i],
				ConnectedFeatures: make([]Connection, 0),
			},
		)
	}

	return placements
}

func RandomPlacement(placements []Placement) *Placement {
	if len(placements) == 0 {
		return nil
	}

	randN := rand.Int() % len(placements)
	return &placements[randN]
}
