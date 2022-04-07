package board

import "beeb/carcassonne/tile"

type Road struct {
	Segments []*tile.RoadSegment
}

//compile a whole road from any given segement of the road
func CompileRoadFromSegment(rs *tile.RoadSegment) Road {

	if rs == nil {
		panic("Send me a road segment damnit")
	}

	r := Road{}

	stack := make([]*tile.RoadSegment, 0, 1)

	//push first element to the stack
	stack = append(stack, rs)

	var stack_rs *tile.RoadSegment

	//while there's something in the stack to look at
	for len(stack) > 0 {

		//pop off the stack
		stack_rs, stack = stack[len(stack)-1], stack[:len(stack)-1]

		if r.ContainsSegment(stack_rs) {
			continue
		}

		//add the stack segment to the road
		r.Segments = append(r.Segments, stack_rs)

		//look at each edge of the tile that is connected to this road segment
		for _, s := range stack_rs.EdgeSegments {

			//if the edge is not connected, continue
			if s == nil {
				continue
			}

			//if the road is connected to another segment, add it to the stack
			stack = append(stack, s)
		}
	}

	return r
}

func (r Road) ContainsSegment(segment *tile.RoadSegment) bool {
	for _, rs := range r.Segments {
		if rs == segment {
			return true
		}
	}

	return false
}
