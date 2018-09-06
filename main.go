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

	position *maputil.Coordinates
}

func newGopher(name string, coord *maputil.Coordinates) Gopher {
	return Gopher{name: name, lifespan: 0, hunger: rand.Intn(1000), position: coord}
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

func CreateMap(width int, height int) map[string]struct{} {

	var m = make(map[string]struct{})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			var s struct{}
			m[maputil.CoordinateMapKey(x, y)] = s
		}
	}

	return m

}

func PerformMoment(wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	//fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	if !g.IsDead() {
		g.lifespan++
		g.applyHunger()
		g.Eat()
		//fmt.Println(g.position)
		//fmt.Println(g.hunger)
		c <- g
	} else {
		//fmt.Println(g.name, "is dead :(")
	}
	wg.Done()

}

func generateGopher(inputChannel chan *Gopher, start *maputil.Coordinates) {
	goph := newGopher(fmt.Sprintf("[1]", rand.Int()), start)

	select {
	case inputChannel <- &goph:
	}

}

func main() {

	//	runtime.GOMAXPROCS(1)
	start := time.Now()

	rand.Perm(1000)

	width := 1000
	height := 1000

	var mymap = CreateMap(width, height)

	keys := make([]string, len(mymap))

	i := 0
	for k := range mymap {
		keys[i] = k
		i++
	}

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	//fmt.Print(keys)

	var wg sync.WaitGroup

	var numGophers = 10000

	channel := make(chan *Gopher, numGophers)

	count := 0

	for len(channel) != cap(channel) {
		var c = maputil.StringToCoordinates(keys[count])
		go generateGopher(channel, &c)
		count++
	}

	i = 0

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
