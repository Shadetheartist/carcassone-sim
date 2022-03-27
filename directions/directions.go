package directions

type Direction int

//specifially indexed to rotate together
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
