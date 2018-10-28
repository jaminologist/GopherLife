package world

import (
	food "gopherlife/food"
	math "gopherlife/math"
	"math/rand"
	"sync"
	"time"
)

const hungerPerMoment = 2
const timeToDecay = 50

type Gender int

type Gopher struct {
	Name     string
	Lifespan int

	Gender Gender

	Decay int

	Hunger int

	Position math.Coordinates

	HeldFood *food.Food

	FoodTargets []math.Coordinates

	MovementPath []math.Coordinates
}

const (
	Male   Gender = 0
	Female Gender = 1
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func (gender Gender) String() string {
	// declare an array of strings
	// ... operator counts how many
	// items in the array (7)
	names := [...]string{
		"Male",
		"Female"}

	if gender < Male || gender > Female {
		return "Unknown"
	}

	return names[gender]
}

var genders = [2]Gender{Male, Female}

func NewGopher(Name string, coord math.Coordinates) Gopher {

	return Gopher{
		Name:     Name,
		Lifespan: 0,
		Hunger:   rand.Intn(100),
		Position: coord,
		Gender:   genders[rand.Intn(len(genders))],
	}
}

func (g *Gopher) SetName(Name string) {
	g.Name = Name
}

func (g *Gopher) IsDead() bool {
	return g.Lifespan >= 5000 || g.Hunger <= 0
}

func (g *Gopher) IsHungry() bool {
	return g.Hunger <= 250
}

func (g *Gopher) IsDecayed() bool {
	return g.Decay >= timeToDecay
}

func (g *Gopher) ApplyHunger() {
	g.Hunger -= hungerPerMoment
}

func (g *Gopher) Move(x int, y int) {
	g.Position.Add(x, y)
}

func (g *Gopher) Eat() {
	if g.HeldFood != nil {
		g.Hunger += g.HeldFood.Energy
		g.HeldFood = nil
	}
}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func (g *Gopher) FindFood(world *World, radius int) {

	var coordsArray = []math.Coordinates{}

	for x := 0; x < radius; x++ {
		for y := 0; y < x; y++ {
			if x == 0 && y == 0 {
				continue
			}

			keySlice := []math.Coordinates{}

			if x == 0 {
				keySlice = append(keySlice, g.Position.RelativeCoordinate(x, y), g.Position.RelativeCoordinate(x, -y))
			} else if y == 0 {
				keySlice = append(keySlice, g.Position.RelativeCoordinate(-x, y), g.Position.RelativeCoordinate(x, y))
			} else {
				keySlice = append(keySlice,
					g.Position.RelativeCoordinate(x, y),
					g.Position.RelativeCoordinate(-x, -y),
					g.Position.RelativeCoordinate(-x, y),
					g.Position.RelativeCoordinate(x, -y),
				)
			}

			for _, element := range keySlice {
				if mapPoint, ok := world.world[element.MapKey()]; ok {
					food := mapPoint.Food
					gopher := mapPoint.Gopher
					if food != nil && gopher == nil {
						coordsArray = append(coordsArray, element)
						//g.FoodTargets = coordsArray
						//return
					}
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
		g.Decay++

	case g.IsHungry():

		switch {
		case g.HeldFood != nil:
			g.Eat()
		case len(g.FoodTargets) > 0:

			target := g.FoodTargets[0]

			moveX, moveY := math.FindNextStep(g.Position, target)

			if moveX == 0 && moveY == 0 {
				g.QueuePickUpFood(world)
				g.ClearFoodTargets()
				break
			}

			g.QueueMovement(world, moveX, moveY)

		case len(g.FoodTargets) < 0:
			g.FindFood(world, 10)
		default:
			g.FindFood(world, 10)
		}
	case !g.IsHungry():
		g.Wander(world)
	}

	if !g.IsDead() {
		g.Lifespan++
		g.ApplyHunger()
	}

	if !g.IsDecayed() {
		channel <- g
	} else {
		g.QueueRemoveGopher(world)
	}

	wg.Done()

}

func (gopher *Gopher) Wander(world *World) {
	world.InputActions <- func() {

		x := rand.Intn(3) - 1
		y := rand.Intn(3) - 1

		success := world.MoveGopher(gopher, x, y)
		_ = success
	}
}

func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []math.Coordinates{}
}

func (gopher *Gopher) QueuePickUpFood(world *World) {

	world.InputActions <- func() {
		food, ok := world.RemoveFoodFromWorld(gopher.Position)
		if ok {
			gopher.HeldFood = food
			world.onFoodPickUp(gopher.Position)
		}
	}

}

func (gopher *Gopher) QueueMovement(world *World, x int, y int) {

	world.InputActions <- func() {
		success := world.MoveGopher(gopher, x, y)
		_ = success
	}

}

func (gopher *Gopher) QueueRemoveGopher(world *World) {

	world.InputActions <- func() {
		if mapPoint, ok := world.world[gopher.Position.MapKey()]; ok {
			mapPoint.Gopher = nil
		}
	}

}
