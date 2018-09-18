package main

import (
	"fmt"
	"gopherlife/animal"
	food "gopherlife/food"
	gopherlife "gopherlife/world"
	"os"
	"strconv"
	"sync"
	"time"
)

const numberOfGophs = 5

const size = 30

func main() {

	//	runtime.GOMAXPROCS(1)
	start := time.Now()

	width, height := size, size

	var world = gopherlife.CreateWorld(width, height)
	world.RenderWorld()
	//fmt.Println(world)

	var wg sync.WaitGroup

	var numGophers = 1

	var channel = world.ActiveGophers

	for numGophers > 0 {

		fmt.Println("Enter Command...")
		var input, value string
		fmt.Scanln(&input, &value)

		switch input {
		case "at", "advancetime":
			time, err := strconv.Atoi(value)
			if err != nil {
				fmt.Println(err)
				os.Exit(2)
			}

			for time > 0 {

				numGophers = len(channel)
				secondChannel := make(chan *animal.Gopher, numGophers)
				for i := 0; i < numGophers; i++ {
					msg := <-channel
					wg.Add(1)
					go gopherlife.PerformMoment(&world, &wg, msg, secondChannel)
				}
				channel = secondChannel
				time--
				wg.Wait()

				wait := true
				for wait {
					select {
					case action := <-world.InputActions:
						action()
					default:
						wait = false
					}
				}

				wait = true

				for wait {
					select {
					case action := <-world.OutputAction:
						action()
					default:
						wait = false
					}
				}

			}

		}
	}

	fmt.Println(time.Since(start))
	fmt.Println("Done")

	fmt.Println(food.NewCarrot().Name)
}
