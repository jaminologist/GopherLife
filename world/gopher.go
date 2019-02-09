package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
)

const hungerPerMoment = 1
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

	a := gopher.Lifespan - 5000

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

func (gopher *Gopher) moveTowardsFood(tileMap TileMap) {

	if len(gopher.FoodTargets) > 0 {

		target := gopher.FoodTargets[0]

		mapPoint, _ := tileMap.Tile(target.GetX(), target.GetY())

		if mapPoint.Food == nil {
			gopher.ClearFoodTargets()
		} else {

			if gopher.Position.IsInRange(target, 0, 0) {
				tileMap.QueuePickUpFood(gopher)
				gopher.ClearFoodTargets()
				return
			}

			moveX, moveY := calc.FindNextStep(gopher.Position, target)
			tileMap.QueueGopherMove(gopher, moveX, moveY)

		}

	} else {
		gopher.LookForFood(tileMap)
	}

}

func (gopher *Gopher) LookForFood(tileMap TileMap) {
	gopher.FoodTargets = tileMap.Search(gopher.Position, 25, 25, 1, SearchForFood)
}

func (gopher *Gopher) handleHunger(tileMap TileMap) {
	switch {
	case gopher.HeldFood != nil:
		gopher.Eat()
	default:
		gopher.moveTowardsFood(tileMap)
	}
}

func (gopher *Gopher) PerformMoment(tileMap TileMap) {

	switch {
	case gopher.IsDead:
		gopher.Decay++
	case gopher.IsHungry:
		gopher.handleHunger(tileMap)
	case !gopher.IsHungry:

		switch {
		case gopher.Gender == Male:
			if gopher.IsLookingForLove() {
				gopher.GopherTargets = tileMap.Search(gopher.Position, 15, 15, 1, SearchForFemaleGopher)
				if len(gopher.GopherTargets) <= 0 {
					gopher.Wander(tileMap)
				} else {

					target := gopher.GopherTargets[0]

					if gopher.Position.IsInRange(target, 1, 1) {
						tileMap.QueueMating(gopher, target)
						break
					}
					moveX, moveY := calc.FindNextStep(gopher.Position, target)
					tileMap.QueueGopherMove(gopher, moveX, moveY)
					gopher.ClearFoodTargets()
				}
			} else {
				gopher.Wander(tileMap)
			}

		default:
			gopher.Wander(tileMap)
		}

	}

	gopher.AdvanceLife()

}

//Wander Randomly decides a diretion for the gopher to move in
func (gopher *Gopher) Wander(tileMap TileMap) {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	tileMap.QueueGopherMove(gopher, x, y)
}

//ClearFoodTargets Clears all food targets from the Gopher
func (gopher *Gopher) ClearFoodTargets() {
	gopher.FoodTargets = []calc.Coordinates{}
}

func (gopher *Gopher) ClearGopherTargets() {
	gopher.GopherTargets = []calc.Coordinates{}
}

type GopherActor struct {
	*Gopher
	QueueableActions
	Searchable
	TileContainer
	Insertable
	PickableTiles
	MoveableActors
	ActorGeneration
	GopherBirthRate int
}

func (gopher *GopherActor) Update() {

	switch {
	case gopher.IsDead:
		gopher.Decay++
	case gopher.IsHungry:
		gopher.handleHunger()
	case !gopher.IsHungry:

		switch {
		case gopher.Gender == Male:
			if gopher.IsLookingForLove() {
				gopher.GopherTargets = gopher.Search(gopher.Position, 15, 15, 1, SearchForFemaleGopher)
				if len(gopher.GopherTargets) <= 0 {
					gopher.Wander()
				} else {

					target := gopher.GopherTargets[0]

					if gopher.Position.IsInRange(target, 1, 1) {
						gopher.QueueMating(target)
						break
					}
					moveX, moveY := calc.FindNextStep(gopher.Position, target)
					gopher.QueueGopherMove(moveX, moveY)
					gopher.ClearFoodTargets()
				}
			} else {
				gopher.Wander()
			}

		default:
			gopher.Wander()
		}

	}

	gopher.AdvanceLife()

}

func (gopher *GopherActor) handleHunger() {
	switch {
	case gopher.HeldFood != nil:
		gopher.Eat()
	default:
		gopher.moveTowardsFood()
	}
}

func (gopher *GopherActor) moveTowardsFood() {

	if len(gopher.FoodTargets) > 0 {

		target := gopher.FoodTargets[0]

		mapPoint, _ := gopher.Tile(target.GetX(), target.GetY())

		if mapPoint.Food == nil {
			gopher.ClearFoodTargets()
		} else {

			if gopher.Position.IsInRange(target, 0, 0) {
				gopher.QueuePickUpFood()
				gopher.ClearFoodTargets()
				return
			}

			moveX, moveY := calc.FindNextStep(gopher.Position, target)
			gopher.QueueGopherMove(moveX, moveY)

		}

	} else {
		gopher.LookForFood()
	}

}

//Wander Randomly decides a diretion for the gopher to move in
func (gopher *GopherActor) Wander() {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	gopher.QueueGopherMove(x, y)
}

func (gopher *GopherActor) LookForFood() {
	gopher.FoodTargets = gopher.Search(gopher.Position, 25, 25, 1, SearchForFood)
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (gopher *GopherActor) QueueRemoveGopher() {

	gopher.Add(func() {
		gopher.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
	})
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (gopher *GopherActor) QueueGopherMove(moveX int, moveY int) {

	gopher.Add(func() {
		success := gopher.MoveGopher(gopher.Gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (gopher *GopherActor) QueuePickUpFood() {

	gopher.Add(func() {
		food, ok := gopher.PickUpFood(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			gopher.ClearFoodTargets()
		}
	})
}

func (gopher *GopherActor) QueueMating(matePosition calc.Coordinates) {

	gopher.Add(func() {

		if mapPoint, ok := gopher.Tile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(gopher.GopherBirthRate)

			emptySpaces := gopher.Search(gopher.Position, 10, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {

						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])

						var gopherActor = GopherActor{
							Gopher:           &newborn,
							GopherBirthRate:  gopher.GopherBirthRate,
							QueueableActions: gopher.QueueableActions,
							Searchable:       gopher.Searchable,
							TileContainer:    gopher.TileContainer,
							Insertable:       gopher.Insertable,
							PickableTiles:    gopher.PickableTiles,
							MoveableActors:   gopher.MoveableActors,
							ActorGeneration:  gopher.ActorGeneration,
						}

						gopher.AddNewGopher(pos.GetX(), pos.GetY(), &gopherActor)

					}

				}

			}

		}
	})

}
