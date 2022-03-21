package main

import (
	"beeb/carcassonne/loader"
	"fmt"

	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
)

var size float64 = 100

func main() {
	tiles := loader.LoadTiles()

	fmt.Println("Loaded Tiles", tiles)

	dc := gg.NewContext(100, 100)

	var x int = 0
	var scale float64 = 2

	for _, v := range tiles {
		scaledImageSize := int(5 * scale)
		rescaledImage := resize.Resize(
			uint(scaledImageSize),
			uint(scaledImageSize),
			v.Image,
			resize.NearestNeighbor,
		)

		dci := gg.NewContextForImage(rescaledImage)
		dc.DrawImage(dci.Image(), x, 0)
		dc.Fill()

		x += int(float64(scaledImageSize) * 1.2)
	}

	fmt.Println("Saving image")

	dc.SavePNG("out.png")

	//todo: the map,
	//tile placement validation,
	//linked features,
	//game structure,
	//player AI
}
