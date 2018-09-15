package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"./food"
	"./maputil"
)

const hungerPerMoment = 2
const numberOfGophs = 1

const size = 30

type MapPoint struct {
	Gopher *Gopher
	Food   *food.Food
}

func (mp *MapPoint) isEmpty() bool {
	return mp.Gopher == &Gopher{} && mp.Food == &food.Food{}
}

type World struct {
	world map[string]*MapPoint

	width  int
	height int

	inputActions chan func()
	outputAction chan func()

	activeGophers chan *Gopher
}

func (world *World) SetUpMapPoints(numberOfGophers int, numberOfFood int) {

	keys := make([]string, len(world.world))

	i := 0
	for k := range world.world {
		keys[i] = k
		i++
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	count := 0

	world.activeGophers = make(chan *Gopher, numberOfGophers)

	for i := 0; i < numberOfGophers; i++ {
		var mapPoint = world.world[keys[count]]

		var goph = newGopher(strconv.Itoa(i), maputil.StringToCoordinates(keys[count]))

		mapPoint.Gopher = &goph
		world.activeGophers <- &goph
		world.world[keys[count]] = mapPoint
	}

	for i := 0; i < numberOfFood; i++ {
		var mapPoint = world.world[keys[count]]

		var food = food.NewPotato()

		mapPoint.Food = &food
		world.world[keys[count]] = mapPoint
	}

}

func CreateWorld(width int, height int) World {

	world := World{width: width, height: height}
	world.inputActions = make(chan func(), 10000)
	world.outputAction = make(chan func(), 10000)

	world.world = make(map[string]*MapPoint)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var point = MapPoint{}
			world.world[maputil.CoordinateMapKey(x, y)] = &point
		}
	}

	world.inputActions <- func() {
		fmt.Println("Hello")
	}

	world.SetUpMapPoints(numberOfGophs, 100)

	return world

}

type Gopher struct {
	name     string
	lifespan int
	hunger   int

	position maputil.Coordinates
}

func newGopher(name string, coord maputil.Coordinates) Gopher {
	return Gopher{name: name, lifespan: 0, hunger: rand.Intn(100), position: coord}
}

func (g *Gopher) SetName(name string) {
	g.name = name
}

func (g *Gopher) IsDead() bool {
	return g.lifespan >= 300 || g.hunger <= 0
}

func (g *Gopher) applyHunger() {
	g.hunger -= hungerPerMoment
}

func (g *Gopher) Move(x int, y int) {
	g.position.Add(x, y)
}

func (g *Gopher) Eat() {

	//for i := 0; i < 100; i++ {

	//}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func CreateMap(width int, height int) map[string]*Gopher {

	var m = make(map[string]*Gopher)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			m[maputil.CoordinateMapKey(x, y)] = &Gopher{}
		}
	}

	return m

}

func QueueMovement(world *World, goph *Gopher, x int, y int) func() {

	return func() {

		//currentPostion := goph.position
		currentMapPoint := world.world[goph.position.MapKey()]

		targetPosition := goph.position.RelativeCoordinate(x, y)
		targetMapPoint, exists := world.world[targetPosition.MapKey()]

		//	fmt.Println("Current", currentMapPoint.Gopher.name, "Position", currentMapPoint.Gopher.position)

		if exists && targetMapPoint.Gopher == nil {

			targetMapPoint.Gopher = goph
			currentMapPoint.Gopher = nil
			goph.position = targetPosition
			fmt.Println("Gopher ", goph.name, "Moves to ", goph.position.MapKey())
		} else {
			world.outputAction <- func() {
				fmt.Println("Gopher ", goph.name, " Can't Move!")
			}
		}
	}

}

func PerformMoment(world *World, wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	//fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	if !g.IsDead() {
		g.lifespan++
		g.applyHunger()
		g.Eat()

		world.inputActions <- QueueMovement(world, g, 1, 1)
		c <- g
	} else {
		//fmt.Println(g.name, "is dead :(")
	}
	wg.Done()

}

func generateGopher(inputChannel chan *Gopher, coordinates string, world map[string]*Gopher) {
	goph := newGopher(fmt.Sprintf("[1]", rand.Int()), maputil.StringToCoordinates(coordinates))

	select {
	case inputChannel <- &goph:
		//world[coordinates] = goph
	}

}

func main() {

	//	runtime.GOMAXPROCS(1)
	start := time.Now()

	width, height := size, size

	var world = CreateWorld(width, height)
	fmt.Println(world)

	var wg sync.WaitGroup

	var numGophers = 1

	var channel = world.activeGophers

	i := 0

	for numGophers > 0 {

		//	fmt.Println("Num Gophs: ", numGophers, " Number of moments: ", i)
		numGophers = len(channel)
		secondChannel := make(chan *Gopher, numGophers)
		for i := 0; i < numGophers; i++ {
			msg := <-channel
			//		fmt.Println("Name: ", msg.name, " Position: ", msg.position)
			wg.Add(1)

			go PerformMoment(&world, &wg, msg, secondChannel)
		}
		channel = secondChannel
		i++
		wg.Wait()

		wait := true
		for wait {
			select {
			case action := <-world.inputActions:
				action()
			default:
				wait = false
			}
		}

		wait = true

		for wait {
			select {
			case action := <-world.outputAction:
				action()
			default:
				wait = false
			}
		}

	}

	fmt.Println(time.Since(start))
	fmt.Println("Done")

	fmt.Println(food.NewCarrot().Name)
}
