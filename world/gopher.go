package world

import (
	food "gopherlife/food"
	math "gopherlife/math"
	"math/rand"
	"sync"
	"time"
)

const hungerPerMoment = 0
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

			if moveX == 0 && moveY == 0 {
				world.InputActions <- g.QueuePickUpFood(world)
				break
			}

			world.InputActions <- g.QueueMovement(world, moveX, moveY)

		default:
			g.FindFood(world, 10)
		}
	}

	if !g.IsDead() {
		g.Lifespan++
		g.ApplyHunger()
	}

	if !g.IsDecayed() {
		channel <- g
	} else {
		world.InputActions <- g.QueueRemoveGopher(world)
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

func (gopher *Gopher) QueueRemoveGopher(world *World) func() {

	return func() {
		if mapPoint, ok := world.world[gopher.Position.MapKey()]; ok {
			mapPoint.Gopher = nil
		}
	}

}
