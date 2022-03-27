package tile_test

import (
	"beeb/carcassonne/directions"
	"beeb/carcassonne/loader"
	"beeb/carcassonne/tile"
	"testing"
)

func loadTiles() map[string]tile.Tile {
	tiles, _ := loader.LoadTiles("../data/tiles.yml", "../data/bitmaps")
	return tiles
}

func TestFeature(t *testing.T) {
	tiles := loadTiles()

	_tile := tiles["RiverCurve"]

	if _tile.Feature(directions.North).Type != tile.Grass {
		t.Error("Err N")
	}

	if _tile.Feature(directions.East).Type != tile.Grass {
		t.Error("Err E")
	}

	if _tile.Feature(directions.South).Type != tile.River {
		t.Error("Err S")
	}

	if _tile.Feature(directions.West).Type != tile.River {
		t.Error("Err W")
	}

	_tile.Placement.Orientation = 90

	g2td := _tile.Placement.GridToTileDir(directions.North)
	if _tile.Feature(g2td).Type != tile.River {
		t.Error("Err grid to tile dir N")
	}

	t2gd := _tile.Placement.TileToGridDir(directions.North)
	if _tile.Feature(t2gd).Type != tile.Grass {
		t.Error("Err tile to grid dir N")
	}

	g2td = _tile.Placement.GridToTileDir(directions.East)
	if _tile.Feature(g2td).Type != tile.Grass {
		t.Error("Err grid to tile dir E")
	}

	t2gd = _tile.Placement.TileToGridDir(directions.East)
	if _tile.Feature(t2gd).Type != tile.River {
		t.Error("Err tile to grid dir N")
	}

	g2td = _tile.Placement.GridToTileDir(directions.South)
	if _tile.Feature(g2td).Type != tile.Grass {
		t.Error("Err grid to  tile dir S")
	}

	t2gd = _tile.Placement.TileToGridDir(directions.South)
	if _tile.Feature(t2gd).Type != tile.River {
		t.Error("Err tile to grid dir N")
	}

	g2td = _tile.Placement.GridToTileDir(directions.West)
	if _tile.Feature(g2td).Type != tile.River {
		t.Error("Err grid to  tile dir W")
	}

	t2gd = _tile.Placement.TileToGridDir(directions.West)
	if _tile.Feature(t2gd).Type != tile.Grass {
		t.Error("Err tile to grid dir N")
	}

}
