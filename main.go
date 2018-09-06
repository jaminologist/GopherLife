package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"./food"
	"./maputil"
)

const hungerPerMoment = 2

type Gopher struct {
	name     string
	lifespan int
	hunger   int
}

func newGopher(name string) Gopher {

	return Gopher{name: name, lifespan: 0, hunger: rand.Intn(1000)}
}

func (g *Gopher) SetName(name string) {
	g.name = name
}

func (g *Gopher) IsDead() bool {
	return g.lifespan >= 300 || g.hunger <= 0
}

func (g *Gopher) applyHunger() {
	g.hunger -= hungerPerMoment
}

func (g *Gopher) Eat() {

	//for i := 0; i < 100; i++ {

	//}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func CreateMap(width int, height int) {

	var m = make(map[string]struct{})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var s struct{}
			m[maputil.CoordinateMapKey(x, y)] = s
		}
	}

}

func PerformMoment(wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	//	fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	if !g.IsDead() {
		g.lifespan++
		g.applyHunger()
		g.Eat()
		c <- g
	} else {
		//fmt.Println(g.name, "is dead :(")
	}
	wg.Done()

}

func generateGophers(inputChannel chan *Gopher) {
	goph := newGopher(fmt.Sprintf("[1]", rand.Int()))

	select {
	case inputChannel <- &goph:
		generateGophers(inputChannel)
	}

}

func a(inputChannel chan *Gopher) {
	for len(inputChannel) != cap(inputChannel) {
		go generateGophers(inputChannel)
	}
}

func main() {

	//	runtime.GOMAXPROCS(1)
	start := time.Now()

	var wg sync.WaitGroup

	var numGophers = 10000

	channel := make(chan *Gopher, numGophers)

	for len(channel) != cap(channel) {
		go generateGophers(channel)
	}

	i := 0

	for numGophers > 0 {

		fmt.Println("Num Gophs: ", numGophers, " Number of moments: ", i)
		numGophers = len(channel)
		secondChannel := make(chan *Gopher, numGophers)
		for i := 0; i < numGophers; i++ {
			msg := <-channel
			wg.Add(1)
			go PerformMoment(&wg, msg, secondChannel)
		}
		channel = secondChannel
		i++
		wg.Wait()
	}

	fmt.Println(time.Since(start))
	fmt.Println("Done")

	fmt.Println(food.NewCarrot().Name)
}
