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

	ActiveActors chan *GopherActor

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

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

	tileMap.ActiveActors = make(chan *GopherActor, statistics.MaximumNumberOfGophers*2)
	tileMap.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

	frp := FoodRespawnPickup{Insertable: &tileMap.BasicGridContainer}
	tileMap.FoodRespawnPickup = frp

	ag := GopherGeneration{
		Insertable:     &tileMap.BasicGridContainer,
		maxGenerations: tileMap.Statistics.MaximumNumberOfGophers,
		ActiveGophers:  tileMap.ActiveActors,
		gopherArray:    tileMap.gopherArray,
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

		tileMap.gopherArray[i] = &gopher
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
	tileMap.Statistics.NumberOfGophers = len(tileMap.ActiveGophers)
	if tileMap.Statistics.NumberOfGophers > 0 {
		tileMap.Moments++
	}

	tileMap.diagnostics.processStopWatch.Stop()

	return true

}

func (tileMap *PartitionTileMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

	numGophers := len(tileMap.ActiveActors)
	tileMap.gopherArray = make([]*Gopher, numGophers)
	tileMap.GopherGeneration.gopherArray = tileMap.gopherArray

	secondChannel := make(chan *GopherActor, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveActors
		tileMap.gopherArray[i] = gopher.Gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.Act(gopher, secondChannel)

	}
	tileMap.ActiveActors = secondChannel
	tileMap.GopherGeneration.ActiveGophers = tileMap.ActiveActors
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
	tileMap.SelectedGopher = tileMap.gopherArray[rand.Intn(len(tileMap.gopherArray))]
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

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (tileMap *PartitionTileMap) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	tileMap.Add(func() {
		success := tileMap.MoveGopher(gopher, moveX, moveY)
		_ = success
	})
}

func (tileMap *PartitionTileMap) Search(position calc.Coordinates, width int, height int, maximumFind int, searchType SearchType) []calc.Coordinates {
	var coordsArray = []calc.Coordinates{}

	x, y := position.GetX(), position.GetY()

	spiral := calc.NewSpiral(width, height)

	var query TileQuery

	switch searchType {
	case SearchForFood:
		locations := queryForFood(tileMap, width, height, x, y)
		calc.SortByNearestFromCoordinate(position, locations)

		if len(locations) >= maximumFind {
			return locations[:maximumFind]
		} else {
			return locations[:len(locations)]
		}

		query = CheckMapPointForFood

	case SearchForEmptySpace:
		query = CheckMapPointForEmptySpace
	case SearchForFemaleGopher:
		locations := queryForFemalePartner(tileMap, width, height, x, y)
		calc.SortByNearestFromCoordinate(position, locations)

		if len(locations) >= maximumFind {
			return locations[:maximumFind]
		} else {
			return locations[:len(locations)]
		}
		//query = CheckMapPointForFemaleGopher
	}

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > maximumFind {
			break
		}

		relativeCoords := calc.RelativeCoordinate(coordinates, x, y)

		if tile, ok := tileMap.Tile(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if query(tile) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(position, coordsArray)

	return coordsArray
}

func queryForFood(tileMap *PartitionTileMap, width int, height int, x int, y int) []calc.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-width, y-height, x+width, y+height

	startX, startY := tileMap.convertToGridCoordinates(worldStartX, worldStartY)
	endX, endY := tileMap.convertToGridCoordinates(worldEndX, worldEndY)

	foodLocations := make([]calc.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {
			if grid, ok := tileMap.Grid(x*tileMap.BasicGridContainer.gridWidth, y*tileMap.BasicGridContainer.gridHeight); ok {
				for key := range grid.foodTileLocations {

					tile := grid.foodTileLocations[key]

					i, j := tile.Food.Position.GetX(), tile.Food.Position.GetY()

					if i >= worldStartX &&
						i < worldEndX &&
						j >= worldStartY &&
						j < worldEndY {
						foodLocations = append(foodLocations, tile.Food.Position)
					}
				}
			}

		}
	}

	return foodLocations
}

func queryForFemalePartner(tileMap *PartitionTileMap, width int, height int, x int, y int) []calc.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-width, y-height, x+width, y+height

	startX, startY := tileMap.convertToGridCoordinates(x-width, y-height)
	endX, endY := tileMap.convertToGridCoordinates(x+width, y+height)

	locations := make([]calc.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {

			if grid, ok := tileMap.Grid(x*tileMap.BasicGridContainer.gridWidth, y*tileMap.BasicGridContainer.gridHeight); ok {
				for key := range grid.gopherTileLocations {

					tile := grid.gopherTileLocations[key]

					i, j := tile.Gopher.Position.GetX(), tile.Gopher.Position.GetY()
					if i >= worldStartX &&
						i < worldEndX &&
						j >= worldStartY &&
						j < worldEndY && tile.Gopher.Gender == Female && tile.Gopher.IsLookingForLove() {
						locations = append(locations, tile.Gopher.Position)
					}
				}
			}

		}
	}

	return locations
}
