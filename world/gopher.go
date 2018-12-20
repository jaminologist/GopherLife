package world

import (
	"gopherlife/calc"
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

	HeldFood *Food

	FoodTargets []calc.Coordinates

	GopherTargets []calc.Coordinates

	MovementPath []calc.Coordinates
}

//NewGopher Creates a new Gopher and the given co-ordinate
func NewGopher(Name string, coord calc.Coordinates) Gopher {

	return Gopher{
		Name:     Name,
		Lifespan: 0,
		Hunger:   rand.Intn(100) + 50,
		Position: coord,
		Gender:   GetRandomGender(),
	}
}

//SetName Sets the name of the gopher
func (gopher *Gopher) SetName(Name string) {
	gopher.Name = Name
}

//IsMature Checks if the gopher is no longer a child
func (gopher *Gopher) IsMature() bool {
	return gopher.Lifespan >= 50
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
	g.FoodTargets = Find(world, g.Position, 25, 1, CheckMapPointForFood)
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
				g.GopherTargets = Find(world, g.Position, 15, 1, g.CheckMapPointForPartner)
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

//Wander Randomly decides a
func (gopher *Gopher) Wander(world *World) {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	world.QueueGopherMove(gopher, x, y)
}

func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []calc.Coordinates{}
}
