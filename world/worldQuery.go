package world

type MapPointCheck func(*MapPoint) bool

//CheckMapPointForFood Checks if the Tile contains food and also does not have a gopher ontop of it
func CheckMapPointForFood(mapPoint *MapPoint) bool {
	return mapPoint.Food != nil && mapPoint.Gopher == nil
}

//CheckMapPointForEmptySpace Checks if the Tile contains nothing
func CheckMapPointForEmptySpace(mapPoint *MapPoint) bool {
	return mapPoint.Food == nil && mapPoint.Gopher == nil
}

//CheckMapPointForPartner Checks if the Tile contains a sutible partner for the querying gopher
func (gopher *Gopher) CheckMapPointForPartner(mapPoint *MapPoint) bool {
	return mapPoint.Gopher != nil && mapPoint.Gopher.IsLookingForLove() && gopher.Gender.Opposite() == mapPoint.Gopher.Gender
}
