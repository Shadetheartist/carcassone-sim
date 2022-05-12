package tile

const EdgeArraySize int = 4

type EdgeArray[T comparable] [EdgeArraySize]T

func (edgeArray *EdgeArray[T]) GetNorth() T {
	return edgeArray[0]
}

func (edgeArray *EdgeArray[T]) GetEast() T {
	return edgeArray[1]
}

func (edgeArray *EdgeArray[T]) GetSouth() T {
	return edgeArray[2]
}

func (edgeArray *EdgeArray[T]) GetWest() T {
	return edgeArray[3]
}

func (edgeArray *EdgeArray[T]) SetNorth(t T) {
	edgeArray[0] = t
}

func (edgeArray *EdgeArray[T]) SetEast(t T) {
	edgeArray[1] = t
}

func (edgeArray *EdgeArray[T]) SetSouth(t T) {
	edgeArray[2] = t
}

func (edgeArray *EdgeArray[T]) SetWest(t T) {
	edgeArray[3] = t
}

func (edgeArray *EdgeArray[T]) Contains(t T) bool {
	for i := 0; i < 4; i++ {
		if edgeArray[i] == t {
			return true
		}
	}

	return false
}
