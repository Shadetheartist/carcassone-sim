package engine

import (
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/util"
	"beeb/carcassonne/util/directions"
	"math/rand"
	"sync"
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

type TilePlacementManager struct {
	engine       *Engine
	agents       []*TilePlacementAgent
	outputBuffer []Placement
}

func NewTilePlacementManager(e *Engine) *TilePlacementManager {

	tpm := &TilePlacementManager{}
	tpm.engine = e

	agentCount := 1

	tpm.outputBuffer = make([]Placement, 0, 128)

	tpm.agents = make([]*TilePlacementAgent, agentCount)

	for i := 0; i < agentCount; i++ {
		tpm.agents[i] = NewTilePlacementAgent(e)
	}

	return tpm
}

func (tpm *TilePlacementManager) PossibleTilePlacements(rtg *tile.ReferenceTileGroup) []Placement {

	if tpm.engine.GameBoard.PlacedTileCount < 1 {
		return tpm.firstTilePlacement(rtg)
	}

	return tpm.agents[0].PossibleTilePlacements(nil, rtg, tpm.engine.GameBoard.OpenPositionsList())
}

type TilePlacementAgent struct {
	engine            *Engine
	placementBuffer   []Placement
	orientationBuffer []*tile.ReferenceTile
	connectionsBuffer []Connection
}

func NewTilePlacementAgent(e *Engine) *TilePlacementAgent {
	tpa := &TilePlacementAgent{}

	tpa.engine = e
	tpa.orientationBuffer = make([]*tile.ReferenceTile, 4)
	tpa.placementBuffer = make([]Placement, 0, 256)
	tpa.connectionsBuffer = make([]Connection, 0, 256)

	return tpa
}

func (tpa *TilePlacementAgent) PossibleTilePlacements(wg *sync.WaitGroup, rtg *tile.ReferenceTileGroup, openPositionsList []util.Point[int]) []Placement {
	if wg != nil {
		defer wg.Done()
	}

	tpa.placementBuffer = tpa.placementBuffer[:0]

	//all connections will share this one buffer, they will be sub-slices
	tpa.connectionsBuffer = tpa.connectionsBuffer[:0]

	lastConnectionLen := 0

	for _, openPos := range openPositionsList {

		//this will fill the connection buffer
		tpa.getPlaceableOrientations(openPos, rtg)

		for _, rt := range tpa.orientationBuffer {
			if rt != nil {
				tpa.placementBuffer = append(tpa.placementBuffer, Placement{
					Position:          openPos,
					ReferenceTile:     rt,
					ConnectedFeatures: tpa.connectionsBuffer[lastConnectionLen:len(tpa.connectionsBuffer)],
				})
			}
		}

		lastConnectionLen = len(tpa.connectionsBuffer)
	}

	return tpa.placementBuffer
}

func (tpa *TilePlacementAgent) getPlaceableOrientations(openPosKey util.Point[int], rtg *tile.ReferenceTileGroup) {
	e := tpa.engine

	openPositionEdgeSignature := e.GameBoard.OpenPositions[openPosKey]

	for i := 0; i < 4; i++ {
		rt := rtg.Orientations[i]
		tileEdgeSignature := rt.EdgeSignature
		if tileEdgeSignature.Compatible(openPositionEdgeSignature) {
			tpa.orientationBuffer[i] = rt

			for edge, feature := range rt.EdgeFeatures {
				edgeDir := directions.Direction(edge)
				neighbourPos := openPosKey.EdgePos(edgeDir)
				otherTile, err := e.GameBoard.TileMatrix.GetPt(neighbourPos)

				if err != nil {
					continue
				}

				complimentDir := directions.Compliment[edgeDir]

				if otherTile != nil {
					tpa.connectionsBuffer = append(tpa.connectionsBuffer, Connection{
						FeatureA: feature,
						EdgeA:    edgeDir,
						FeatureB: otherTile.EdgeFeatures[complimentDir],
						EdgeB:    complimentDir,
					})
				}
			}
		} else {
			tpa.orientationBuffer[i] = nil
		}
	}
}

func (tpm *TilePlacementManager) firstTilePlacement(rtg *tile.ReferenceTileGroup) []Placement {
	e := tpm.engine
	placements := make([]Placement, 0, len(e.GameBoard.OpenPositions))

	middle := e.GameBoard.TileMatrix.Size() / 2
	quarter := e.GameBoard.TileMatrix.Size() / 4
	// for each quadrant to start it
	for i := 0; i < 4; i++ {
		top := i < 2     //true, true, false, false
		left := i%2 == 0 //true, false, true, false

		y := middle
		if top {
			y -= quarter
		} else {
			y += quarter
		}

		x := middle
		if left {
			x -= quarter
		} else {
			x += quarter
		}

		dirs := directions.Inner(top, left)
		for _, d := range dirs {
			placements = append(placements,
				Placement{
					Position: util.Point[int]{
						X: x,
						Y: y,
					},
					ReferenceTile:     rtg.Orientations[d],
					ConnectedFeatures: make([]Connection, 0),
				},
			)
		}

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
