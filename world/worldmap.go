package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

type World struct {
	grid [][]*MapPoint

	actionQueue   chan func()
	ActiveGophers chan *Gopher

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

	Statistics Statistics

	globalStopWatch  calc.StopWatch
	inputStopWatch   calc.StopWatch
	gopherStopWatch  calc.StopWatch
	processStopWatch calc.StopWatch
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

func CreateWorldCustom(statistics Statistics) World {

	world := World{}
	world.Statistics = statistics
	world.actionQueue = make(chan func(), statistics.MaximumNumberOfGophers*2)

	world.grid = make([][]*MapPoint, statistics.Width)

	for i := 0; i < statistics.Width; i++ {
		world.grid[i] = make([]*MapPoint, statistics.Height)

		for j := 0; j < statistics.Height; j++ {
			mp := MapPoint{nil, nil}
			world.grid[i][j] = &mp
		}
	}

	world.ActiveGophers = make(chan *Gopher, statistics.NumberOfGophers)
	world.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

	var wg sync.WaitGroup
	world.GopherWaitGroup = &wg

	world.SetUpMapPoints()
	return world

}

func CreateWorld() World {
	world := CreateWorldCustom(
		Statistics{
			Width:                  3000,
			Height:                 3000,
			NumberOfGophers:        5000,
			NumberOfFood:           1000000,
			MaximumNumberOfGophers: 100000,
			GopherBirthRate:        7,
		},
	)
	return world
}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the world
//If there is not a gopher at the give coordinates this function returns zero.
func (world *World) SelectEntity(x int, y int) (*Gopher, bool) {

	world.SelectedGopher = nil

	if mapPoint, ok := world.GetMapPoint(x, y); ok {
		if mapPoint.Gopher != nil {
			world.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, true
}

//GetMapPoint Gets the given MapPoint at position (x,y)
func (world *World) GetMapPoint(x int, y int) (*MapPoint, bool) {

	if x < 0 || x >= world.Statistics.Width || y < 0 || y >= world.Statistics.Height {
		return nil, false
	}

	return world.grid[x][y], true
}

//AddFunctionToWorldInputActions is used to store functions that write data to the world.
func (world *World) AddFunctionToWorldInputActions(inputFunction func()) {
	world.actionQueue <- inputFunction
}

//InsertGopher Inserts the given gopher into the world at the specified co-ordinate
func (world *World) InsertGopher(gopher *Gopher, x int, y int) bool {

	if mp, ok := world.GetMapPoint(x, y); ok {
		if !mp.HasGopher() {
			mp.SetGopher(gopher)
			return true
		}
	}

	return false

}

//InsertFood Inserts the given food into the world at the specified co-ordinate
func (world *World) InsertFood(food *Food, x int, y int) bool {

	if mp, ok := world.GetMapPoint(x, y); ok {
		if !mp.HasFood() {
			mp.SetFood(food)
			return true
		}
	}
	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (world *World) RemoveFoodFromWorld(x int, y int) (*Food, bool) {

	if mapPoint, ok := world.GetMapPoint(x, y); ok {
		if mapPoint.HasFood() {
			var food = mapPoint.Food
			mapPoint.ClearFood()
			return food, true
		}
	}

	return nil, false
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (world *World) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentMapPoint, exists := world.GetMapPoint(gopher.Position.GetX(), gopher.Position.GetY())

	if !exists {
		return false
	}

	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)
	targetMapPoint, exists := world.GetMapPoint(targetPosition.GetX(), targetPosition.GetY())

	if exists && targetMapPoint.Gopher == nil {

		targetMapPoint.Gopher = gopher
		currentMapPoint.Gopher = nil

		gopher.Position = targetPosition

		return true
	}

	return false
}

func (world *World) SelectRandomGopher() {
	world.SelectedGopher = world.gopherArray[rand.Intn(len(world.gopherArray))]
}

func (world *World) UnSelectGopher() {
	world.SelectedGopher = nil
}

func (world *World) SetUpMapPoints() {

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		world.Statistics.Width, world.Statistics.Height)

	count := 0

	for i := 0; i < world.Statistics.NumberOfGophers; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.GetCuteName(), pos)

		world.InsertGopher(&gopher, pos.GetX(), pos.GetY())

		if i == 0 {
			world.SelectedGopher = &gopher
		}

		world.gopherArray[i] = &gopher
		world.ActiveGophers <- &gopher
		count++
	}

	for i := 0; i < world.Statistics.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		world.InsertFood(&food, pos.GetX(), pos.GetY())
		count++
	}

}

