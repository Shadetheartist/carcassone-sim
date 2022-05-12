package deck

import (
	"beeb/carcassonne/tile"
	"errors"
	"math/rand"
)

type Deck struct {
	Index int
	Tiles []*tile.ReferenceTileGroup
}

func (d *Deck) Scry() (*tile.ReferenceTileGroup, error) {

	if d.Index >= len(d.Tiles) {
		return nil, errors.New("Deck is Empty")
	}

	tile := d.Tiles[d.Index]

	return tile, nil
}

func (d *Deck) Pop() (*tile.ReferenceTileGroup, error) {

	if d.Index >= len(d.Tiles) {
		return nil, errors.New("Deck is Empty")
	}

	tile := d.Tiles[d.Index]

	d.Index++

	return tile, nil
}

func (d *Deck) Remaining() int {
	return len(d.Tiles) - d.Index
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
