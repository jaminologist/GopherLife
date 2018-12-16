package world

import (
	"gopherlife/calc"
	food "gopherlife/food"
	"math/rand"
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
	return g.Lifespan >= 50
}

func (gopher *Gopher) SetIsDead() {

	if gopher.IsDead {
		return
	}

	chance := rand.Intn(101)

	a := gopher.Lifespan - 500

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

		if mapPoint, ok := world.GetMapPoint(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if mapPointCheck(mapPoint) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(g.Position, coordsArray)

	return coordsArray

}

func (g *Gopher) moveTowardsFood(world *World) {

	if len(g.FoodTargets) > 0 {

		target := g.FoodTargets[0]

		mapPoint, _ := world.GetMapPoint(target.GetX(), target.GetY())

		if mapPoint.Food == nil {
			g.ClearFoodTargets()
		} else {

			if g.Position.IsInRange(target, 0, 0) {
				world.QueuePickUpFood(g)
				return
			}

			moveX, moveY := calc.FindNextStep(g.Position, target)
			world.QueueGopherMove(g, moveX, moveY)

		}

	} else {
		g.LookForFood(world)
	}

}

func (g *Gopher) LookForFood(world *World) {
	g.FoodTargets = g.Find(world, 25, 1, CheckMapPointForFood)
}

func (g *Gopher) handleHunger(world *World) {
	switch {
	case g.HeldFood != nil:
		g.Eat()
	default:
		g.moveTowardsFood(world)
	}
}

func (g *Gopher) PerformMoment(world *World) {

	switch {
	case g.IsDead:
		g.Decay++
	case g.IsHungry:
		g.handleHunger(world)
	case !g.IsHungry:

		switch {
		case g.Gender == Male:

			if g.IsLookingForLove() {
				g.GopherTargets = g.Find(world, 15, 1, g.CheckMapPointForPartner)
				if len(g.GopherTargets) <= 0 {
					g.Wander(world)
				} else {
					target := g.GopherTargets[0]

					if g.Position.IsInRange(target, 1, 1) {
						world.QueueMating(g, target)
						break
					}
					moveX, moveY := calc.FindNextStep(g.Position, target)
					world.QueueGopherMove(g, moveX, moveY)
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

func (gopher *Gopher) Wander(world *World) {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	world.QueueGopherMove(gopher, x, y)
}

func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []calc.Coordinates{}
}
