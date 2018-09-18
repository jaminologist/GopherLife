package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"./food"
	"./maputil"
)

const hungerPerMoment = 2
const numberOfGophs = 5

const size = 30

type MapPoint struct {
	Gopher *Gopher
	Food   *food.Food
}

func (mp *MapPoint) isEmpty() bool {
	return mp.Gopher == nil && mp.Food == nil
}

type World struct {
	world map[string]*MapPoint

	width  int
	height int

	inputActions chan func()
	outputAction chan func()

	activeGophers chan *Gopher
}

func (world *World) RenderWorld() {

	for k := range world.world {

		if !world.world[k].isEmpty() {

			mapGopher := world.world[k].Gopher
			var a string
			if mapGopher == nil {
				a = "NOT SET"
			} else {
				a = mapGopher.name
			}

			mapFood := world.world[k].Food
			var b string

			if mapFood == nil {
				b = "NOT SET"
			} else {
				b = mapFood.Name
			}

			fmt.Println("At ", k, " Gopher is: ", a, "Food is: ", b)
		}
	}
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

	world.activeGophers = make(chan *Gopher, numberOfGophers)

	for i := 0; i < numberOfGophers; i++ {
		var mapPoint = world.world[keys[count]]

		var goph = newGopher(strconv.Itoa(i), maputil.StringToCoordinates(keys[count]))

		mapPoint.Gopher = &goph
		world.activeGophers <- &goph
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

	heldFood *food.Food

	foodTargets []maputil.Coordinates

	movementPath []maputil.Coordinates
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

	if g.heldFood != nil {

		prev := g.hunger
		g.hunger += g.heldFood.Energy

		foodName := g.heldFood.Name

		g.heldFood = nil

		fmt.Println("Gopher "+g.name, " is eating a ",
			foodName, ". Hunger restored to ", g.hunger, " from ", prev)
	}
	//for i := 0; i < 100; i++ {

	//}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func (g *Gopher) FindFood(world *World, radius int) {

	fmt.Println("Gopher ", g.name, " is Searching for Food...")

	var coordsArray = []maputil.Coordinates{}

	for x := 0; x < radius; x++ {
		for y := 0; y < radius; y++ {
			if x == 0 && y == 0 {
				continue
			}

			key := g.position.RelativeCoordinate(x, y)

			if mapPoint, ok := world.world[key.MapKey()]; ok {
				food := mapPoint.Food

				if food != nil {
					coordsArray = append(coordsArray, key)
				}
			}

		}
	}

	sort.Sort(maputil.ByNearest(coordsArray))

	if len(coordsArray) > 0 {
		g.foodTargets = coordsArray
		fmt.Println("Gopher ", g.name, " Found food at nearest location ", coordsArray[0])
	}

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

func QueuePickUpFood(world *World, goph *Gopher) func() {

	return func() {

		//currentPostion := goph.position
		currentMapPoint := world.world[goph.position.MapKey()]

		if currentMapPoint.Food == nil {
			goph.foodTargets = nil
		} else {
			goph.heldFood = currentMapPoint.Food
			currentMapPoint.Food = nil
		}
	}

}

func QueueRemoveGopher(world *World, goph *Gopher) func() {

	return func() {

		if mapPoint, ok := world.world[goph.position.MapKey()]; ok {
			mapPoint.Gopher = nil
		}
		//currentPostion := goph.position
		currentMapPoint := world.world[goph.position.MapKey()]

		if currentMapPoint.Food == nil {
			goph.foodTargets = nil
		} else {
			goph.heldFood = currentMapPoint.Food
			currentMapPoint.Food = nil
		}
	}

}

func PerformMoment(world *World, wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	//fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	switch {
	case g.IsDead():
	case g.hunger < 1000:

		switch {
		case g.heldFood != nil:
			g.Eat()
		case len(g.foodTargets) > 0:
			target := g.foodTargets[0]

			diffX := g.position.GetX() - target.GetX()
			diffY := g.position.GetY() - target.GetY()

			moveX := 0
			moveY := 0

			if moveX == 0 && moveY == 0 {
				world.inputActions <- QueuePickUpFood(world, g)
			}

			if diffX > 0 {
				moveX = -1
			} else if diffX < 0 {
				moveX = 1
			}

			if diffY > 0 {
				moveY = -1
			} else if diffY < 0 {
				moveY = 1
			}

			world.inputActions <- QueueMovement(world, g, moveX, moveY)
		default:
			g.FindFood(world, 100)
		}
	}

	if !g.IsDead() {
		g.lifespan++
		g.applyHunger()
		c <- g
	} else {
		fmt.Println("Gopher ", g.name, " is dead :(")
		g = nil
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
	world.RenderWorld()
	//fmt.Println(world)

	var wg sync.WaitGroup

	var numGophers = 1

	var channel = world.activeGophers

	for numGophers > 0 {

		fmt.Println("Enter Command...")
		var input, value string
		fmt.Scanln(&input, &value)

		switch input {
		case "at", "advancetime":
			time, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			for time > 0 {

				numGophers = len(channel)
				secondChannel := make(chan *Gopher, numGophers)
				for i := 0; i < numGophers; i++ {
					msg := <-channel
					wg.Add(1)
					go PerformMoment(&world, &wg, msg, secondChannel)
				}
				channel = secondChannel
				time--
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

		}
	}

	fmt.Println(time.Since(start))
	fmt.Println("Done")

	fmt.Println(food.NewCarrot().Name)
}
