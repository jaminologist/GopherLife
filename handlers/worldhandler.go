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

type UpdateableRender interface {
	Update() bool
	json.Marshaler
	world.Controllable
}

type container struct {
	tileMap  UpdateableRender
	tileMaps map[string]func(world.Statistics) UpdateableRender
}

type SelectReturn struct {
	WorldRender string
	Gopher      *world.Gopher
}

type PageVariables struct {
	Data template.HTML
}

func SetUpPage() {

	stats := world.Statistics{
		Width:                  50,
		Height:                 50,
		NumberOfGophers:        5,
		NumberOfFood:           200,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        7,
	}

	/*stats := world.Statistics{
		Width:                  3000,
		Height:                 3000,
		NumberOfGophers:        5000,
		NumberOfFood:           1000000,
		MaximumNumberOfGophers: 1000000,
		GopherBirthRate:        7,
	}*/

	var tileMap = world.NewSpiralMapController(stats)

	tileMapFunctions := make(map[string]func(world.Statistics) UpdateableRender)
	tileMapFunctions["a"] = func(s world.Statistics) UpdateableRender {
		d := world.NewGopherMapWithSpiralSearch(s)
		return &d
	}
	tileMapFunctions["b"] = func(s world.Statistics) UpdateableRender {
		d := world.NewGopherMapWithParitionGridAndSearch(s)
		return &d
	}
	container := container{
		tileMap:  &tileMap,
		tileMaps: tileMapFunctions,
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
			//Data: template.HTML(container.tileMap.Render()),
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
			jsonData, _ := container.tileMap.MarshalJSON()
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

		mapSelection := r.FormValue("mapSelection")

		stats := world.Statistics{
			Width:                  int(width),
			Height:                 int(height),
			NumberOfGophers:        int(numberOfGophers),
			NumberOfFood:           int(numberOfFood),
			MaximumNumberOfGophers: int(maxPopulation),
			GopherBirthRate:        int(birthRate),
		}

		var tileMap UpdateableRender

		if tileMapFunc, ok := container.tileMaps[mapSelection]; ok {
			tileMap = tileMapFunc(stats)
		} else {
			adr := world.NewGopherMapWithSpiralSearch(stats)
			tileMap = &adr
		}

		container.tileMap = tileMap
	}
}

func ajaxSelectGopher(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		position := calc.StringToCoordinates(r.FormValue("position"))
		container.tileMap.Click(position.GetX(), position.GetY())
	}
}

func ajaxHandleWorldInput(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		keydown := r.FormValue("keydown")
		key, err := strconv.ParseInt(keydown, 10, 64)
		if err != nil {
			container.tileMap.KeyPress(world.Keys(key))
		}
	}
}
