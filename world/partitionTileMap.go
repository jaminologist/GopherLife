package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

//PartitionTileMap is cool

const (
	gridWidth  = 5
	gridHeight = 5
)

type PartitionTileMap struct {
	QueueableActions
	BasicGridContainer
	FoodRespawnPickup
	GopherGeneration

	GopherSliceAndChannel

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher

	Moments  int
	IsPaused bool

	Statistics
	diagnostics Diagnostics
}

func CreatePartitionTileMapCustom(statistics Statistics) *PartitionTileMap {

	tileMap := PartitionTileMap{}
	tileMap.Statistics = statistics

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)
	tileMap.QueueableActions = &qa

	tileMap.BasicGridContainer = NewBasicGridContainer(statistics.Width,
		statistics.Height,
		gridWidth,
		gridHeight,
	)

	tileMap.GopherSliceAndChannel = GopherSliceAndChannel{
		ActiveActors: make(chan *GopherActor, statistics.MaximumNumberOfGophers*2),
		ActiveArray:  make([]*GopherActor, statistics.NumberOfGophers),
	}

	frp := FoodRespawnPickup{Insertable: &tileMap.BasicGridContainer}
	tileMap.FoodRespawnPickup = frp

	ag := GopherGeneration{
		Insertable:            &tileMap.BasicGridContainer,
		maxGenerations:        tileMap.Statistics.MaximumNumberOfGophers,
		GopherSliceAndChannel: &tileMap.GopherSliceAndChannel,
	}

	tileMap.GopherGeneration = ag

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

	tileMap.setUpTiles()
	return &tileMap
}

func CreatePartitionTileMap() *PartitionTileMap {
	tileMap := CreatePartitionTileMapCustom(
		Statistics{
			Width:                  3000,
			Height:                 3000,
			NumberOfGophers:        5000,
			NumberOfFood:           50000,
			MaximumNumberOfGophers: 100000,
			GopherBirthRate:        7,
		},
	)
	return tileMap
}

func (tileMap *PartitionTileMap) setUpTiles() {

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		tileMap.Statistics.Width, tileMap.Statistics.Height)

	count := 0

	for i := 0; i < tileMap.Statistics.NumberOfGophers; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.CuteName(), pos)

		tileMap.BasicGridContainer.InsertGopher(pos.GetX(), pos.GetY(), &gopher)

		if i == 0 {
			tileMap.SelectedGopher = &gopher
			tileMap.SelectedGopher.Position = pos
		}

		var gopherActor = GopherActor{
			Gopher:           &gopher,
			GopherBirthRate:  tileMap.Statistics.GopherBirthRate,
			QueueableActions: tileMap.QueueableActions,
			Searchable:       tileMap,
			TileContainer:    &tileMap.BasicGridContainer,
			Insertable:       &tileMap.BasicGridContainer,
			PickableTiles:    tileMap,
			MoveableActors:   tileMap,
			ActorGeneration:  &tileMap.GopherGeneration,
		}

		tileMap.ActiveArray[i] = &gopherActor
		tileMap.ActiveActors <- &gopherActor
		count++
	}

	for i := 0; i < tileMap.Statistics.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		tileMap.BasicGridContainer.InsertFood(pos.GetX(), pos.GetY(), &food)
		count++
	}
}

//TogglePause Toggles the pause
func (tileMap *PartitionTileMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

func (tileMap *PartitionTileMap) Update() bool {

	if tileMap.IsPaused {
		return false
	}

	if tileMap.SelectedGopher != nil && tileMap.SelectedGopher.IsDecayed() {
		tileMap.SelectRandomGopher()
	}

	if !tileMap.diagnostics.globalStopWatch.IsStarted() {
		tileMap.diagnostics.globalStopWatch.Start()
	}

	tileMap.diagnostics.processStopWatch.Start()
	tileMap.processGophers()
	tileMap.processQueuedTasks()
	tileMap.Statistics.NumberOfGophers = len(tileMap.ActiveActors)
	if tileMap.Statistics.NumberOfGophers > 0 {
		tileMap.Moments++
	}

	tileMap.diagnostics.processStopWatch.Stop()

	return true

}

func (tileMap *PartitionTileMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

	numGophers := len(tileMap.ActiveActors)
	tileMap.GopherSliceAndChannel.ActiveArray = make([]*GopherActor, numGophers)

	secondChannel := make(chan *GopherActor, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveActors
		tileMap.ActiveArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.Act(gopher, secondChannel)

	}
	tileMap.ActiveActors = secondChannel
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *PartitionTileMap) Act(gopher *GopherActor, channel chan *GopherActor) {
	gopher.Update()
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		gopher.QueueRemoveGopher()
	}
	tileMap.GopherWaitGroup.Done()
}

func (tileMap *PartitionTileMap) processQueuedTasks() {
	tileMap.diagnostics.inputStopWatch.Start()
	tileMap.QueueableActions.Process()
	tileMap.diagnostics.inputStopWatch.Stop()
}

func (tileMap *PartitionTileMap) SelectedTile() (*Tile, bool) {

	if tileMap.SelectedGopher != nil {
		if tile, ok := tileMap.Tile(tileMap.SelectedGopher.Position.GetX(), tileMap.SelectedGopher.Position.GetY()); ok {
			if tile.HasGopher() {
				return tile, ok
			}
		}
	}
	return nil, false

}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the tileMap
//If there is not a gopher at the give coordinates this function returns zero.
func (tileMap *PartitionTileMap) SelectEntity(x int, y int) (*Gopher, bool) {

	tileMap.SelectedGopher = nil

	if mapPoint, ok := tileMap.Tile(x, y); ok {
		if mapPoint.Gopher != nil {
			tileMap.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, true
}

func (tileMap *PartitionTileMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.ActiveArray[rand.Intn(len(tileMap.ActiveArray))].Gopher
}

func (tileMap *PartitionTileMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *PartitionTileMap) Stats() *Statistics {
	return &tileMap.Statistics
}

func (tileMap *PartitionTileMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

func (tileMap *PartitionTileMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := calc.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if tileMap.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		tileMap.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false
}

func (tileMap *PartitionTileMap) Search(position calc.Coordinates, width int, height int, maximumFind int, searchType SearchType) []calc.Coordinates {

	x, y := position.GetX(), position.GetY()

	var locations []calc.Coordinates

	switch searchType {
	case SearchForFood:
		locations = queryForFood(tileMap, width, height, x, y)
	case SearchForFemaleGopher:
		locations = queryForFemalePartner(tileMap, width, height, x, y)
	case SearchForEmptySpace:
		sts := SpiralTileSearch{TileContainer: &tileMap.BasicGridContainer}
		return sts.Search(position, width, height, maximumFind, searchType)
	}

	calc.SortByNearestFromCoordinate(position, locations)

	if len(locations) >= maximumFind {
		return locations[:maximumFind]
	}
	return locations[:len(locations)]
}

func queryForFood(tileMap *PartitionTileMap, width int, height int, x int, y int) []calc.Coordinates {
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

func queryForFemalePartner(tileMap *PartitionTileMap, width int, height int, x int, y int) []calc.Coordinates {

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

func gridQuery(tileMap *PartitionTileMap, width int, height int, x int, y int,
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
