package tile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type TileData struct {
	Tiles map[string]Tile
	Deck  map[string]uint8
}

func LoadTiles(fileName string) TileData {

	fmt.Println("Loading Tile Data")

	tiles := TileData{}

	fileContent, err := os.ReadFile(fileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileContent, &tiles)

	if err != nil {
		panic(err)
	}

	fmt.Println("Tile Data Loaded")

	for k, v := range tiles.Tiles {
		fmt.Println(k, v)
	}

	return tiles

}
