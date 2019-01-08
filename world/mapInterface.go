package world

type TileMapInterface interface {
	MapPoint(x int, y int) (*Tile, bool)
}
