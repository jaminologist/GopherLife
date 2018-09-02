package main

import (
	"fmt"
	"sync"
	"time"
)

type Gopher struct {
	name     string
	lifespan int

	hunger int
}

func (g *Gopher) SetName(name string) {
	g.name = name
}

func (g *Gopher) IsDead() bool {
	return g.lifespan >= 4
}

func (g *Gopher) Eat() {

	for i := 0; i < 1000000; i++ {

	}

}

func (g *Gopher) Dig() {
	time.Sleep(100)
}

func StartLife(wg *sync.WaitGroup, g *Gopher, c chan *Gopher) {

	fmt.Println(g.name, "is alive :) lifespan is: ", g.lifespan)

	if !g.IsDead() {
		g.lifespan++
		g.Eat()
		c <- g
	} else {
		fmt.Println(g.name, "is dead :(")
	}
	wg.Done()

}

func main() {

	//	runtime.GOMAXPROCS(1)

	start := time.Now()

	var wg sync.WaitGroup

	var numGophers = 10000

	channel := make(chan *Gopher, numGophers)

	for i := 0; i < numGophers; i++ {
		wg.Add(1)
		go StartLife(&wg, &Gopher{fmt.Sprintf("%v", i), 1, 1}, channel)
	}

	for numGophers > 0 {
		wg.Wait()
		numGophers = len(channel)
		secondChannel := make(chan *Gopher, numGophers)
		for i := 0; i < numGophers; i++ {
			msg := <-channel
			wg.Add(1)
			go StartLife(&wg, msg, secondChannel)
		}
		channel = secondChannel
	}

	fmt.Println("Done")

	fmt.Println(time.Since(start))

}
