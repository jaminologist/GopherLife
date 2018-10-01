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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&world))
	http.HandleFunc("/PollWorld", ajaxProcessWorld(&world))
	http.HandleFunc("/ShiftWorldView", ajaxHandleWorldInput(&world))

	fmt.Println("Listening")
	http.ListenAndServe(":8080", nil)

	//world.ProcessWorld()

	fmt.Println(time.Since(start))
	fmt.Println("The World Lasted for ", world.Moments, " Moments")
	fmt.Println("Done")
}

func worldToHTML(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		pageVariables := PageVariables{
			Data: template.HTML(world.RenderWorld()),
		}

		t, err := template.ParseFiles("static/index.html")

		if err != nil {
			log.Print("Template parsing error: ", err)
		}

		err = t.Execute(w, pageVariables)

		if err != nil {
			log.Printf("Template executing error: ", err)
		}

		//var renderString = world.RenderWorld()
		//w.Header().Add("Content-Type", "text/html")
		//w.Header().Add("Content-Length", strconv.Itoa(len(renderString)))

		//fmt.Fprintln(w, world.RenderWorld())
	}

}

func ajaxProcessWorld(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		world.ProcessWorld()

		w.Write([]byte(world.RenderWorld()))
	}
}

func ajaxHandleWorldInput(world *gopherlife.World) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var leftArrow = "37"
		var rightArrow = "39"
		var upArrow = "38"
		var downArrow = "40"

		keydown := r.FormValue("keydown")

		switch keydown {
		case leftArrow:
			world.StartX--
		case rightArrow:
			world.StartX++
		case upArrow:
			world.StartY--
		case downArrow:
			world.StartY++
		}

	}
}

type PageVariables struct {
	Data template.HTML
}
