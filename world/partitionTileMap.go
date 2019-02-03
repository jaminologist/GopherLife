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

	gridWidth  int
	gridHeight int

	numberOfGridsWide   int
	numberOfGridsHeight int

	ActiveGophers chan *Gopher

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

	tileMap.numberOfGridsWide = statistics.Width / gridWidth

	if tileMap.numberOfGridsWide*gridWidth < statistics.Width {
		tileMap.numberOfGridsWide++
	}

	tileMap.numberOfGridsHeight = statistics.Height / gridHeight

	if tileMap.numberOfGridsWide*gridHeight < statistics.Height {
		tileMap.numberOfGridsHeight++
	}

	tileMap.gridWidth = gridWidth
	tileMap.gridHeight = gridHeight

	tileMap.BasicGridContainer = NewBasicGridContainer(statistics.Width,
		statistics.Height,
		gridWidth,
		gridHeight,
	)

	tileMap.ActiveGophers = make(chan *Gopher, statistics.NumberOfGophers)
	tileMap.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

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
		}

		tileMap.gopherArray[i] = &gopher
		tileMap.ActiveGophers <- &gopher
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

	numGophers := len(tileMap.ActiveGophers)
	tileMap.gopherArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveGophers
		tileMap.gopherArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.performEntityAction(gopher, secondChannel)

	}
	tileMap.ActiveGophers = secondChannel
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *PartitionTileMap) performEntityAction(gopher *Gopher, channel chan *Gopher) {

	gopher.PerformMoment(tileMap)

	if !gopher.IsDecayed() {

		wait := true

		for wait {
			select {
			case channel <- gopher:
				wait = false
			default:
				//	fmt.Println("Can't Write")
			}
		}

	} else {
		tileMap.QueueRemoveGopher(gopher)
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
			return tile, ok
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

	currentPosition := gopher.Position
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if tileMap.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		tileMap.RemoveGopher(currentPosition.GetX(), currentPosition.GetY(), gopher)
		gopher.Position.Set(targetPosition.GetX(), targetPosition.GetY())
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

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (tileMap *PartitionTileMap) QueuePickUpFood(gopher *Gopher) {
	tileMap.Add(func() {
		food, ok := tileMap.removeFoodFromWorld(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			tileMap.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (tileMap *PartitionTileMap) QueueRemoveGopher(gopher *Gopher) {

	tileMap.Add(func() {
		tileMap.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY(), gopher)
	})
}

func (tileMap *PartitionTileMap) removeFoodFromWorld(x int, y int) (*Food, bool) {

	if tile, ok := tileMap.Tile(x, y); ok {
		food := tile.Food
		tileMap.RemoveFood(x, y, food)
		return food, true
	}

	return nil, false
}

func (tileMap *PartitionTileMap) onFoodPickUp(location calc.Coordinates) {

	size := 50

	xrange := rand.Perm(size)
	yrange := rand.Perm(size)

loop:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {

			newX := location.GetX() + xrange[i] - size/2
			newY := location.GetY() + yrange[j] - size/2

			food := NewPotato()
			if tileMap.InsertFood(newX, newY, &food) {
				break loop
			}

		}
	}

}

func (tileMap *PartitionTileMap) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	tileMap.Add(func() {

		if mapPoint, ok := tileMap.Tile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(tileMap.Statistics.GopherBirthRate)

			emptySpaces := tileMap.Search(gopher.Position, 10, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])

						if len(tileMap.gopherArray) <= tileMap.Statistics.MaximumNumberOfGophers {
							if tileMap.InsertGopher(pos.GetX(), pos.GetY(), &newborn) {
								tileMap.ActiveGophers <- &newborn
							}
						}
					}

				}

			}

		}
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
		return locations
	case SearchForEmptySpace:
		query = CheckMapPointForEmptySpace
	case SearchForFemaleGopher:
		locations := queryForFemalePartner(tileMap, width, height, x, y)
		calc.SortByNearestFromCoordinate(position, locations)
		return locations
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

	startX, startY := tileMap.convertToGridCoordinates(x-width, y-height)
	endX, endY := tileMap.convertToGridCoordinates(x+width, y+height)

	foodLocations := make([]calc.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {

			if grid, ok := tileMap.Grid(x, y); ok {
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

			if grid, ok := tileMap.Grid(x, y); ok {
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
