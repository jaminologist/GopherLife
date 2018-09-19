package main

import (
	"fmt"
	food "gopherlife/food"
	gopherlife "gopherlife/world"
	"os"
	"strconv"
	"time"
)

const numberOfGophs = 5

const size = 30

func main() {

	//	runtime.GOMAXPROCS(1)
	start := time.Now()

	width, height := size, size
	var world = gopherlife.CreateWorld(width, height)

	for len(world.ActiveGophers) > 0 {

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
				world.ProcessWorld()
				time--
			}

		}
	}

	fmt.Println(time.Since(start))
	fmt.Println("The World Lasted for ", world.Moments, " Moments")
	fmt.Println("Done")

	fmt.Println(food.NewCarrot().Name)
}
