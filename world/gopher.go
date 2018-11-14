package world

import (
	"gopherlife/calc"
	food "gopherlife/food"
	"gopherlife/names"
	"math/rand"
	"time"
)

const hungerPerMoment = 5
const timeToDecay = 10

type Gender int

type Gopher struct {
	Name     string
	Lifespan int

	Gender Gender

	Decay int

	Hunger int

	IsMated bool

	IsDead bool

	IsHungry bool

	CounterTillReadyToFindLove int

	Position calc.Coordinates

	HeldFood *food.Food

	FoodTargets []calc.Coordinates

	GopherTargets []calc.Coordinates

	MovementPath []calc.Coordinates
}

const (
	Male   Gender = 0
	Female Gender = 1
)

func (gender Gender) String() string {

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

func NewGopher(Name string, coord calc.Coordinates) Gopher {

	i := rand.Intn(len(genders))

	return Gopher{
		Name:     Name,
		Lifespan: 0,
		Hunger:   rand.Intn(100) + 50,
		Position: coord,
		Gender:   genders[i],
	}
}

func (g *Gopher) SetName(Name string) {
	g.Name = Name
}

func (g *Gopher) IsMature() bool {
	return g.Lifespan >= 100
}

func (gopher *Gopher) SetIsDead() {

	if gopher.IsDead {
		return
	}

	chance := rand.Intn(101)

	a := gopher.Lifespan - 300

	if a > chance || gopher.Hunger <= 0 {
		gopher.IsDead = true
	}

}

func (gopher *Gopher) SetIsHungry() {

	if gopher.IsHungry {

		if gopher.Hunger > 300 {
			gopher.IsHungry = false
		}

	} else {

		if gopher.Hunger < 150 {
			gopher.IsHungry = true
		}
	}

}

func (g *Gopher) IsLookingForLove() bool {
	return g.IsMature() && !g.IsMated
}

func (g *Gopher) IsTakingABreakFromLove() bool {
	return g.IsMature() && g.IsMated && g.CounterTillReadyToFindLove < 100
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

func (gopher *Gopher) AdvanceLife() {

	if !gopher.IsDead {
		gopher.Lifespan++
		gopher.ApplyHunger()

		if gopher.Gender == Female && gopher.IsTakingABreakFromLove() {
			gopher.CounterTillReadyToFindLove++

			if !gopher.IsTakingABreakFromLove() {
				gopher.IsMated = false
			}

		}

		gopher.SetIsDead()
		gopher.SetIsHungry()

	}

}

type MapPointCheck func(*MapPoint) bool

func (g *Gopher) Find(world *World, radius int, maximumFind int, mapPointCheck MapPointCheck) []calc.Coordinates {

	var coordsArray = []calc.Coordinates{}

	spiral := calc.NewSpiral(radius, radius)

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > maximumFind {
			break
		}

		if coordinates.X == 0 && coordinates.Y == 0 {
			continue
		}

		relativeCoords := g.Position.RelativeCoordinate(coordinates.X, coordinates.Y)

		if mapPoint, ok := world.world[relativeCoords.MapKey()]; ok {
			if mapPointCheck(mapPoint) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(g.Position, coordsArray)

	return coordsArray

}

func (g *Gopher) PerformMoment(world *World) {

	switch {
	case g.IsDead:
		g.Decay++

	case g.IsHungry:

		switch {
		case g.HeldFood != nil:
			g.Eat()
		case len(g.FoodTargets) > 0:

			target := g.FoodTargets[0]

			moveX, moveY := calc.FindNextStep(g.Position, target)

			if moveX == 0 && moveY == 0 {
				g.QueuePickUpFood(world)
				g.ClearFoodTargets()
				break
			}

			g.QueueMovement(world, moveX, moveY)

		case len(g.FoodTargets) < 0:
			g.FoodTargets = g.Find(world, 15, 1, g.CheckMapPointForFood)
			//If no foodtargets wander?
		default:
			g.FoodTargets = g.Find(world, 15, 1, g.CheckMapPointForFood)
		}
	case !g.IsHungry:

		switch {
		case g.Gender == Male:

			if g.IsLookingForLove() {
				g.GopherTargets = g.Find(world, 15, 1, g.CheckMapPointForPartner)
				if len(g.GopherTargets) <= 0 {
					g.Wander(world)
				} else {
					target := g.GopherTargets[0]

					diffX, diffY := g.Position.Difference(target)

					moveX, moveY := calc.FindNextStep(g.Position, target)

					/*if world.SelectedGopher.Position.Equals(&g.Position) {
						fmt.Println(target, " actual position is ", g.Position)
						fmt.Println(moveX, ",", moveY)
						fmt.Println("diff ", "(", diffX, ",", diffY, ")")
					}*/

					if calc.Abs(diffX) <= 1 && calc.Abs(diffY) <= 1 {
						g.QueueMating(world, target)
						break
					}

					g.QueueMovement(world, moveX, moveY)
				}
			} else {
				g.Wander(world)
			}

		default:
			g.Wander(world)
		}

	}

	g.AdvanceLife()
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

	world.AddFunctionToWorldInputActions(func() {

		x := rand.Intn(3) - 1
		y := rand.Intn(3) - 1

		success := world.MoveGopher(gopher, x, y)
		_ = success
	})

}

func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []calc.Coordinates{}
}

func (gopher *Gopher) QueuePickUpFood(world *World) {

	world.AddFunctionToWorldInputActions(func() {
		food, ok := world.RemoveFoodFromWorld(gopher.Position)
		if ok {
			gopher.HeldFood = food
			world.onFoodPickUp(gopher.Position)
		}
	})

}

func (gopher *Gopher) QueueMovement(world *World, x int, y int) {

	world.AddFunctionToWorldInputActions(func() {
		success := world.MoveGopher(gopher, x, y)
		_ = success
	})

}

func (gopher *Gopher) QueueMating(world *World, matePosition calc.Coordinates) {

	world.AddFunctionToWorldInputActions(func() {

		if mapPoint, ok := world.world[matePosition.MapKey()]; ok {

			mate := mapPoint.Gopher

			if mate == nil {
				return
			}

			litterNumber := rand.Intn(7)

			emptySpaces := gopher.Find(world, 10, litterNumber, gopher.CheckMapPointForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				//gopher.IsMated = true
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						newborn := NewGopher(names.GetCuteName(), emptySpaces[i])
						world.AddNewGopher(&newborn)
					}

				}

			}

		}
	})

}
