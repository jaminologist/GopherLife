package world

type TileQuery func(*Tile) bool

//CheckMapPointForFood Checks if the Tile contains food and also does not have a gopher ontop of it
func CheckMapPointForFood(tile *Tile) bool {
	return tile.Food != nil && tile.Gopher == nil
}

//CheckMapPointForEmptySpace Checks if the Tile contains nothing
func CheckMapPointForEmptySpace(tile *Tile) bool {
	return tile.Food == nil && tile.Gopher == nil
}

//CheckMapPointForPartner Checks if the Tile contains a sutible partner for the querying gopher
func (gopher *Gopher) CheckMapPointForPartner(tile *Tile) bool {
	return tile.Gopher != nil && tile.Gopher.IsLookingForLove() && gopher.Gender.Opposite() == tile.Gopher.Gender
}
