package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Gopher struct {
	name     string
	lifespan int
	hunger   int
}

func newGopher(name string) Gopher { //Gophers live for 3 years
	return Gopher{name: name, lifespan: 0, hunger: 100}
}

func (g *Gopher) SetName(name string) {
	g.name = name
}

func (g *Gopher) IsDead() bool {
	return g.lifespan >= 300
}

func (g *Gopher) Eat() {

	for i := 0; i < 100; i++ {

	}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func StartLife(wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	//	fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	if !g.IsDead() {
		g.lifespan++
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
			go StartLife(&wg, msg, secondChannel)
		}
		channel = secondChannel
		i++
		wg.Wait()
	}

	fmt.Println(time.Since(start))
	fmt.Println("Done")

}
