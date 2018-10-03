package main

import (
	"fmt"
	"strconv"
	"time"
)

type Gopher struct {
	name     string
	lifetime int
	isDead   bool
}

func (gopher Gopher) Live() {

	fmt.Println("I, Gopher ", gopher.name, " am alive! :)")

	gopher.lifetime++

	if gopher.lifetime > 300 {
		gopher.isDead = true
	}

	fmt.Println("I, Gopher ", gopher.name, " am now the dead :(")
}

func main() {

	input := make(chan func(), 10000)
	output := make(chan func(), 10000)

	for i := 0; i < 10; i++ {
		gopher := Gopher{name: strconv.Itoa(i)}
		go gopher.Live()
	}

	time.Sleep(time.Duration(2) * time.Second)

}
