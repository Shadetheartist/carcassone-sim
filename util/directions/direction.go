package directions

type Direction int

// specifially indexed to rotate together
const (
	North Direction = iota
	East
	South
	West
)

var StrMap = map[string]Direction{
	"north": North,
	"east":  East,
	"south": South,
	"west":  West,
}

var IntMap = map[Direction]string{
	North: "north",
	East:  "east",
	South: "south",
	West:  "west",
}

var Compliment = map[Direction]Direction{
	North: South,
	South: North,
	East:  West,
	West:  East,
}

var List = [...]Direction{
	North,
	East,
	South,
	West,
}

// Inner gets the two inward facing directions of corners of a square
// ex: for the top-left corner of a square, the inward directions are south and east
func Inner(top bool, left bool) []Direction {
	dirs := make([]Direction, 2)

	if !top {
		dirs[0] = South
	} else {
		dirs[0] = North
	}

	if !left {
		dirs[1] = East
	} else {
		dirs[1] = West
	}

	return dirs
}
