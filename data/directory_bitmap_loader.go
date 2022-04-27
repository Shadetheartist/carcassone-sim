package data

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/image/bmp"
)

type DirectoryBitmapLoader struct {
	bitmaps map[string]image.Image
}

func (dbl *DirectoryBitmapLoader) Keys() []string {
	keys := make([]string, 0, len(dbl.bitmaps))

	for k := range dbl.bitmaps {
		keys = append(keys, k)
	}

	return keys
}

func (dbl *DirectoryBitmapLoader) GetTileBitmap(tileName string) (image.Image, error) {

	if dbl.bitmaps == nil {
		return nil, fmt.Errorf("image data has not been loaded")
	}

	if img, exists := dbl.bitmaps[tileName]; exists {
		return img, nil
	}

	return nil, fmt.Errorf("there is no bitmap with the name %s", tileName)
}

func (dbl *DirectoryBitmapLoader) LoadBitmapsFromDirectory(bitmapDir string) {
	files, err := ioutil.ReadDir(bitmapDir)

	if err != nil {
		panic(err)
	}

	bitmaps := make(map[string]image.Image)

	for _, file := range files {

		fileName := filepath.Join(bitmapDir, file.Name())
		tileName := strings.Split(file.Name(), ".bmp")[0]

		reader, err := os.Open(fileName)

		if err != nil {
			panic(err)
		}

		image, err := bmp.Decode(reader)

		if err != nil {
			panic(err)
		}

		bitmaps[tileName] = image
	}

	dbl.bitmaps = bitmaps
}
