package world

import (
	"fmt"
	food "gopherlife/food"
	math "gopherlife/math"
	"gopherlife/names"
	"math/rand"
	"time"
)

const hungerPerMoment = 2
const timeToDecay = 1

type Gender int

type Gopher struct {
	Name     string
	Lifespan int

	Gender Gender

	Decay int

	Hunger int

	IsMated bool

	Position math.Coordinates

	HeldFood *food.Food

	FoodTargets []math.Coordinates

	GopherTargets []math.Coordinates

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

func (gender Gender) Opposite() Gender {

	switch gender {
	case Male:
		return Female
	case Female:
		return Male
	default:
		return Male
	}
}

var genders = [2]Gender{Male, Female}

func NewGopher(Name string, coord math.Coordinates) Gopher {

	i := rand.Intn(len(genders))

	return Gopher{
		Name:     Name,
		Lifespan: 0,
		Hunger:   rand.Intn(100),
		Position: coord,
		Gender:   genders[i],
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

func (g *Gopher) IsLookingForLove() bool {
	return g.Lifespan > 150 && !g.IsMated
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

type MapPointCheck func(*MapPoint) bool

func (g *Gopher) Find(world *World, radius int, check MapPointCheck) []math.Coordinates {

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
					check(mapPoint)
					food := mapPoint.Food
					gopher := mapPoint.Gopher
					if food != nil && gopher == nil {
						coordsArray = append(coordsArray, element)
					}
				}
			}
		}
	}

	math.SortCoordinatesUsingCoordinate(g.Position, coordsArray)

	return coordsArray

}

func (g *Gopher) PerformMoment(world *World) {

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
			g.FoodTargets = g.Find(world, 15, g.CheckMapPointForFood)
			//If no foodtargets wander?
		default:
			g.FoodTargets = g.Find(world, 15, g.CheckMapPointForFood)
		}
	case !g.IsHungry():

		switch {
		case g.Gender == Male && g.IsLookingForLove():
			g.GopherTargets = g.Find(world, 15, g.CheckMapPointForPartner)
			if len(g.GopherTargets) == 0 {
				g.Wander(world)
			} else {
				target := g.GopherTargets[0]
				moveX, moveY := math.FindNextStep(g.Position, target)

				if moveX == 1 || moveY == 1 {
					g.QueueMating(world, target)
					break
				}

				g.QueueMovement(world, moveX, moveY)
			}
		default:
			g.Wander(world)
		}

	}

	if !g.IsDead() {
		g.Lifespan++
		g.ApplyHunger()
	}
}

func (gopher *Gopher) CheckMapPointForFood(mapPoint *MapPoint) bool {
	return mapPoint.Food != nil && mapPoint.Gopher == nil
}

func (gopher *Gopher) CheckMapPointForPartner(mapPoint *MapPoint) bool {
	return mapPoint.Gopher != nil && mapPoint.Gopher.IsLookingForLove() && gopher.Gender.Opposite() == mapPoint.Gopher.Gender
}

func (gopher *Gopher) CheckMapPointForEmptySpace(mapPoint *MapPoint) bool {
	return mapPoint.Food == nil && mapPoint.Gopher == nil
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

func (gopher *Gopher) QueueMating(world *World, matePosition math.Coordinates) {

	world.InputActions <- func() {

		if mapPoint, ok := world.world[matePosition.MapKey()]; ok {

			mate := mapPoint.Gopher

			if mate == nil {
				return
			}

			emptySpaces := gopher.Find(world, 4, gopher.CheckMapPointForEmptySpace)

			fmt.Println("HERE")
			fmt.Println(len(emptySpaces))

			if mate.Gender == Female && len(emptySpaces) > 0 {
				gopher.IsMated = true
				mate.IsMated = true

				fmt.Println("BIRTH ME")

				newborn := NewGopher(names.GetCuteName(), emptySpaces[0])

				world.AddNewGopher(&newborn)

			}

		}
	}

}
