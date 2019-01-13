package handler

import (
	"encoding/json"
	"fmt"
	"gopherlife/calc"
	"gopherlife/world"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type container struct {
	tileMap  world.TileMap
	renderer *world.Renderer
}

type SelectReturn struct {
	WorldRender string
	Gopher      *world.Gopher
}

type PageVariables struct {
	Data template.HTML
}

func SetUpPage() {

	var tileMap = world.CreatePartitionTileMap()
	renderer := world.NewRenderer()

	container := container{
		tileMap:  &tileMap,
		renderer: &renderer,
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&container))
	http.HandleFunc("/ProcessWorld", ajaxProcessWorld(&container))
	http.HandleFunc("/ShiftWorldView", ajaxHandleWorldInput(&container))
	http.HandleFunc("/SelectGopher", ajaxSelectGopher(&container))
	http.HandleFunc("/ResetWorld", resetWorld(&container))

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

func worldToHTML(container *container) func(w http.ResponseWriter, r *http.Request) {

	//	renderer := gopherlife.NewRenderer()

	return func(w http.ResponseWriter, r *http.Request) {

		pageVariables := PageVariables{
			Data: template.HTML(container.renderer.RenderWorld(container.tileMap).WorldRender),
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

func ajaxProcessWorld(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		container.tileMap.Update()

		if true {
			jsonData, _ := json.Marshal(container.renderer.RenderWorld(container.tileMap))
			w.Write(jsonData)
		} else {
			w.WriteHeader(404)
		}

	}
}

func resetWorld(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		width, _ := strconv.ParseInt(r.FormValue("width"), 10, 64)
		height, _ := strconv.ParseInt(r.FormValue("height"), 10, 64)
		numberOfGophers, _ := strconv.ParseInt(r.FormValue("numberOfGophers"), 10, 64)
		numberOfFood, _ := strconv.ParseInt(r.FormValue("numberOfFood"), 10, 64)
		birthRate, _ := strconv.ParseInt(r.FormValue("birthRate"), 10, 64)
		maxPopulation, _ := strconv.ParseInt(r.FormValue("maxPopulation"), 10, 64)

		tileMap := world.CreateWorldCustom(
			world.Statistics{
				Width:                  int(width),
				Height:                 int(height),
				NumberOfGophers:        int(numberOfGophers),
				NumberOfFood:           int(numberOfFood),
				MaximumNumberOfGophers: int(maxPopulation),
				GopherBirthRate:        int(birthRate),
			},
		)

		container.tileMap = &tileMap
	}
}

func ajaxSelectGopher(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		position := calc.StringToCoordinates(r.FormValue("position"))

		if _, ok := container.tileMap.SelectEntity(position.GetX(), position.GetY()); ok {
			w.Header().Set("Content-Type", "application/json")
			jsonData, _ := json.Marshal(container.renderer.RenderWorld(container.tileMap))

			w.Write(jsonData)
		} else {
			w.WriteHeader(404)
		}

	}
}

func ajaxHandleWorldInput(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var leftArrow = "37"
		var rightArrow = "39"
		var upArrow = "38"
		var downArrow = "40"

		var pKey = "80"

		var qKey = "81"
		var wKey = "87"

		keydown := r.FormValue("keydown")

		switch keydown {
		case wKey:
			container.tileMap.UnSelectGopher()
		case qKey:
			container.tileMap.SelectRandomGopher()
		case pKey:
			container.tileMap.TogglePause()
		case leftArrow:
			container.renderer.ShiftRenderer(-1, 0)
		case rightArrow:
			container.renderer.ShiftRenderer(1, 0)
		case upArrow:
			container.renderer.ShiftRenderer(0, -1)
		case downArrow:
			container.renderer.ShiftRenderer(0, 1)
		}

	}
}
