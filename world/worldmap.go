package world

import (
	"gopherlife/food"
	"gopherlife/math"
	"math/rand"
	"strconv"
	"sync"
)

const numberOfGophs = 1000
const numberOfFoods = 2000
const worldSize = 200

type World struct {
	world map[string]*MapPoint

	width  int
	height int

	InputActions chan func()
	OutputAction chan func()

	GopherWaitGroup sync.WaitGroup

	ActiveGophers chan *Gopher

	SelectedGopher *Gopher

	Moments int

	IsPaused bool
}

func CreateMap(width int, height int) map[string]*Gopher {

	var m = make(map[string]*Gopher)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			m[math.CoordinateMapKey(x, y)] = &Gopher{}
		}
	}

	return m

}

func CreateWorld() World {

	world := World{width: worldSize, height: worldSize}
	world.InputActions = make(chan func(), 10000)
	world.OutputAction = make(chan func(), 10000)

	world.world = make(map[string]*MapPoint)

	for x := 0; x < worldSize; x++ {
		for y := 0; y < worldSize; y++ {
			var point = MapPoint{}
			world.world[math.CoordinateMapKey(x, y)] = &point
		}
	}

	world.SetUpMapPoints(numberOfGophs, numberOfFoods)

	return world

}

func (world *World) SelectEntity(mapKey string) (*Gopher, bool) {

	if mapPoint, ok := world.world[mapKey]; ok {
		if mapPoint.Gopher != nil {
			world.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, false
}

func (world *World) SetUpMapPoints(numberOfGophers int, numberOfFood int) {

	keys := make([]string, len(world.world))

	i := 0
	for k := range world.world {
		keys[i] = k
		i++
	}

	rand.Seed(1)

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	count := 0

	world.ActiveGophers = make(chan *Gopher, numberOfGophers)

	for i := 0; i < numberOfGophers; i++ {
		var mapPoint = world.world[keys[count]]

		var goph = NewGopher(strconv.Itoa(i), math.StringToCoordinates(keys[count]))

		mapPoint.Gopher = &goph
		world.ActiveGophers <- &goph
		world.world[keys[count]] = mapPoint
		count++
	}

	for i := 0; i < numberOfFood; i++ {
		var mapPoint = world.world[keys[count]]

		var food = food.NewPotato()

		mapPoint.Food = &food
		world.world[keys[count]] = mapPoint
		count++
	}

}

func (world *World) onFoodPickUp(location math.Coordinates) {

	size := 50

	xrange := rand.Perm(size)
	yrange := rand.Perm(size)

	for i := 0; i < size; i++ {

		isDone := false

		for j := 0; j < size; j++ {
			newFoodLocation := math.NewCoordinate(
				location.GetX()+xrange[i]-size/2,
				location.GetY()+yrange[j]-size/2)

			if mapPoint, ok := world.world[newFoodLocation.MapKey()]; ok {

				if mapPoint.Food == nil {
					var food = food.NewPotato()
					world.world[newFoodLocation.MapKey()].Food = &food

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

func (world *World) ProcessWorld() bool {

	if world.IsPaused {
		return false
	}

	numGophers := len(world.ActiveGophers)

	secondChannel := make(chan *Gopher, numGophers)
	for i := 0; i < numGophers; i++ {
		gopher := <-world.ActiveGophers
		world.GopherWaitGroup.Add(1)
		go gopher.PerformMoment(world, &world.GopherWaitGroup, secondChannel)

	}

	world.ActiveGophers = secondChannel

	world.GopherWaitGroup.Wait()

	wait := true
	for wait {
		select {
		case action := <-world.InputActions:
			action()
		case action := <-world.OutputAction:
			action()
		default:
			wait = false
		}
	}

	if numGophers > 0 {
		world.Moments++
	}

	return true

}

func (world *World) TogglePause() {
	world.IsPaused = !world.IsPaused
}
