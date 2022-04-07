package board

import (
	dir "beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"errors"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Board struct {
	TileOptions map[string]tile.Tile
	Tiles       map[tile.Position]*tile.Tile
	RoadCount   int

	//dir mapped
	OpenPositions map[tile.Position][]tile.FeatureType

	BoardImage         *ebiten.Image
	OpenPositionsImage *ebiten.Image
	RoadsImage         *ebiten.Image
}

func New(tileOptions map[string]tile.Tile, imageW int, imageH int) Board {
	b := Board{}
	b.TileOptions = tileOptions
	b.Tiles = make(map[tile.Position]*tile.Tile)
	b.OpenPositions = make(map[tile.Position][]tile.FeatureType, 4)

	b.BoardImage = ebiten.NewImage(imageW, imageH)
	b.RoadsImage = ebiten.NewImage(imageW, imageH)
	b.OpenPositionsImage = ebiten.NewImage(imageW, imageH)

	return b
}

//be careful not to add the same tile by reference accidentally,
//always make sure that each tile sent to this function has its own memory
func (b *Board) AddTile(t *tile.Tile, p tile.Placement) error {

	//verify placement
	if _, err := b.IsTilePlaceable(t, p.Position); err != nil {
		return errors.New("Can't add a tile here, it is not able to be placed.")
	}

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

	//track changes to open positions
	b.manageOpenPositions(t, p, neighbours)

	//track changes and additions to roads
	//this must occur after the other tile state changes as it
	//depends on the tile's neighbouring tiles
	t.IntegrateRoads()

	return nil
}

func (b *Board) manageOpenPositions(t *tile.Tile, p tile.Placement, neighbours []*tile.Tile) {
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
