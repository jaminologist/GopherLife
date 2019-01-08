package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

type TileMap struct {
	grid [][]*Tile

	actionQueue   chan func()
	ActiveGophers chan *Gopher

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

	Statistics
	Diagnostics
}

type Statistics struct {
	Width                  int
	Height                 int
	NumberOfGophers        int
	MaximumNumberOfGophers int
	GopherBirthRate        int
	NumberOfFood           int
}

type Diagnostics struct {
	globalStopWatch  calc.StopWatch
	inputStopWatch   calc.StopWatch
	gopherStopWatch  calc.StopWatch
	processStopWatch calc.StopWatch
}

func CreateWorldCustom(statistics Statistics) TileMap {

	tileMap := TileMap{}
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

	tileMap.SetUpMapPoints()
	return tileMap

}

func CreateTileMap() TileMap {
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
func (tileMap *TileMap) SelectEntity(x int, y int) (*Gopher, bool) {

	tileMap.SelectedGopher = nil

	if mapPoint, ok := tileMap.GetTile(x, y); ok {
		if mapPoint.Gopher != nil {
			tileMap.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, true
}

//GetTile Gets the given MapPoint at position (x,y)
func (tileMap *TileMap) GetTile(x int, y int) (*Tile, bool) {

	if x < 0 || x >= tileMap.Statistics.Width || y < 0 || y >= tileMap.Statistics.Height {
		return nil, false
	}

	return tileMap.grid[x][y], true
}

//AddFunctionToWorldInputActions is used to store functions that write data to the tileMap.
func (tileMap *TileMap) AddFunctionToWorldInputActions(inputFunction func()) {
	tileMap.actionQueue <- inputFunction
}

//InsertGopher Inserts the given gopher into the tileMap at the specified co-ordinate
func (tileMap *TileMap) InsertGopher(gopher *Gopher, x int, y int) bool {

	if tile, ok := tileMap.GetTile(x, y); ok {
		if !tile.HasGopher() {
			tile.SetGopher(gopher)
			return true
		}
	}

	return false

}

//InsertFood Inserts the given food into the tileMap at the specified co-ordinate
func (tileMap *TileMap) InsertFood(food *Food, x int, y int) bool {

	if tile, ok := tileMap.GetTile(x, y); ok {
		if !tile.HasFood() {
			tile.SetFood(food)
			return true
		}
	}
	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (tileMap *TileMap) RemoveFoodFromWorld(x int, y int) (*Food, bool) {

	if mapPoint, ok := tileMap.GetTile(x, y); ok {
		if mapPoint.HasFood() {
			var food = mapPoint.Food
			mapPoint.ClearFood()
			return food, true
		}
	}

	return nil, false
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (tileMap *TileMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentMapPoint, exists := tileMap.GetTile(gopher.Position.GetX(), gopher.Position.GetY())

	if !exists {
		return false
	}

	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)
	tarGetTile, exists := tileMap.GetTile(targetPosition.GetX(), targetPosition.GetY())

	if exists && tarGetTile.Gopher == nil {

		tarGetTile.Gopher = gopher
		currentMapPoint.Gopher = nil

		gopher.Position = targetPosition

		return true
	}

	return false
}

func (tileMap *TileMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.gopherArray[rand.Intn(len(tileMap.gopherArray))]
}

func (tileMap *TileMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *TileMap) SetUpMapPoints() {

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

func (tileMap *TileMap) onFoodPickUp(location calc.Coordinates) {

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

func (tileMap *TileMap) PerformEntityAction(gopher *Gopher, wg *sync.WaitGroup, channel chan *Gopher) {

	gopher.PerformMoment(tileMap)

	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		tileMap.QueueRemoveGopher(gopher)
	}

	wg.Done()
}

func (tileMap *TileMap) Update() bool {

	if tileMap.IsPaused {
		return false
	}

	if tileMap.SelectedGopher != nil && tileMap.SelectedGopher.IsDecayed() {
		tileMap.SelectRandomGopher()
	}

	if !tileMap.globalStopWatch.IsStarted() {
		tileMap.globalStopWatch.Start()
	}

	tileMap.processStopWatch.Start()
	tileMap.processGophers()
	tileMap.processQueuedTasks()

	if tileMap.Statistics.NumberOfGophers > 0 {
		tileMap.Moments++
	}

	tileMap.processStopWatch.Stop()

	return true

}

func (tileMap *TileMap) processGophers() {

	tileMap.gopherStopWatch.Start()

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

	tileMap.gopherStopWatch.Stop()
}

func (tileMap *TileMap) processQueuedTasks() {

	tileMap.inputStopWatch.Start()

	wait := true
	for wait {
		select {
		case action := <-tileMap.actionQueue:
			action()
		default:
			wait = false
		}
	}

	tileMap.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *TileMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (tileMap *TileMap) QueueRemoveGopher(gopher *Gopher) {

	tileMap.AddFunctionToWorldInputActions(func() {
		//gopher = nil
		if mapPoint, ok := tileMap.GetTile(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			mapPoint.Gopher = nil
		}
	})
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (tileMap *TileMap) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	tileMap.AddFunctionToWorldInputActions(func() {
		success := tileMap.MoveGopher(gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (tileMap *TileMap) QueuePickUpFood(gopher *Gopher) {

	tileMap.AddFunctionToWorldInputActions(func() {
		food, ok := tileMap.RemoveFoodFromWorld(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			tileMap.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})

}

func (tileMap *TileMap) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	tileMap.AddFunctionToWorldInputActions(func() {

		if mapPoint, ok := tileMap.GetTile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(tileMap.Statistics.GopherBirthRate)

			emptySpaces := Find(tileMap, gopher.Position, 10, litterNumber, CheckMapPointForEmptySpace)

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
