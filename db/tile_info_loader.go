package db

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

type TileInfoLoader interface {
	GetTileInfo(string) (TileInfo, error)
}
