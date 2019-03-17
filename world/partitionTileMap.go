package world

import (
	"gopherlife/geometry"
)

//PartitionTileMap is cool

const (
	gridWidth  = 5
	gridHeight = 5
)

func CreatePartitionTileMapCustom(settings GopherMapSettings) *GopherWorld {

	gc := NewBasicGridContainer(settings.Width,
		settings.Height,
		gridWidth,
		gridHeight,
	)

	search := GridTileSearch{
		BasicGridContainer: &gc,
	}

	tileMap := NewGopherWorld(&settings, &search, &gc, &gc, &gc, &gc, &gc)

	tileMap.setUpTiles()
	return &tileMap
}

type GridTileSearch struct {
	*BasicGridContainer
}

func (searcher *GridTileSearch) Search(position geometry.Coordinates, width int, height int, maximumFind int, searchType SearchType) []geometry.Coordinates {

	x, y := position.GetX(), position.GetY()

	var locations []geometry.Coordinates

	switch searchType {
	case SearchForFood:
		locations = queryForFood(searcher, width, height, x, y)
	case SearchForFemaleGopher:
		locations = queryForFemalePartner(searcher, width, height, x, y)
	case SearchForEmptySpace:
		sts := SpiralTileSearch{TileContainer: searcher.BasicGridContainer}
		return sts.Search(position, width, height, maximumFind, searchType)
	}

	geometry.SortByNearestFromCoordinate(position, locations)

	if len(locations) >= maximumFind {
		return locations[:maximumFind]
	}
	return locations[:len(locations)]
}

func queryForFood(tileMap *GridTileSearch, width int, height int, x int, y int) []geometry.Coordinates {
	return gridQuery(tileMap, width, height, x, y,

		func(container *TrackedTileContainer) map[int]*Tile {
			return container.foodTileLocations
		},

		func(tile *Tile) (int, int) {
			return tile.Food.Position.GetX(), tile.Food.Position.GetY()
		},

		func(tile *Tile) bool {
			return true
		},
	)
}

func queryForFemalePartner(tileMap *GridTileSearch, width int, height int, x int, y int) []geometry.Coordinates {

	return gridQuery(tileMap, width, height, x, y,

		func(container *TrackedTileContainer) map[int]*Tile {
			return container.gopherTileLocations
		},

		func(tile *Tile) (int, int) {
			return tile.Gopher.Position.GetX(), tile.Gopher.Position.GetY()
		},

		func(tile *Tile) bool {
			return tile.Gopher.Gender == Female && tile.Gopher.IsLookingForLove()
		},
	)
}

func gridQuery(tileMap *GridTileSearch, width int, height int, x int, y int,
	gridSearchFunc func(*TrackedTileContainer) map[int]*Tile,
	coordsFromTile func(*Tile) (int, int),
	tileCheck func(*Tile) bool) []geometry.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-width, y-height, x+width, y+height

	startX, startY := tileMap.convertToGridCoordinates(x-width, y-height)
	endX, endY := tileMap.convertToGridCoordinates(x+width, y+height)

	locations := make([]geometry.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {

			if grid, ok := tileMap.Grid(x*tileMap.BasicGridContainer.gridWidth, y*tileMap.BasicGridContainer.gridHeight); ok {
				potentialLocations := gridSearchFunc(grid)

				for key := range potentialLocations {

					tile := potentialLocations[key]

					i, j := coordsFromTile(tile)

					if i >= worldStartX &&
						i < worldEndX &&
						j >= worldStartY &&
						j < worldEndY && tileCheck(tile) {
						locations = append(locations, geometry.Coordinates{i, j})
					}

				}
			}

		}
	}

	return locations
}
