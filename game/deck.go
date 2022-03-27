package game

import (
	"beeb/carcassonne/tile"
	"errors"
	"math/rand"
)

type Deck struct {
	Index int
	Tiles []tile.Tile
}

func (d *Deck) Pop() (tile.Tile, error) {

	if d.Index >= len(d.Tiles) {
		return tile.Tile{}, errors.New("Deck is Empty")
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

func (d *Deck) Prepend(t tile.Tile) {
	d.Tiles = append([]tile.Tile{t}, d.Tiles...)
}

func (d *Deck) Append(t tile.Tile) {
	d.Tiles = append(d.Tiles, t)
}
