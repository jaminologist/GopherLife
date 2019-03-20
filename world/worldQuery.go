package world

//TileQuery used for functions that check the contents of a tile
type TileQuery func(*GopherWorldTile) bool

//CheckMapPointForFood Checks if the Tile contains food and also does not have a gopher ontop of it
func CheckMapPointForFood(tile *GopherWorldTile) bool {
	return tile.Food != nil && tile.Gopher == nil
}

//CheckMapPointForEmptySpace Checks if the Tile contains nothing
func CheckMapPointForEmptySpace(tile *GopherWorldTile) bool {
	return tile.Food == nil && tile.Gopher == nil
}

//CheckMapPointForFemaleGopher Checks if a tile contians a female gopher that is looking for a mate
func CheckMapPointForFemaleGopher(tile *GopherWorldTile) bool {
	return tile.Gopher != nil && tile.Gopher.IsLookingForLove() && Female == tile.Gopher.Gender
}

//CheckMapPointForPartner Checks if the Tile contains a sutible partner for the querying gopher
func (gopher *Gopher) CheckMapPointForPartner(tile *GopherWorldTile) bool {
	return tile.Gopher != nil && tile.Gopher.IsLookingForLove() && gopher.Gender.Opposite() == tile.Gopher.Gender
}
