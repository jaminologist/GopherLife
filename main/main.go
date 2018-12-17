package main

import (
	handlers "gopherlife/handlers"
	"math/rand"
	"time"
)

func main() {
	//runtime.GOMAXPROCS(1)
	rand.Seed(time.Now().UnixNano())
	rand.Seed(1)
	handlers.SetUpPage()
}
