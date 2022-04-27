package imageHelpers

import "image"

func Neighbours(img image.Image, p image.Point) []image.Point {
	neighbours := make([]image.Point, 0, 8)

	neighbours = append(neighbours, OrthogonalNeighbours(img, p)...)

	//northwest
	if p.Y > 0 && p.X > 0 {
		neighbours = append(neighbours, image.Point{
			p.X - 1,
			p.Y - 1,
		})
	}

	//northeast
	if p.Y > 0 && p.X < img.Bounds().Dx()-1 {
		neighbours = append(neighbours, image.Point{
			p.X + 1,
			p.Y - 1,
		})
	}

	//southwest
	if p.Y < img.Bounds().Dy()-1 && p.X > 0 {
		neighbours = append(neighbours, image.Point{
			p.X - 1,
			p.Y + 1,
		})
	}

	//southeast
	if p.Y < img.Bounds().Dy()-1 && p.X < img.Bounds().Dx()-1 {
		neighbours = append(neighbours, image.Point{
			p.X + 1,
			p.Y + 1,
		})
	}

	return neighbours
}

func OrthogonalNeighbours(img image.Image, p image.Point) []image.Point {
	neighbours := make([]image.Point, 0, 4)

	//north
	if p.Y > 0 {
		neighbours = append(neighbours, image.Point{
			p.X,
			p.Y - 1,
		})
	}

	//west
	if p.X > 0 {
		neighbours = append(neighbours, image.Point{
			p.X - 1,
			p.Y,
		})
	}

	//south
	if p.Y < img.Bounds().Dy()-1 {
		neighbours = append(neighbours, image.Point{
			p.X,
			p.Y + 1,
		})
	}

	//east
	if p.X < img.Bounds().Dx()-1 {
		neighbours = append(neighbours, image.Point{
			p.X + 1,
			p.Y,
		})
	}

	return neighbours
}
