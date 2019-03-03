package world

import (
	"gopherlife/calc"
	"sync"
)

//PartitionTileMap is cool

const (
	gridWidth  = 5
	gridHeight = 5
)

func CreatePartitionTileMapCustom(statistics Statistics) *GopherMap {

	tileMap := GopherMap{}
	tileMap.Statistics = statistics

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)
	tileMap.QueueableActions = &qa

	gridContainer := NewBasicGridContainer(statistics.Width,
		statistics.Height,
		gridWidth,
		gridHeight,
	)

	tileMap.TileContainer = &gridContainer
	tileMap.InsertableGophers = &gridContainer
	tileMap.InsertableFood = &gridContainer
	tileMap.Searchable = &GridTileSearch{
		BasicGridContainer: &gridContainer,
	}

	tileMap.GopherSliceAndChannel = GopherSliceAndChannel{
		ActiveActors: make(chan *GopherActor, statistics.MaximumNumberOfGophers*2),
		ActiveArray:  make([]*GopherActor, statistics.NumberOfGophers),
	}

	frp := FoodRespawnPickup{InsertableFood: tileMap.InsertableFood}
	tileMap.FoodRespawnPickup = frp

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

	ag := GopherGeneration{
		InsertableGophers:     tileMap.InsertableGophers,
		maxGenerations:        tileMap.Statistics.MaximumNumberOfGophers,
		GopherSliceAndChannel: &tileMap.GopherSliceAndChannel,
	}

	tileMap.GopherGeneration = ag

	tileMap.setUpTiles()
	return &tileMap
}

type GridTileSearch struct {
	*BasicGridContainer
}

func (searcher *GridTileSearch) Search(position calc.Coordinates, width int, height int, maximumFind int, searchType SearchType) []calc.Coordinates {

	x, y := position.GetX(), position.GetY()

	var locations []calc.Coordinates

	switch searchType {
	case SearchForFood:
		locations = queryForFood(searcher, width, height, x, y)
	case SearchForFemaleGopher:
		locations = queryForFemalePartner(searcher, width, height, x, y)
	case SearchForEmptySpace:
		sts := SpiralTileSearch{TileContainer: searcher.BasicGridContainer}
		return sts.Search(position, width, height, maximumFind, searchType)
	}

	calc.SortByNearestFromCoordinate(position, locations)

	if len(locations) >= maximumFind {
		return locations[:maximumFind]
	}
	return locations[:len(locations)]
}

func queryForFood(tileMap *GridTileSearch, width int, height int, x int, y int) []calc.Coordinates {
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

func queryForFemalePartner(tileMap *GridTileSearch, width int, height int, x int, y int) []calc.Coordinates {

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
	tileCheck func(*Tile) bool) []calc.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-width, y-height, x+width, y+height

	startX, startY := tileMap.convertToGridCoordinates(x-width, y-height)
	endX, endY := tileMap.convertToGridCoordinates(x+width, y+height)

	locations := make([]calc.Coordinates, 0)

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
						locations = append(locations, calc.Coordinates{i, j})
					}

				}
			}

		}
	}

	return locations
}
