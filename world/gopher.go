package world

import (
	food "gopherlife/food"
	math "gopherlife/math"
	"math/rand"
	"sync"
	"time"
)

const hungerPerMoment = 10

type Gopher struct {
	Name     string
	Lifespan int
	Hunger   int

	Position math.Coordinates

	HeldFood *food.Food

	FoodTargets []math.Coordinates

	MovementPath []math.Coordinates
}

func NewGopher(Name string, coord math.Coordinates) Gopher {
	return Gopher{Name: Name, Lifespan: 0, Hunger: rand.Intn(100), Position: coord}
}

func (g *Gopher) SetName(Name string) {
	g.Name = Name
}

func (g *Gopher) IsDead() bool {
	return false //g.Lifespan >= 300 || g.Hunger <= 0
}

func (g *Gopher) ApplyHunger() {
	g.Hunger -= hungerPerMoment
}

func (g *Gopher) Move(x int, y int) {
	g.Position.Add(x, y)
}

func (g *Gopher) Eat() {

	if g.HeldFood != nil {

		//prev := g.Hunger
		g.Hunger += g.HeldFood.Energy

		//foodName := g.HeldFood.Name

		g.HeldFood = nil

		//fmt.Println("Gopher "+g.Name, " is eating a ",
		//foodName, ". Hunger restored to ", g.Hunger, " from ", prev)
	}
	//for i := 0; i < 100; i++ {

	//}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func (g *Gopher) FindFood(world *World, radius int) {

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
					//g.FoodTargets = coordsArray
					//return
				}
			}

		}
	}

	math.SortCoordinatesUsingCoordinate(g.Position, coordsArray)

	if len(coordsArray) > 0 {
		g.FoodTargets = coordsArray
	}

}

func (g *Gopher) PerformMoment(world *World, wg *sync.WaitGroup, channel chan *Gopher) {

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
				world.InputActions <- g.QueuePickUpFood(world)
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

			world.InputActions <- g.QueueMovement(world, moveX, moveY)
		default:
			g.FindFood(world, 10)
		}
	}

	if !g.IsDead() {
		g.Lifespan++
		g.ApplyHunger()
		channel <- g
	} else {
		g = nil
	}
	wg.Done()

}

func (gopher *Gopher) QueuePickUpFood(world *World) func() {

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

func (gopher *Gopher) QueueMovement(world *World, x int, y int) func() {

	return func() {

		//currentPostion := goph.position
		currentMapPoint := world.world[gopher.Position.MapKey()]

		targetPosition := gopher.Position.RelativeCoordinate(x, y)
		targetMapPoint, exists := world.world[targetPosition.MapKey()]
		if exists && targetMapPoint.Gopher == nil {

			targetMapPoint.Gopher = gopher
			currentMapPoint.Gopher = nil
			gopher.Position = targetPosition
		} else {
			world.OutputAction <- func() {
			}
		}
	}

}
