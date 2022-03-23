package loader

import (
	"beeb/carcassonne/directions"
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

	info := TileInfoFile{}

	fileContent, err := os.ReadFile(infoFileName)

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(fileContent, &info)

	if err != nil {
		panic(err)
	}

	return info
}

func LoadTiles() (map[string]tile.Tile, TileInfoFile) {
	fmt.Println("Loading Tiles")

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
			Name:       tileName,
			Image:      bitmaps[tileInfo.Image],
			Features:   tileInfo.Features,
			Edges:      tileInfo.Edges,
			Neighbours: make(map[directions.Direction]*tile.Tile),
			Placement: tile.Placement{
				Position:    tile.Position{},
				Orientation: 0,
			},
		}
	}

	return tiles, info
}

type RiverDeck struct {
	Begin string
	End   string
	Deck  map[string]int
}

type TileInfoFile struct {
	Tiles     map[string]TileInfo
	Deck      map[string]int
	RiverDeck RiverDeck
}

type TileInfo struct {
	Image    string
	Features map[int]tile.Feature
	Edges    map[directions.Direction]int
}
