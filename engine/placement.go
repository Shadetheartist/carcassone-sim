package engine

import (
	"beeb/carcassonne/engine/tile"
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

type PlacementFunction func(e *Engine, rtg *tile.ReferenceTileGroup) []Placement

func (e *Engine) PossibleTilePlacements(rtg *tile.ReferenceTileGroup) []Placement {
	return e.placementFunction(e, rtg)
}

func possibleTilePlacementsNonDeterministic(e *Engine, rtg *tile.ReferenceTileGroup) []Placement {

	if e.GameBoard.PlacedTileCount < 1 {
		return e.defaultPlacements(rtg)
	}

	e.placementBuffer = e.placementBuffer[:0]
	//all connections will share this one buffer, they will be sub-slices
	e.connectionsBuffer = e.connectionsBuffer[:0]

	lastConnectionLen := 0
	for openPos := range e.GameBoard.OpenPositions {

		//this will fill the connection buffer
		e.getPlaceableOrientations(openPos, rtg)

		for _, rt := range e.orientationBuffer {
			if rt != nil {
				e.placementBuffer = append(e.placementBuffer, Placement{
					Position:          openPos,
					ReferenceTile:     rt,
					ConnectedFeatures: e.connectionsBuffer[lastConnectionLen:len(e.connectionsBuffer)],
				})
			}
		}

		lastConnectionLen = len(e.connectionsBuffer)
	}

	return e.placementBuffer
}

func possibleTilePlacementsDeterministic(e *Engine, rtg *tile.ReferenceTileGroup) []Placement {

	if e.GameBoard.PlacedTileCount < 1 {
		return e.defaultPlacements(rtg)
	}

	e.placementBuffer = e.placementBuffer[:0]

	//all connections will share this one buffer, they will be sub-slices
	e.connectionsBuffer = e.connectionsBuffer[:0]

	lastConnectionLen := 0

	for _, openPos := range e.GameBoard.OpenPositionsList {

		//this will fill the connection buffer
		e.getPlaceableOrientations(openPos, rtg)

		for _, rt := range e.orientationBuffer {
			if rt != nil {
				e.placementBuffer = append(e.placementBuffer, Placement{
					Position:          openPos,
					ReferenceTile:     rt,
					ConnectedFeatures: e.connectionsBuffer[lastConnectionLen:len(e.connectionsBuffer)],
				})
			}
		}

		lastConnectionLen = len(e.connectionsBuffer)
	}

	return e.placementBuffer
}

func (e *Engine) getPlaceableOrientations(openPosKey util.Point[int], rtg *tile.ReferenceTileGroup) {
	openPositionEdgeSignature := e.GameBoard.OpenPositions[openPosKey]

	for i := 0; i < 4; i++ {
		rt := rtg.Orientations[i]
		tileEdgeSignature := rt.EdgeSignature
		if tileEdgeSignature.Compatible(openPositionEdgeSignature) {
			e.orientationBuffer[i] = rt

			for edge, feature := range rt.EdgeFeatures {
				edgeDir := directions.Direction(edge)
				neighbourPos := openPosKey.EdgePos(edgeDir)
				otherTile, err := e.GameBoard.TileMatrix.GetPt(neighbourPos)

				if err != nil {
					continue
				}

				complimentDir := directions.Compliment[edgeDir]

				if otherTile != nil {
					e.connectionsBuffer = append(e.connectionsBuffer, Connection{
						FeatureA: feature,
						EdgeA:    edgeDir,
						FeatureB: otherTile.EdgeFeatures[complimentDir],
						EdgeB:    complimentDir,
					})
				}
			}
		} else {
			e.orientationBuffer[i] = nil
		}
	}
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

//eek!

type Evaluation struct {
	Score  int
	Meeple int
}

func EvaluatePlacement(placement Placement) Evaluation {
	return Evaluation{}
}