func (world *World) onFoodPickUp(location calc.Coordinates) {

	size := 50
	xrange, yrange := rand.Perm(size), rand.Perm(size)
	food := NewPotato()

loop:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newX, newY := location.GetX()+xrange[i]-size/2, location.GetY()+yrange[j]-size/2
			if world.InsertFood(&food, newX, newY) {
				break loop
			}
		}
	}
}

func (world *World) PerformEntityAction(gopher *Gopher, wg *sync.WaitGroup, channel chan *Gopher) {

	gopher.PerformMoment(world)

	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		world.QueueRemoveGopher(gopher)
	}

	wg.Done()
}

func (world *World) ProcessWorld() bool {

	if world.IsPaused {
		return false
	}

	if world.SelectedGopher != nil && world.SelectedGopher.IsDecayed() {
		world.SelectRandomGopher()
	}

	if !world.globalStopWatch.IsStarted() {
		world.globalStopWatch.Start()
	}

	world.processStopWatch.Start()
	world.processGophers()
	world.processQueuedTasks()

	if world.Statistics.NumberOfGophers > 0 {
		world.Moments++
	}

	world.processStopWatch.Stop()

	return true

}

func (world *World) processGophers() {

	world.gopherStopWatch.Start()

	numGophers := len(world.ActiveGophers)
	world.gopherArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-world.ActiveGophers
		world.gopherArray[i] = gopher
		world.GopherWaitGroup.Add(1)
		go world.PerformEntityAction(gopher, world.GopherWaitGroup, secondChannel)

	}
	world.ActiveGophers = secondChannel
	world.GopherWaitGroup.Wait()

	world.gopherStopWatch.Stop()
}

func (world *World) processQueuedTasks() {

	world.inputStopWatch.Start()

	wait := true
	for wait {
		select {
		case action := <-world.actionQueue:
			action()
		default:
			wait = false
		}
	}

	world.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (world *World) TogglePause() {
	world.IsPaused = !world.IsPaused
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (world *World) QueueRemoveGopher(gopher *Gopher) {

	world.AddFunctionToWorldInputActions(func() {
		//gopher = nil
		if mapPoint, ok := world.GetMapPoint(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			mapPoint.Gopher = nil
		}
	})
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (world *World) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	world.AddFunctionToWorldInputActions(func() {
		success := world.MoveGopher(gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (world *World) QueuePickUpFood(gopher *Gopher) {

	world.AddFunctionToWorldInputActions(func() {
		food, ok := world.RemoveFoodFromWorld(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			world.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})

}

func (world *World) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	world.AddFunctionToWorldInputActions(func() {

		if mapPoint, ok := world.GetMapPoint(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(world.Statistics.GopherBirthRate)

			emptySpaces := Find(world, gopher.Position, 10, litterNumber, CheckMapPointForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						pos := emptySpaces[i]
						newborn := NewGopher(names.GetCuteName(), emptySpaces[i])

						if len(world.gopherArray) <= world.Statistics.MaximumNumberOfGophers {
							if world.InsertGopher(&newborn, pos.GetX(), pos.GetY()) {
								world.ActiveGophers <- &newborn
							}
						}
					}

				}

			}

		}
	})

}
