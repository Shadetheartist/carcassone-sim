package data

import (
	"beeb/carcassonne/engine/tile"
	"beeb/carcassonne/imageHelpers"
	"beeb/carcassonne/matrix"
	"beeb/carcassonne/util"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

type DeckInfo struct {
	Deck map[string]int
}

type GameData struct {
	TileNames []string
	Bitmaps   map[string]image.Image
	DeckInfo  DeckInfo

	//mapped by name, then by orientation (0 = 0, 1 = 90, 2 = 180, 3 = 270 degrees)
	ReferenceTileGroups map[string]*tile.ReferenceTileGroup
}

func LoadGameData(bitmapDirectory string, deckFilePath string) *GameData {
	gameData := &GameData{}

	gameData.loadBitmaps(bitmapDirectory)
	gameData.loadDeckInfo(deckFilePath)
	gameData.compileReferenceTiles()

	return gameData
}

func (gd *GameData) loadDeckInfo(deckFilePath string) {
	di, err := LoadDeckInfo(deckFilePath)

	if err != nil {
		panic(err)
	}

	for deckFileTileName, _ := range di.Deck {
		bitmapExists := false
		for _, bitmapTileName := range gd.TileNames {
			if deckFileTileName == bitmapTileName {
				bitmapExists = true
				break
			}
		}

		if !bitmapExists {
			panic(fmt.Sprint("Missing bmp file for deck tile: ", deckFileTileName))
		}
	}

	gd.DeckInfo = di
}

func (gd *GameData) loadBitmaps(bitmapDirectory string) {
	bitmapLoader := DirectoryBitmapLoader{}
	bitmapLoader.LoadBitmapsFromDirectory(bitmapDirectory)

	gd.Bitmaps = bitmapLoader.bitmaps
	gd.TileNames = bitmapLoader.Keys()

	sort.SliceStable(gd.TileNames, func(i, j int) bool {
		return gd.TileNames[i] < gd.TileNames[j]
	})
}

func (gd *GameData) compileReferenceTiles() {
	gd.ReferenceTileGroups = make(map[string]*tile.ReferenceTileGroup, len(gd.TileNames))

	for _, tileName := range gd.TileNames {
		img := gd.Bitmaps[tileName]

		orientedReferenceTiles := make([]*tile.ReferenceTile, 4)
		compiledReferenceTile := gd.compileTile(tileName)
		compiledReferenceTile.Image = img
		orientedReferenceTiles[0] = &compiledReferenceTile

		rotated90 := compiledReferenceTile
		rotated90.Orientation = 90
		rotated90.FeatureMatrix = compiledReferenceTile.FeatureMatrix.Copy()
		rotated90.FeatureMatrix.Rotate90()
		//this rotates counter-clockwise, so we use the inverse
		rotated90.Image = imaging.Rotate270(img)

		orientedReferenceTiles[1] = &rotated90

		rotated180 := compiledReferenceTile
		rotated180.Orientation = 180
		rotated180.FeatureMatrix = compiledReferenceTile.FeatureMatrix.Copy()
		rotated180.FeatureMatrix.Rotate180()
		rotated180.Image = imaging.Rotate180(img)
		orientedReferenceTiles[2] = &rotated180

		rotated270 := compiledReferenceTile
		rotated270.Orientation = 270
		rotated270.FeatureMatrix = compiledReferenceTile.FeatureMatrix.Copy()
		rotated270.FeatureMatrix.Rotate270()
		//this rotates counter-clockwise, so we use the inverse
		rotated270.Image = imaging.Rotate90(img)
		orientedReferenceTiles[3] = &rotated270

		rtg := &tile.ReferenceTileGroup{}
		rtg.Orientations = orientedReferenceTiles
		rtg.Name = tileName
		gd.ReferenceTileGroups[tileName] = rtg
	}

	for _, tileName := range gd.TileNames {
		rtg := gd.ReferenceTileGroups[tileName]

		compileAvgFeaturePositions(rtg)

		for i := 0; i < 4; i++ {
			rt := rtg.Orientations[i]
			rt.EdgeFeatures = determineEdgeFeatures(rt.FeatureMatrix)
			rt.EdgeSignature = compileEdgeSignature(rt.EdgeFeatures)
		}

		rtg.Features = rtg.Orientations[0].Features
		for _, f := range rtg.Features {
			f.ParentRefenceTileGroup = rtg
		}
	}

}

func compileAvgFeaturePositions(rtg *tile.ReferenceTileGroup) {
	for _, rt := range rtg.Orientations {
		rt.AvgFeaturePos = make(map[*tile.Feature]util.Point[float64])
		for _, f := range rt.Features {
			rt.AvgFeaturePos[f] = avgFeaturePos(f, rt.FeatureMatrix)
		}
	}
}

func avgFeaturePos(f *tile.Feature, m *matrix.Matrix[*tile.Feature]) util.Point[float64] {
	var totalX int
	var totalY int
	var n float64

	m.Iterate(func(mf *tile.Feature, x int, y int, idx int) {
		if mf == f {
			totalX += x
			totalY += y
			n++
		}
	})

	//no dividing by zero!
	if n == 0 {
		return util.Point[float64]{X: -1, Y: -1}
	}

	avgX := float64(totalX) / n
	avgY := float64(totalY) / n

	return util.Point[float64]{X: avgX, Y: avgY}
}

func compileEdgeSignature(edgeFeatures *tile.EdgeArray[*tile.Feature]) *tile.EdgeSignature {
	sig := &tile.EdgeSignature{}

	for i, f := range edgeFeatures {
		sig[i] = f.Type
	}

	return sig
}

func (gd *GameData) compileTile(tileName string) tile.ReferenceTile {
	img := gd.Bitmaps[tileName]

	rt := tile.ReferenceTile{}

	rt.Name = tileName

	rt.FeatureMatrix, rt.Features = gd.buildMatrix(img)

	return rt
}

func isRoadColor(c color.Color) bool {
	r, g, b, a := c.RGBA()
	roadR, roadG, roadB, roadA := color.RGBA{R: 255, G: 255, B: 255, A: 255}.RGBA()
	return r == roadR && g == roadG && b == roadB && a == roadA
}

func isFarmColor(c color.Color) bool {
	r, g, b, a := c.RGBA()
	farmR, farmG, farmB, farmA := color.RGBA{R: 106, G: 190, B: 48, A: 255}.RGBA()
	return r == farmR && g == farmG && b == farmB && a == farmA
}

func isRiverColor(c color.Color) bool {
	r, g, b, a := c.RGBA()
	riverR, riverG, riverB, riverA := color.RGBA{R: 91, G: 110, B: 225, A: 255}.RGBA()
	return r == riverR && g == riverG && b == riverB && a == riverA
}

func isCastleDark(c color.Color) bool {
	r, g, b, a := c.RGBA()
	castleR, castleG, castleB, castleA := color.RGBA{R: 102, G: 57, B: 49, A: 255}.RGBA()
	return r == castleR && g == castleG && b == castleB && a == castleA
}

func isCloisterColor(c color.Color) bool {
	r, g, b, a := c.RGBA()

	cloisterR, cloisterG, cloisterB, cloisterA := color.RGBA{R: 63, G: 63, B: 116, A: 255}.RGBA()
	cloisterR2, cloisterG2, cloisterB2, cloisterA2 := color.RGBA{R: 50, G: 60, B: 57, A: 255}.RGBA()

	return (r == cloisterR || r == cloisterR2) &&
		(g == cloisterG || g == cloisterG2) &&
		(b == cloisterB || b == cloisterB2) &&
		(a == cloisterA || a == cloisterA2)
}

func isShieldColor(c color.Color) bool {
	r, g, b, a := c.RGBA()
	riverR, riverG, riverB, riverA := color.RGBA{R: 99, G: 155, B: 255, A: 255}.RGBA()
	return r == riverR && g == riverG && b == riverB && a == riverA
}

func (gd *GameData) buildMatrix(img image.Image) (*matrix.Matrix[*tile.Feature], []*tile.Feature) {

	featureMatrix := matrix.NewMatrix[*tile.Feature](img.Bounds().Dx())

	features := make([]*tile.Feature, 0, 4)

	var feature *tile.Feature
	var featureColor color.Color
	var fillOrthoganallyOnly bool

	//there's only ever one river segment per tile
	//however sometimes it's cut in half by a road,
	//so we have to have special handling to only use one river segment at any time
	var riverFeature *tile.Feature

	segmentCallback := func(img image.Image, p image.Point, idx int) bool {
		featureColor = img.At(p.X, p.Y)
		feature = &tile.Feature{}
		feature.Id = uuid.New()

		fillOrthoganallyOnly = true

		if isRoadColor(featureColor) {
			feature.Type = tile.Road
			fillOrthoganallyOnly = false
		} else if isRiverColor(featureColor) {
			if riverFeature == nil {
				riverFeature = &tile.Feature{}
			}

			feature = riverFeature
			feature.Type = tile.River
			fillOrthoganallyOnly = false
		} else if isFarmColor(featureColor) {
			feature.Type = tile.Farm
		} else if isCastleDark(featureColor) {
			feature.Type = tile.Castle
		} else if isShieldColor(featureColor) {
			feature.Type = tile.Shield
		} else if isCloisterColor(featureColor) {
			feature.Type = tile.Cloister
		} else {
			feature = nil
		}

		if feature != nil {
			features = append(features, feature)
		}

		//return true if orthogonal neighbour behaviour is desired
		return fillOrthoganallyOnly
	}

	fillProcessor := func(img image.Image, p image.Point, idx int) bool {
		c := img.At(p.X, p.Y)

		sameColor := c == featureColor || (isCloisterColor(c) && isCloisterColor(featureColor))

		if sameColor {
			featureMatrix.Set(p.X, p.Y, feature)
		}

		return sameColor
	}

	imageHelpers.FillAll(img, segmentCallback, fillProcessor)

	return featureMatrix, features
}

func determineEdgeFeatures(featureMatrix *matrix.Matrix[*tile.Feature]) *tile.EdgeArray[*tile.Feature] {
	centerPixel := 3
	size := featureMatrix.Size()

	edgeArray := &tile.EdgeArray[*tile.Feature]{}

	edgeArray.SetNorth(featureMatrix.Get(centerPixel, 0))
	edgeArray.SetSouth(featureMatrix.Get(centerPixel, size-1))
	edgeArray.SetWest(featureMatrix.Get(0, centerPixel))
	edgeArray.SetEast(featureMatrix.Get(size-1, centerPixel))

	return edgeArray
}
