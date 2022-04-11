package main

import (
	"beeb/carcassonne/db"
	"beeb/carcassonne/game"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var size float64 = 100

func ymlConfigFilePath() string {
	return filepath.Join(exeDir(), "data/tiles.yml")
}

func bitmapDirectory() string {
	return filepath.Join(exeDir(), "data/bitmaps")
}

func exeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	return filepath.Dir(exePath)
}

func main() {
	rand.Seed(time.Now().Unix())

	ebiten.SetWindowSize(1200, 900)
	ebiten.SetWindowTitle("Carcassonne Simulator")
	ebiten.SetScreenClearedEveryFrame(false)

	tileInfoLoader := &db.ConfigFileDataLoader{}
	configPath := ymlConfigFilePath()
	tileInfoLoader.LoadData(configPath)

	bitmapLoader := &db.DirectoryBitmapLoader{}
	bitmapDir := bitmapDirectory()
	bitmapLoader.LoadBitmapsFromDirectory(bitmapDir)

	game := game.CreateGame(tileInfoLoader, tileInfoLoader, bitmapLoader)

	game.Setup()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
