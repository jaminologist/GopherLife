package main

import (
	"fmt"
	"strconv"
	"sync"
)

type Gopher struct {
	name     string
	lifetime int
}

func (gopher Gopher) Live(aliveChannel chan Gopher, wg *sync.WaitGroup) {

	gopher.lifetime++

	if gopher.lifetime > 30000 {
		fmt.Println("I, Gopher ", gopher.name, " am now the dead :(")
	} else {
		aliveChannel <- gopher
	}

	wg.Done()
}

func main() {

	var wg sync.WaitGroup

	numGophers := 10

	inputChan, outputChan := make(chan Gopher, numGophers), make(chan Gopher, numGophers)

	for i := 0; i < numGophers; i++ {
		newGopher := Gopher{name: strconv.Itoa(i)}
		inputChan <- newGopher
		fmt.Println("I, Gopher ", newGopher.name, " am alive! :)")
	}

	for numGophers > 0 {

		for i := 0; i < numGophers; i++ {
			wg.Add(1)
			gopher := <-inputChan
			go gopher.Live(outputChan, &wg)
		}

		wg.Wait()

		numGophers = len(outputChan)
		inputChan = outputChan
		outputChan = make(chan Gopher, numGophers)

	}

}
