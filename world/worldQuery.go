package world

import "gopherlife/calc"

type MapPointQuery func(*Tile) bool

//CheckMapPointForFood Checks if the Tile contains food and also does not have a gopher ontop of it
func CheckMapPointForFood(mapPoint *Tile) bool {
	return mapPoint.Food != nil && mapPoint.Gopher == nil
}

//CheckMapPointForEmptySpace Checks if the Tile contains nothing
func CheckMapPointForEmptySpace(mapPoint *Tile) bool {
	return mapPoint.Food == nil && mapPoint.Gopher == nil
}

//CheckMapPointForPartner Checks if the Tile contains a sutible partner for the querying gopher
func (gopher *Gopher) CheckMapPointForPartner(mapPoint *Tile) bool {
	return mapPoint.Gopher != nil && mapPoint.Gopher.IsLookingForLove() && gopher.Gender.Opposite() == mapPoint.Gopher.Gender
}

func Find(tileMap *TileMap, startPosition calc.Coordinates, radius int, maximumFind int, mapPointCheck MapPointQuery) []calc.Coordinates {

	var coordsArray = []calc.Coordinates{}

	spiral := calc.NewSpiral(radius, radius)

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > maximumFind {
			break
		}

		/*if coordinates.X == 0 && coordinates.Y == 0 {
			continue
		}*/

		relativeCoords := startPosition.RelativeCoordinate(coordinates.X, coordinates.Y)

		if mapPoint, ok := tileMap.GetTile(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if mapPointCheck(mapPoint) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(startPosition, coordsArray)

	return coordsArray

}
