package world

type Tile struct {
	Gopher *Gopher
	Food   *Food
}

//NewTile Returns a new tile which hold the given gopher and food
func NewTile(gopher *Gopher, food *Food) Tile {
	return Tile{Gopher: gopher, Food: food}
}

func (tile *Tile) IsEmpty() bool {
	return tile.Gopher == nil && tile.Food == nil
}

//HasGopher Checks if this tile contains a gopher
func (tile *Tile) HasGopher() bool {
	return tile.Gopher != nil
}

//HasFood Checks if this tile contains food
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
