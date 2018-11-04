package main

import (
	"fmt"
	handlers "gopherlife/handlers"
	"math/rand"
	"time"
)

func main() {

	//runtime.GOMAXPROCS(1)
	rand.Seed(time.Now().UnixNano())
	start := time.Now()

	handlers.SetUpPage()

	fmt.Println(time.Since(start))
	fmt.Println("Done")
}
