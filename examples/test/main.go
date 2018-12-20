package main

import (
	"fmt"
	"sync"
)

func main() {

	slice := make([]int, 5)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for i := 0; i < 10000; i++ {
			slice[0] = 5
		}

		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 10000; i++ {
			slice[1] = 10
		}

		wg.Done()
	}()

	wg.Wait()

	fmt.Println("done")
}
