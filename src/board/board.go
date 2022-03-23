package board

import (
	dir "beeb/carcassonne/directions"
	"beeb/carcassonne/tile"
	"errors"
	"image"
)

type Board struct {
	Tiles map[tile.Position]*tile.Tile
}

func CreateBoard() Board {
	b := Board{}
	b.Tiles = make(map[tile.Position]*tile.Tile)

	return b
}

func (b Board) CompareEdge(t *tile.Tile, p tile.Placement, d dir.Direction) (*tile.Tile, error) {
	tileDir := p.TileDirection(d)
	neighbourPos := p.Position.EdgePos(d)

	if neighbour, exists := b.Tiles[neighbourPos]; exists {
		neighbourComplimentDir := dir.Compliment[d]
		neighbourTileDir := neighbour.Placement.TileDirection(neighbourComplimentDir)

		if t.Feature(tileDir).Type != neighbour.Feature(neighbourTileDir).Type {
			return neighbour, errors.New(string(d) + " Feature Type Mismatch")
		}

		return neighbour, nil
	}

	return nil, nil
}

//returns a map of all valid neighbours
func (b Board) VerifyPlacement(t *tile.Tile, position tile.Placement) (map[dir.Direction]*tile.Tile, error) {

	neighbours := make(map[dir.Direction]*tile.Tile)

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

func (b Board) ConnectedFeatures(t *tile.Tile, p tile.Placement) map[dir.Direction]tile.FeatureType {
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

func (b Board) AddTile(t *tile.Tile, p tile.Placement) error {
	neighbours, err := b.VerifyPlacement(t, p)

	if err != nil {
		return err
	}

	//if there are no mismatched tiles,
	//set the neighbour relationships between the tiles for later use
	for _, d := range dir.List {
		if neighbours[d] != nil {
			t.Neighbours[d] = neighbours[d]
			neighbours[d].Neighbours[dir.Compliment[d]] = t
		}
	}

	//place the tile at the position for itself and the board
	t.Placement = p
	t.Neighbours = neighbours
	b.Tiles[p.Position] = t

	return nil
}

func (b Board) PossibleTilePlacements(t *tile.Tile) []tile.Placement {
	placments := make([]tile.Placement, 0)

	//for every tile on the board
	for boardTilePosition := range b.Tiles {

		//observe the orthagonal neighbours of the tile
		neighbours := b.OrthoganalNeighbours(boardTilePosition)

		for neighbourPosition, neighbour := range neighbours {
			//neighbour being nil means it's an open placement location
			if neighbour == nil {
				//try placing the piece in every orientation
				for _, orientation := range tile.OrientationList {

					//try placing
					_, err := b.VerifyPlacement(t, tile.Placement{
						Position:    neighbourPosition,
						Orientation: orientation,
					})

					//if there is no error we can validly place the piece here
					if err == nil {
						placments = append(placments, tile.Placement{
							Position:    neighbourPosition,
							Orientation: orientation,
						})
					}
				}
			}
		}
	}

	return placments
}

func (b Board) OrthoganalNeighbours(p tile.Position) map[tile.Position]*tile.Tile {
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

func (b Board) BoundingBox() image.Rectangle {
	rect := image.Rectangle{}

	for pos := range b.Tiles {
		rect.Max.X = max(pos.X, rect.Max.X)
		rect.Max.Y = max(pos.Y, rect.Max.Y)
		rect.Min.X = min(pos.X, rect.Min.X)
		rect.Min.Y = min(pos.X, rect.Min.X)
	}

	return rect
}
