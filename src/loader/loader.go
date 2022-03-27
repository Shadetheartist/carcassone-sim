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

type FeatureInfo struct {
	Type   string
	Shield bool
	Curve  bool
}

type TileInfo struct {
	Image    string
	Features map[int]FeatureInfo
	Edges    map[string]int
}

func loadBitmaps(bitmapDir string) map[string]image.Image {
	files, err := ioutil.ReadDir(bitmapDir)

	if err != nil {
		panic(err)
	}

	bitmaps := make(map[string]image.Image)

	for _, file := range files {

		fileName := filepath.Join(bitmapDir, file.Name())

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

func LoadTiles(ymlPath string, bitmapDirectory string) (map[string]tile.Tile, TileInfoFile) {
	bitmaps := loadBitmaps(bitmapDirectory)
	info := loadTileInfo(ymlPath)

	tiles := make(map[string]tile.Tile)

	for tileName, tileInfo := range info.Tiles {

		if _, ok := bitmaps[tileInfo.Image]; !ok {
			panic(fmt.Sprint("Bitmap not found", tileInfo.Image))
		}

		edges := make(map[directions.Direction]int)

		for edgeStr, featureId := range tileInfo.Edges {
			edges[directions.StrMap[edgeStr]] = featureId
		}

		features := make(map[int]*tile.Feature)

		for featureId, featureInfo := range tileInfo.Features {

			edgesForFeature := make([]directions.Direction, 0)

			for dir, fid := range edges {
				if fid == featureId {
					edgesForFeature = append(edgesForFeature, dir)
				}
			}

			features[featureId] = &tile.Feature{
				Type:   tile.FeatureTypeStrMap[featureInfo.Type],
				Shield: featureInfo.Shield,
				Curve:  featureInfo.Curve,
				Edges:  edgesForFeature,
			}
		}

		t := tile.Tile{
			Name:       tileName,
			Image:      bitmaps[tileInfo.Image],
			Features:   features,
			Edges:      edges,
			Neighbours: make([]*tile.Tile, 4),
			Placement: tile.Placement{
				Position:    tile.Position{},
				Orientation: 0,
			},
		}

		t.EdgeFeatureTypes = t.ComputeEdgeFeatureTypes()

		tiles[tileName] = t

	}

	return tiles, info
}
