package db

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFileDataLoader struct {
	info GameConfig
}

func (dl *ConfigFileDataLoader) GetAllTileNames() []string {
	keys := make([]string, 0, len(dl.info.Tiles))
	for k := range dl.info.Tiles {
		keys = append(keys, k)
	}

	return keys
}

func (dl *ConfigFileDataLoader) GetGameConfig() GameConfig {
	return dl.info
}

func (dl *ConfigFileDataLoader) GetTileInfo(tileName string) (TileInfo, error) {

	if dl.info.Tiles == nil {
		return TileInfo{}, fmt.Errorf("Tile data has not been loaded")
	}

	if t, exists := dl.info.Tiles[tileName]; exists {
		return t, nil
	}

	return TileInfo{}, fmt.Errorf("There is no tile with the name %s", tileName)
}

func (dl *ConfigFileDataLoader) LoadData(configFilePath string) error {

	info := GameConfig{}

	fileContent, err := os.ReadFile(configFilePath)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(fileContent, &info)

	if err != nil {
		return err
	}

	dl.info = info

	return nil
}
