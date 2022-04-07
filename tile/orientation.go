package tile

var OrientationList = [...]uint16{0, 90, 180, 270}

func LimitToOrientation(n int) uint16 {
	if n < 0 {
		panic("Was not expecting a negative number")
	}

	return uint16((n * 90) % 360)
}
