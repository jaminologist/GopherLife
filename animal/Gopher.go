package animal

import (
	food "gopherlife/food"
	math "gopherlife/math"
	"math/rand"
	"time"
)

const hungerPerMoment = 0

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
	return g.Lifespan >= 300 || g.Hunger <= 0
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
