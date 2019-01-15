package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

type SpiralSearchTileMap struct {
	grid [][]*Tile

	actionQueue   chan func()
	ActiveGophers chan *Gopher

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

	Statistics
	diagnostics Diagnostics
}

func CreateWorldCustom(statistics Statistics) SpiralSearchTileMap {

	tileMap := SpiralSearchTileMap{}
	tileMap.Statistics = statistics
	tileMap.actionQueue = make(chan func(), statistics.MaximumNumberOfGophers*2)

	tileMap.grid = make([][]*Tile, statistics.Width)

	for i := 0; i < statistics.Width; i++ {
		tileMap.grid[i] = make([]*Tile, statistics.Height)

		for j := 0; j < statistics.Height; j++ {
			tile := Tile{nil, nil}
			tileMap.grid[i][j] = &tile
		}
	}

	tileMap.ActiveGophers = make(chan *Gopher, statistics.NumberOfGophers)
	tileMap.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

	tileMap.setUpMapPoints()
	return tileMap

}

func CreateTileMap() SpiralSearchTileMap {
	tileMap := CreateWorldCustom(
		Statistics{
			Width:                  3000,
			Height:                 3000,
			NumberOfGophers:        5000,
			NumberOfFood:           1000000,
			MaximumNumberOfGophers: 100000,
			GopherBirthRate:        7,
		},
	)
	return tileMap
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

func (tileMap *SpiralSearchTileMap) SelectedTile() (*Tile, bool) {

	if tileMap.SelectedGopher != nil {
		if tile, ok := tileMap.Tile(tileMap.SelectedGopher.Position.GetX(), tileMap.SelectedGopher.Position.GetY()); ok {
			return tile, ok
		}
	}
	return nil, false

}

//GetTile Gets the given MapPoint at position (x,y)
func (tileMap *SpiralSearchTileMap) Tile(x int, y int) (*Tile, bool) {

	if x < 0 || x >= tileMap.Statistics.Width || y < 0 || y >= tileMap.Statistics.Height {
		return nil, false
	}

	return tileMap.grid[x][y], true
}

//AddFunctionToWorldInputActions is used to store functions that write data to the tileMap.
func (tileMap *SpiralSearchTileMap) AddFunctionToWorldInputActions(inputFunction func()) {
	tileMap.actionQueue <- inputFunction
}

//InsertGopher Inserts the given gopher into the tileMap at the specified co-ordinate
func (tileMap *SpiralSearchTileMap) InsertGopher(gopher *Gopher, x int, y int) bool {

	if tile, ok := tileMap.Tile(x, y); ok {
		if !tile.HasGopher() {
			tile.SetGopher(gopher)
			return true
		}
	}

	return false

}

//InsertFood Inserts the given food into the tileMap at the specified co-ordinate
func (tileMap *SpiralSearchTileMap) InsertFood(food *Food, x int, y int) bool {

	if tile, ok := tileMap.Tile(x, y); ok {
		if !tile.HasFood() {
			tile.SetFood(food)
			return true
		}
	}
	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (tileMap *SpiralSearchTileMap) RemoveFoodFromWorld(x int, y int) (*Food, bool) {

	if mapPoint, ok := tileMap.Tile(x, y); ok {
		if mapPoint.HasFood() {
			var food = mapPoint.Food
			mapPoint.ClearFood()
			return food, true
		}
	}

	return nil, false
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (tileMap *SpiralSearchTileMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentMapPoint, exists := tileMap.Tile(gopher.Position.GetX(), gopher.Position.GetY())

	if !exists {
		return false
	}

	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)
	tarGetTile, exists := tileMap.Tile(targetPosition.GetX(), targetPosition.GetY())

	if exists && tarGetTile.Gopher == nil {

		tarGetTile.Gopher = gopher
		currentMapPoint.Gopher = nil

		gopher.Position = targetPosition

		return true
	}

	return false
}

func (tileMap *SpiralSearchTileMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.gopherArray[rand.Intn(len(tileMap.gopherArray))]
}

func (tileMap *SpiralSearchTileMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *SpiralSearchTileMap) Stats() *Statistics {
	return &tileMap.Statistics
}

func (tileMap *SpiralSearchTileMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

func (tileMap *SpiralSearchTileMap) setUpMapPoints() {

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		tileMap.Statistics.Width, tileMap.Statistics.Height)

	count := 0

	for i := 0; i < tileMap.Statistics.NumberOfGophers; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.CuteName(), pos)

		tileMap.InsertGopher(&gopher, pos.GetX(), pos.GetY())

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
		tileMap.InsertFood(&food, pos.GetX(), pos.GetY())
		count++
	}

}

