package world

import (
	animal "gopherlife/animal"
	"gopherlife/food"
	"gopherlife/math"
	"math/rand"
	"strconv"
	"sync"
)

const numberOfGophs = 100
const numberOfFoods = 500
const worldSize = 200

type World struct {
	world map[string]*MapPoint

	width  int
	height int

	InputActions chan func()
	OutputAction chan func()

	GopherWaitGroup sync.WaitGroup

	ActiveGophers chan *animal.Gopher

	//RenderOutput chan string

	SelectedGopher *animal.Gopher

	StartX int
	StartY int

	IsPaused bool

	RenderSize int

	Moments int
}

func CreateMap(width int, height int) map[string]*animal.Gopher {

	var m = make(map[string]*animal.Gopher)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			m[math.CoordinateMapKey(x, y)] = &animal.Gopher{}
		}
	}

	return m

}

func CreateWorld() World {

	world := World{width: worldSize, height: worldSize}
	world.InputActions = make(chan func(), 10000)
	world.OutputAction = make(chan func(), 10000)

	world.IsPaused = false

	world.world = make(map[string]*MapPoint)

	world.RenderSize = 25

	for x := 0; x < worldSize; x++ {
		for y := 0; y < worldSize; y++ {
			var point = MapPoint{}
			world.world[math.CoordinateMapKey(x, y)] = &point
		}
	}

	world.SetUpMapPoints(numberOfGophs, numberOfFoods)

	return world

}

func (world *World) SelectEntity(mapKey string) {

	if mapPoint, ok := world.world[mapKey]; ok {
		if mapPoint.Gopher != nil {
			world.SelectedGopher = mapPoint.Gopher
		}
	}

}

func (world *World) RenderWorld() string {

	renderString := ""

	startX := 0
	startY := 0

	if world.SelectedGopher != nil {
		startX = world.SelectedGopher.Position.GetX() - world.RenderSize/2
		startY = world.SelectedGopher.Position.GetY() - world.RenderSize/2
	} else {
		startX = world.StartX
		startY = world.StartY
	}

	renderString += "<h1>"
	for y := startY; y < startY+world.RenderSize; y++ {

		for x := startX; x < startX+world.RenderSize; x++ {

			key := math.CoordinateMapKey(x, y)

			if mapPoint, ok := world.world[key]; ok {

				switch {
				case mapPoint.isEmpty():
					renderString += "<span class='grass'>O</span>"
				case mapPoint.Gopher != nil:

					isSelected := false

					if world.SelectedGopher != nil {
						isSelected = world.SelectedGopher.Position.MapKey() == key
					}
					if isSelected {
						renderString += "<span id='" + key + "' style='color:yellow;'>G</span>"
					} else if mapPoint.Gopher.IsDead() {
						renderString += "<span id='" + key + "' class='gopher'>X</span>"
					} else {
						renderString += "<span id='" + key + "' class='gopher'>G</span>"
					}

				case mapPoint.Food != nil:
					renderString += "<span class='food'>F</span>"
				}
			} else {
				renderString += "<span class='food'>X</span>"
			}

		}
		renderString += "<br />"
	}

	renderString += "</h1>"

	return renderString
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

	world.ActiveGophers = make(chan *animal.Gopher, numberOfGophers)

	for i := 0; i < numberOfGophers; i++ {
		var mapPoint = world.world[keys[count]]

		var goph = animal.NewGopher(strconv.Itoa(i), math.StringToCoordinates(keys[count]))

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

func FindFood(world *World, g *animal.Gopher, radius int) {

	//fmt.Println("Gopher ", g.Name, " is Searching for Food...")

	var coordsArray = []math.Coordinates{}

	for x := -radius; x < radius; x++ {
		for y := -radius; y < radius; y++ {
			if x == 0 && y == 0 {
				continue
			}

			key := g.Position.RelativeCoordinate(x, y)

			if mapPoint, ok := world.world[key.MapKey()]; ok {
				food := mapPoint.Food

				if food != nil {
					coordsArray = append(coordsArray, key)
				}
			}

		}
	}

	math.SortCoordinatesUsingCoordinate(g.Position, coordsArray)

	if len(coordsArray) > 0 {
		g.FoodTargets = coordsArray
		//fmt.Println("Gopher ", g.Name, " Found food at nearest location ", coordsArray[0])
	}

}

func QueueMovement(world *World, goph *animal.Gopher, x int, y int) func() {

	return func() {

		//currentPostion := goph.position
		currentMapPoint := world.world[goph.Position.MapKey()]

		targetPosition := goph.Position.RelativeCoordinate(x, y)
		targetMapPoint, exists := world.world[targetPosition.MapKey()]
		if exists && targetMapPoint.Gopher == nil {

			targetMapPoint.Gopher = goph
			currentMapPoint.Gopher = nil
			goph.Position = targetPosition
		} else {
			world.OutputAction <- func() {
			}
		}
	}

}

func QueuePickUpFood(world *World, gopher *animal.Gopher) func() {

	return func() {

		//currentPostion := goph.position
		currentMapPoint := world.world[gopher.Position.MapKey()]

		if currentMapPoint.Food == nil {
			gopher.FoodTargets = nil
		} else {
			gopher.HeldFood = currentMapPoint.Food
			currentMapPoint.Food = nil
			world.onFoodPickUp(gopher.Position)

		}
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

func QueueRemoveGopher(world *World, goph *animal.Gopher) func() {

	return func() {

		if mapPoint, ok := world.world[goph.Position.MapKey()]; ok {
			mapPoint.Gopher = nil
		}
		//currentPostion := goph.position
		currentMapPoint := world.world[goph.Position.MapKey()]

		if currentMapPoint.Food == nil {
			goph.FoodTargets = nil
		} else {
			goph.HeldFood = currentMapPoint.Food
			currentMapPoint.Food = nil
		}
	}

}

func PerformMoment(world *World, wg *sync.WaitGroup, g *animal.Gopher, c chan *animal.Gopher) {

	//fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	switch {
	case g.IsDead():
	case g.Hunger < 1000:

		switch {
		case g.HeldFood != nil:
			g.Eat()
		case len(g.FoodTargets) > 0:

			target := g.FoodTargets[0]

			diffX := g.Position.GetX() - target.GetX()
			diffY := g.Position.GetY() - target.GetY()

			moveX := 0
			moveY := 0

			if moveX == 0 && moveY == 0 {
				world.InputActions <- QueuePickUpFood(world, g)
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

			world.InputActions <- QueueMovement(world, g, moveX, moveY)
		default:
			FindFood(world, g, 50)
		}
	}

	if !g.IsDead() {
		g.Lifespan++
		g.ApplyHunger()
		c <- g
	} else {
		//	fmt.Println("Gopher ", g.Name, " is dead :(")
		g = nil
	}
	wg.Done()

}

func (world *World) TogglePause() {
	world.IsPaused = !world.IsPaused
}

func (world *World) ProcessWorld() {

	if world.IsPaused {
		return
	}

	numGophers := len(world.ActiveGophers)
	secondChannel := make(chan *animal.Gopher, numGophers)
	for i := 0; i < numGophers; i++ {
		gopher := <-world.ActiveGophers
		world.GopherWaitGroup.Add(1)
		go PerformMoment(world, &world.GopherWaitGroup, gopher, secondChannel)
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

}
