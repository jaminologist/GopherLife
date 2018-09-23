package main

import (
	"fmt"
	gopherlife "gopherlife/world"
	"net/http"
	"time"
)

const numberOfGophs = 5

const size = 1000

func main() {

	//runtime.GOMAXPROCS(1)
	start := time.Now()

	var world = gopherlife.CreateWorld()

	for len(world.ActiveGophers) > 0 {
		world.ProcessWorld()
		//mux := http.NewServeMux()

		http.HandleFunc("/", worldToHTML(&world))
		//mux.HandleFunc("/", worldToHTML(&world))
		http.ListenAndServe(":8080", nil)
	}

	fmt.Println(time.Since(start))
	fmt.Println("The World Lasted for ", world.Moments, " Moments")
	fmt.Println("Done")
}

func worldToHTML(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, world.RenderWorld())
	}

}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<p> Hello </p>")
	fmt.Fprintln(w, "<p> Hello  Again</p>")
	fmt.Fprintln(w, "<p> Hello WURTU Again</p>")
}
