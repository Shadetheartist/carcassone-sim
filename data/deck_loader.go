package data

import (
	"os"

	"gopkg.in/yaml.v2"
)

func LoadDeckInfo(deckFilePath string) (DeckInfo, error) {

	info := DeckInfo{}

	fileContent, err := os.ReadFile(deckFilePath)

	if err != nil {
		return info, err
	}

	err = yaml.Unmarshal(fileContent, &info)

	if err != nil {
		return info, err
	}

	return info, nil
}
