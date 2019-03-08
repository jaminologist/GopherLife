package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
)

const hungerPerMoment = 1
const timeToDecay = 10

type Gender int

//Gopher The whole point of the project
type Gopher struct {
	Name string

	Lifespan                   int
	Decay                      int
	Hunger                     int
	CounterTillReadyToFindLove int

	IsMated  bool
	IsDead   bool
	IsHungry bool

	Position calc.Coordinates

	HeldFood *Food

	Gender Gender

	FoodTargets   []calc.Coordinates
	GopherTargets []calc.Coordinates
	MovementPath  []calc.Coordinates
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

func (gopher *Gopher) IsLookingForLove() bool {
	return gopher.IsMature() && !gopher.IsMated
}

func (gopher *Gopher) IsTakingABreakFromLove() bool {
	return gopher.IsMature() && gopher.IsMated && gopher.CounterTillReadyToFindLove < 100
}

func (gopher *Gopher) IsDecayed() bool {
	return gopher.Decay >= timeToDecay
}

func (gopher *Gopher) ApplyHunger() {
	gopher.Hunger -= hungerPerMoment
}

func (gopher *Gopher) Eat() {
	if gopher.HeldFood != nil {
		gopher.Hunger += gopher.HeldFood.Energy
		gopher.HeldFood = nil
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

//ClearFoodTargets Clears all food targets from the Gopher
func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []calc.Coordinates{}
}

func (gopher *Gopher) ClearGopherTargets() {
	gopher.GopherTargets = []calc.Coordinates{}
}

type GopherActor struct {
	ActionQueuer
	Searchable
	GopherContainer
	FoodContainer
	PickableTiles
	MoveableGophers
	ActorGeneration
	GopherBirthRate int
}

func (actor *GopherActor) Update(gopher *Gopher) {
	switch {
	case gopher.IsDead:
		gopher.Decay++
	case gopher.IsHungry:
		actor.handleHunger(gopher)
	case !gopher.IsHungry:

		switch {
		case gopher.Gender == Male:
			if gopher.IsLookingForLove() {
				gopher.GopherTargets = actor.Search(gopher.Position, 15, 15, 1, SearchForFemaleGopher)
				if len(gopher.GopherTargets) <= 0 {
					actor.Wander(gopher)
				} else {

					target := gopher.GopherTargets[0]

					if gopher.Position.IsInRange(target, 1, 1) {
						actor.QueueMating(gopher, target)
						break
					}
					moveX, moveY := calc.FindNextStep(gopher.Position, target)
					actor.QueueGopherMove(moveX, moveY, gopher)
					gopher.ClearFoodTargets()
				}
			} else {
				actor.Wander(gopher)
			}

		default:
			actor.Wander(gopher)
		}

	}

	gopher.AdvanceLife()

}

func (actor *GopherActor) handleHunger(gopher *Gopher) {
	switch {
	case gopher.HeldFood != nil:
		gopher.Eat()
	default:
		actor.moveTowardsFood(gopher)
	}
}

func (actor *GopherActor) moveTowardsFood(gopher *Gopher) {

	if len(gopher.FoodTargets) > 0 {

		target := gopher.FoodTargets[0]

		if _, ok := actor.HasFood(target.GetX(), target.GetY()); ok {

			if gopher.Position.IsInRange(target, 0, 0) {
				actor.QueuePickUpFood(gopher)
				gopher.ClearFoodTargets()
				return
			}

			moveX, moveY := calc.FindNextStep(gopher.Position, target)
			actor.QueueGopherMove(moveX, moveY, gopher)
		} else {
			gopher.ClearFoodTargets()
		}

	} else {
		actor.LookForFood(gopher)
	}

}

//Wander Randomly decides a diretion for the gopher to move in
func (actor *GopherActor) Wander(gopher *Gopher) {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	actor.QueueGopherMove(x, y, gopher)
}

func (actor *GopherActor) LookForFood(gopher *Gopher) {
	gopher.FoodTargets = actor.Search(gopher.Position, 25, 25, 1, SearchForFood)
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (actor *GopherActor) QueueGopherMove(moveX int, moveY int, gopher *Gopher) {

	actor.Add(func() {
		success := actor.MoveGopher(gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (actor *GopherActor) QueuePickUpFood(gopher *Gopher) {

	actor.Add(func() {
		food, ok := actor.PickUpFood(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			gopher.ClearFoodTargets()
		}
	})
}

func (actor *GopherActor) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	actor.Add(func() {

		if mate, ok := actor.HasGopher(matePosition.GetX(), matePosition.GetY()); ok {

			litterNumber := 0

			if actor.GopherBirthRate > 0 {
				litterNumber = rand.Intn(actor.GopherBirthRate)
			}

			emptySpaces := actor.Search(gopher.Position, 10, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {

						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])
						actor.AddNewGopher(pos.GetX(), pos.GetY(), &newborn)

					}

				}

			}

		}
	})

}
