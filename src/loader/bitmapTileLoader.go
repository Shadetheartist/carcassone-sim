package loader

import (
	"beeb/carcassonne/tile"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/image/bmp"
	"gopkg.in/yaml.v2"
)

func loadBitmaps(bitmapDir string) map[string]image.Image {
	files, err := ioutil.ReadDir(bitmapDir)

	if err != nil {
		panic(err)
	}

	bitmaps := make(map[string]image.Image)

	for _, file := range files {

		fileName := filepath.Join(bitmapDir, file.Name())

		fmt.Println("Loading BMP", fileName)

		reader, err := os.Open(fileName)

		if err != nil {
			panic(err)
		}

		image, err := bmp.Decode(reader)

		if err != nil {
			panic(err)
		}

		bitmaps[file.Name()] = image

	}

	return bitmaps
}

func loadTileInfo(infoFileName string) TileInfoFile {

	tiles := TileInfoFile{}

	fileContent, err := os.ReadFile(infoFileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileContent, &tiles)

	if err != nil {
		panic(err)
	}

	return tiles

}

func LoadTiles() map[string]tile.Tile {
	bitmapDirectory := "../data/bitmaps"
	infoFileName := "../data/tiles.yml"

	bitmaps := loadBitmaps(bitmapDirectory)
	info := loadTileInfo(infoFileName)

	tiles := make(map[string]tile.Tile)

	for tileName, tileInfo := range info.Tiles {

		if _, ok := bitmaps[tileInfo.Image]; !ok {
			panic(fmt.Sprint("Bitmap not found", tileInfo.Image))
		}

		tiles[tileName] = tile.Tile{
			Name:        tileName,
			Orientation: tile.OrientationZero,
			Image:       bitmaps[tileInfo.Image],
		}

	}

	return tiles
}

type TileInfoFile struct {
	Tiles map[string]TileInfo
	Deck  map[string]int
}

type TileFeature struct {
	Id     int
	Type   string
	Shield bool
}

type TileInfo struct {
	Image    string
	Features []TileFeature
	Edges    map[string]int
}
