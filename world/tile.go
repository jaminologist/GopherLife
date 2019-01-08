package world

type Tile struct {
	Gopher *Gopher
	Food   *Food
}

func (tile *Tile) isEmpty() bool {
	return tile.Gopher == nil && tile.Food == nil
}

func (tile *Tile) HasGopher() bool {
	return tile.Gopher != nil
}

func (tile *Tile) HasFood() bool {
	return tile.Food != nil
}

func (tile *Tile) SetGopher(g *Gopher) {
	tile.Gopher = g
}

func (tile *Tile) SetFood(f *Food) {
	tile.Food = f
}

func (tile *Tile) ClearGopher() {
	tile.Gopher = nil
}

func (tile *Tile) ClearFood() {
	tile.Food = nil
}
