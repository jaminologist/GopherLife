package world

import (
	"gopherlife/calc"
	"gopherlife/food"
	"gopherlife/names"
	"math/rand"
	"sync"
)

const numberOfGophs = 5000
const numberOfFoods = 1000000
const worldSize = 3000

type World struct {
	grid [][]*MapPoint

	width  int
	height int

	InputActions chan func()

	GopherWaitGroup sync.WaitGroup

	ActiveGophers chan *Gopher

	SelectedGopher *Gopher

	gopherArray []*Gopher

	Moments int

	IsPaused bool

	avgProcessingTime int

	globalStopWatch  calc.StopWatch
	inputStopWatch   calc.StopWatch
	gopherStopWatch  calc.StopWatch
	processStopWatch calc.StopWatch
}

func CreateWorldCustom(width int, height int, gophers int, food int) World {

	world := World{width: width, height: height}
	world.InputActions = make(chan func(), 1000000)

	world.grid = make([][]*MapPoint, worldSize)

	for i := 0; i < width; i++ {
		world.grid[i] = make([]*MapPoint, height)

		for j := 0; j < height; j++ {
			mp := MapPoint{}
			mp.Gopher = nil
			mp.Food = nil
			world.grid[i][j] = &mp
		}
	}

	world.ActiveGophers = make(chan *Gopher, gophers)
	world.gopherArray = make([]*Gopher, gophers)

	world.SetUpMapPoints(gophers, food)
	return world

}

func CreateWorld() World {
	world := CreateWorldCustom(worldSize, worldSize, numberOfGophs, numberOfFoods)
	return world
}

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

func (world *World) GetMapPoint(x int, y int) (*MapPoint, bool) {

	if x < 0 || x >= worldSize || y < 0 || y >= worldSize {
		return nil, false
	}

	return world.grid[x][y], true
}

func (world *World) AddFunctionToWorldInputActions(inputFunction func()) {

loop:
	for {
		select {
		case world.InputActions <- inputFunction:
			break loop
		default:
			//fmt.Println("Stuck here forever ", rand.Intn(100))
		}

	}

}

func (world *World) RemoveFoodFromWorld(position calc.Coordinates) (*food.Food, bool) {

	if mapPoint, ok := world.GetMapPoint(position.GetX(), position.GetY()); ok {
		if mapPoint.Food != nil {

			var food = mapPoint.Food
			mapPoint.Food = nil
			return food, true
		}
	}

	return nil, false
}

func (world *World) InsertGopher(gopher *Gopher, x int, y int) bool {

	if mp, ok := world.GetMapPoint(x, y); ok {
		if !mp.HasGopher() {
			mp.SetGopher(gopher)
		}
	}

	return false

}

func (world *World) MoveGopher(gopher *Gopher, x int, y int) bool {

	currentMapPoint, exists := world.GetMapPoint(gopher.Position.GetX(), gopher.Position.GetY())

	if !exists {
		return false
	}

	targetPosition := gopher.Position.RelativeCoordinate(x, y)
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

func (world *World) SetUpMapPoints(numberOfGophers int, numberOfFood int) {

	keys := make([]calc.Coordinates, worldSize*worldSize)

	keycount := 0

	for i := 0; i < worldSize; i++ {
		for j := 0; j < worldSize; j++ {
			keys[keycount] = calc.NewCoordinate(i, j)
			keycount++
		}
	}

	rand.Seed(1)

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	count := 0

	for i := 0; i < numberOfGophers; i++ {

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

	for i := 0; i < numberOfFood; i++ {
		pos := keys[count]
		mapPoint, _ := world.GetMapPoint(pos.GetX(), pos.GetY())
		var food = food.NewPotato()

		mapPoint.Food = &food
		count++
	}

}

func (world *World) onFoodPickUp(location calc.Coordinates) {

	size := 50

	xrange := rand.Perm(size)
	yrange := rand.Perm(size)

	for i := 0; i < size; i++ {

		isDone := false

		for j := 0; j < size; j++ {
			newFoodLocation := calc.NewCoordinate(
				location.GetX()+xrange[i]-size/2,
				location.GetY()+yrange[j]-size/2)

			if mapPoint, ok := world.GetMapPoint(newFoodLocation.GetX(), newFoodLocation.GetY()); ok {

				if mapPoint.Food == nil {
					var food = food.NewPotato()
					mapPoint.Food = &food

					isDone = true
					break
				}

			}

		}

		if isDone {
			break

		}
	}

}

func (world *World) PerformEntityAction(gopher *Gopher, wg *sync.WaitGroup, channel chan *Gopher) {

	gopher.PerformMoment(world)

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
		world.QueueRemoveGopher(gopher)
	}

	wg.Done()

}

func (world *World) ProcessWorld() bool {

	if world.IsPaused {
		return false
	}

	if !world.globalStopWatch.IsStarted() {
		world.globalStopWatch.Start()
	}

	world.processStopWatch.Start()
	world.gopherStopWatch.Start()

	numGophers := len(world.ActiveGophers)

	world.gopherArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-world.ActiveGophers
		world.gopherArray[i] = gopher
		world.GopherWaitGroup.Add(1)
		go world.PerformEntityAction(gopher, &world.GopherWaitGroup, secondChannel)

	}

	world.ActiveGophers = secondChannel

	world.GopherWaitGroup.Wait()
	world.gopherStopWatch.Stop()

	world.inputStopWatch.Start()
	readingInputActionsUsingJustChannel(world)
	world.inputStopWatch.Stop()

	if numGophers > 0 {
		world.Moments++
	}

	world.processStopWatch.Stop()

	return true

}

func readingInputActionsUsingGoFunc(world *World) {

	var mutex = &sync.Mutex{}

	wait := true
	for wait {
		select {
		case action := <-world.InputActions:
			go func() {
				mutex.Lock()
				action()
				mutex.Unlock()
			}()
		default:
			wait = false
		}
	}

}

func readingInputActionsUsingJustChannel(world *World) {

	wait := true
	for wait {
		select {
		case action := <-world.InputActions:
			action()
		default:
			wait = false
		}
	}

}

func (world *World) TogglePause() {
	world.IsPaused = !world.IsPaused
}

func (world *World) AddNewGopher(gopher *Gopher) {

	world.AddFunctionToWorldInputActions(func() {
		world.ActiveGophers <- gopher
	})

}

func (world *World) QueueRemoveGopher(gopher *Gopher) {

	world.AddFunctionToWorldInputActions(func() {
		//gopher = nil
		if mapPoint, ok := world.GetMapPoint(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			mapPoint.Gopher = nil
		}
	})
}
