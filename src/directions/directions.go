package directions

type Direction string

const (
	North Direction = "north"
	East            = "east"
	South           = "south"
	West            = "west"
)

const (
	IntNorth int = 0
	IntEast      = 1
	IntSouth     = 2
	IntWest      = 3
)

var IntMap = map[int]Direction{
	0: North,
	1: East,
	2: South,
	3: West,
}

var StrMap = map[Direction]int{
	North: 0,
	East:  1,
	South: 2,
	West:  3,
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
