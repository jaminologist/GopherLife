package main

import (
	"fmt"
	gopherlife "gopherlife/world"
	"html/template"
	"log"
	"net/http"
	"time"
)

const numberOfGophs = 5

const size = 1000

func main() {

	//runtime.GOMAXPROCS(1)
	start := time.Now()

	var world = gopherlife.CreateWorld()
	renderer := gopherlife.NewRenderer()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&world))
	http.HandleFunc("/PollWorld", ajaxProcessWorld(&world, &renderer))
	http.HandleFunc("/ShiftWorldView", ajaxHandleWorldInput(&world, &renderer))
	http.HandleFunc("/SelectGopher", ajaxSelectGopher(&world, &renderer))

	fmt.Println("Listening")
	http.ListenAndServe(":8080", nil)

	fmt.Println(time.Since(start))
	fmt.Println("Done")
}

func worldToHTML(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	renderer := gopherlife.NewRenderer()

	return func(w http.ResponseWriter, r *http.Request) {

		pageVariables := PageVariables{
			Data: template.HTML(renderer.RenderWorld(world)),
		}

		t, err := template.ParseFiles("static/index.html")

		if err != nil {
			log.Print("Template parsing error: ", err)
		}

		err = t.Execute(w, pageVariables)

		if err != nil {
			log.Printf("Template executing error: ", err)
		}

	}

}

func ajaxProcessWorld(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		world.ProcessWorld()
		w.Write([]byte(renderer.RenderWorld(world)))
	}
}

func ajaxSelectGopher(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		position := r.FormValue("position")

		if world.SelectEntity(position) {
			w.Write([]byte(renderer.RenderWorld(world)))
		} else {
			w.WriteHeader(404)
		}

	}
}

func ajaxHandleWorldInput(world *gopherlife.World, renderer *gopherlife.Renderer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var leftArrow = "37"
		var rightArrow = "39"
		var upArrow = "38"
		var downArrow = "40"

		var pKey = "80"

		keydown := r.FormValue("keydown")

		switch keydown {
		case pKey:
			world.TogglePause()
		case leftArrow:
			renderer.StartX--
		case rightArrow:
			renderer.StartX++
		case upArrow:
			renderer.StartY--
		case downArrow:
			renderer.StartY++
		}

	}
}

type PageVariables struct {
	Data template.HTML
}
