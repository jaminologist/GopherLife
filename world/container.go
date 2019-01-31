package world

type TileContainer interface {
	Tile(x int, y int) (*Tile, bool)
}

type Basic2DContainer struct {
	grid   [][]*Tile
	width  int
	height int
}

func NewBasic2DContainer(width int, height int) Basic2DContainer {

	container := Basic2DContainer{width: width, height: height}

	container.grid = make([][]*Tile, width)

	for i := 0; i < width; i++ {
		container.grid[i] = make([]*Tile, height)

		for j := 0; j < height; j++ {
			tile := Tile{nil, nil}
			container.grid[i][j] = &tile
		}
	}

	return container
}

func (container *Basic2DContainer) Tile(x int, y int) (*Tile, bool) {

	if x < 0 || x >= container.width || y < 0 || y >= container.height {
		return nil, false
	}

	return container.grid[x][y], true
}
