package deck

import (
	"beeb/carcassonne/engine/tile"
	"errors"
	"math/rand"
)

type Deck struct {
	Tiles []*tile.ReferenceTileGroup
}

func (d *Deck) Scry() (*tile.ReferenceTileGroup, error) {

	if len(d.Tiles) < 1 {
		return nil, errors.New("Deck is Empty")
	}

	return d.Tiles[0], nil
}

func (d *Deck) Pop() (*tile.ReferenceTileGroup, error) {

	if len(d.Tiles) < 1 {
		return nil, errors.New("Deck is Empty")
	}

	tile := d.Tiles[0]

	d.Tiles = d.Tiles[1:]

	return tile, nil
}

func (d *Deck) Remaining() int {
	return len(d.Tiles)
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.Tiles), func(i, j int) {
		d.Tiles[i], d.Tiles[j] = d.Tiles[j], d.Tiles[i]
	})
}

func (d *Deck) Prepend(t *tile.ReferenceTileGroup) {
	d.Tiles = append([]*tile.ReferenceTileGroup{t}, d.Tiles...)
}

func (d *Deck) Append(t *tile.ReferenceTileGroup) {
	d.Tiles = append(d.Tiles, t)
}
