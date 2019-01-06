package world

type entity interface {
	PerformAction(tileMap *TileMap)
}
