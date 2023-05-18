package aiLink

import (
	"beeb/carcassonne/engine"
	"beeb/carcassonne/engine/tile"
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sort"
)

type AILink struct {
	engine    *engine.Engine
	tileIndex map[string]map[int]int
}

func NewAILink(engine *engine.Engine) *AILink {

	// get unique types
	names := make([]string, 0, len(engine.GameData.DeckInfo.Deck))
	for name, _ := range engine.GameData.DeckInfo.Deck {
		names = append(names, name)
	}
	sort.Strings(names)

	tileIndex := make(map[string]map[int]int)
	for i := 0; i < len(names); i++ {
		tileIndex[names[i]] = make(map[int]int)
		tileIndex[names[i]][0] = i * 4
		tileIndex[names[i]][90] = i*4 + 1
		tileIndex[names[i]][180] = i*4 + 2
		tileIndex[names[i]][270] = i*4 + 3
	}

	return &AILink{
		engine:    engine,
		tileIndex: tileIndex,
	}
}

// totalInputs
// the number of inputs expected for this board state
// n: the number of possible positions to place a tile
func (ai *AILink) totalInputs(n int) int {
	numTileTypes := len(ai.engine.GameData.DeckInfo.Deck)
	numTileOrientations := 4
	return n * numTileTypes * numTileOrientations
}

// totalInputBytes
// same thing as total inputs, but calculates the number of bytes needed, for packing bits
func (ai *AILink) totalInputBytes(n int) int {
	return int(math.Ceil(float64(ai.totalInputs(n)) / 8))
}

// Inputs
// the board state as one-hot encoding
func (ai *AILink) Inputs() []byte {

	matrix := ai.engine.GameBoard.TileMatrix
	l := matrix.Len()

	bytesPerTile := ai.totalInputBytes(1)
	inputs := make([]byte, 0, bytesPerTile*l)

	for i := 0; i < l; i++ {
		t := matrix.GetI(i)
		tInputs := ai.inputsForTile(t)

		//printBytes(tInputs)

		inputs = append(inputs, tInputs...)
	}

	ai.saveImg(inputs)

	return inputs
}

func printBytes(bs []byte) {
	fmt.Println(hex.EncodeToString(bs))
}

func (ai *AILink) saveImg(data []byte) {
	matrix := ai.engine.GameBoard.TileMatrix
	l := matrix.Len()

	img := image.NewGray(image.Rect(0, 0, ai.totalInputs(1), l))

	stride := ai.totalInputs(1)
	for y := 0; y < l; y++ {
		for x := 0; x < stride; x++ {
			bitI := stride*y + x
			byteI := bitI / 8
			bitOff := bitI - byteI*8
			bitVal := (data[byteI] >> uint(7-bitOff)) & 1
			// Set black pixel for 0 bit and white pixel for 1 bit
			if bitVal == 0 {
				img.SetGray(x, y, color.Gray{Y: 0})
			} else {
				img.SetGray(x, y, color.Gray{Y: 255})
			}
		}
	}

	_ = saveImageAsPNG(img, "img.png")
}

func (ai *AILink) inputsForTile(t *tile.Tile) []byte {

	inputBytes := make([]byte, ai.totalInputBytes(1))

	if t == nil {
		return inputBytes
	}

	index, ok := ai.tileIndex[t.Reference.Name][t.Reference.Orientation]
	if !ok {
		panic("invalid orientation")
	}

	byteIdx := index / 8
	byteOffset := index - byteIdx*8
	val := byte(0x1) << byteOffset

	inputBytes[byteIdx] = val

	return inputBytes
}

func saveImageAsPNG(img image.Image, filename string) error {
	// Create a new file for writing
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the image as a PNG file
	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	fmt.Printf("PNG file '%s' successfully saved.\n", filename)
	return nil
}
