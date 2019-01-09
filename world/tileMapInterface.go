package world

//TileMapInterface is cool
type TileMapInterface interface {
	Update() bool

	TogglePause()

	Tile(x int, y int) (*Tile, bool)
	SelectedTile() (*Tile, bool)

	SelectEntity(x int, y int) (*Gopher, bool)
	UnSelectGopher()
	SelectRandomGopher()

	Stats() *Statistics
	Diagnostics() *Diagnostics
}