func (tileMap *SpiralSearchTileMap) onFoodPickUp(location calc.Coordinates) {

	size := 50
	xrange, yrange := rand.Perm(size), rand.Perm(size)
	food := NewPotato()

loop:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newX, newY := location.GetX()+xrange[i]-size/2, location.GetY()+yrange[j]-size/2
			if tileMap.InsertFood(&food, newX, newY) {
				break loop
			}
		}
	}
}

func (tileMap *SpiralSearchTileMap) PerformEntityAction(gopher *Gopher, wg *sync.WaitGroup, channel chan *Gopher) {

	gopher.PerformMoment(tileMap)

	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		tileMap.QueueRemoveGopher(gopher)
	}

	wg.Done()
}

func (tileMap *SpiralSearchTileMap) Update() bool {

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

func (tileMap *SpiralSearchTileMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

	numGophers := len(tileMap.ActiveGophers)
	tileMap.gopherArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveGophers
		tileMap.gopherArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.PerformEntityAction(gopher, tileMap.GopherWaitGroup, secondChannel)

	}
	tileMap.ActiveGophers = secondChannel
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *SpiralSearchTileMap) processQueuedTasks() {

	tileMap.diagnostics.inputStopWatch.Start()

	wait := true
	for wait {
		select {
		case action := <-tileMap.actionQueue:
			action()
		default:
			wait = false
		}
	}

	tileMap.diagnostics.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *SpiralSearchTileMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (tileMap *SpiralSearchTileMap) QueueRemoveGopher(gopher *Gopher) {

	tileMap.AddFunctionToWorldInputActions(func() {
		//gopher = nil
		if mapPoint, ok := tileMap.Tile(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			mapPoint.Gopher = nil
		}
	})
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (tileMap *SpiralSearchTileMap) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	tileMap.AddFunctionToWorldInputActions(func() {
		success := tileMap.MoveGopher(gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (tileMap *SpiralSearchTileMap) QueuePickUpFood(gopher *Gopher) {

	tileMap.AddFunctionToWorldInputActions(func() {
		food, ok := tileMap.RemoveFoodFromWorld(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			tileMap.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})
}

func (tileMap *SpiralSearchTileMap) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	tileMap.AddFunctionToWorldInputActions(func() {

		if mapPoint, ok := tileMap.Tile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(tileMap.Statistics.GopherBirthRate)

			emptySpaces := tileMap.Search(gopher.Position, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])

						if len(tileMap.gopherArray) <= tileMap.Statistics.MaximumNumberOfGophers {
							if tileMap.InsertGopher(&newborn, pos.GetX(), pos.GetY()) {
								tileMap.ActiveGophers <- &newborn
							}
						}
					}

				}

			}

		}
	})

}

func (tileMap *SpiralSearchTileMap) Search(startPosition calc.Coordinates, radius int, maximumFind int, searchType SearchType) []calc.Coordinates {
	var coordsArray = []calc.Coordinates{}

	spiral := calc.NewSpiral(radius, radius)

	var query TileQuery

	switch searchType {
	case SearchForFood:
		query = CheckMapPointForFood
	case SearchForEmptySpace:
		query = CheckMapPointForEmptySpace
	case SearchForFemaleGopher:
		query = CheckMapPointForFemaleGopher
	}

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > maximumFind {
			break
		}

		/*if coordinates.X == 0 && coordinates.Y == 0 {
			continue
		}*/

		relativeCoords := startPosition.RelativeCoordinate(coordinates.X, coordinates.Y)

		if tile, ok := tileMap.Tile(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if query(tile) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(startPosition, coordsArray)

	return coordsArray
}
