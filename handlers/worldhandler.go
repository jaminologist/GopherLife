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
	data     WorldPageData
}

type SelectReturn struct {
	WorldRender string
	Gopher      *world.Gopher
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

	data := WorldPageData{
		PageTitle: "G O P H E R L I F E",
		FormData:  GetFormDataFromStats(stats),
		MapData:   []MapData{},
	}

	var tileMap = world.NewGopherMapWithSpiralSearch(stats)

	tileMapFunctions := make(map[string]func(world.Statistics) UpdateableRender)
	tileMapFunctions["GopherMap With Spiral Search"] = func(s world.Statistics) UpdateableRender {
		d := world.NewGopherMapWithSpiralSearch(s)
		return &d
	}
	tileMapFunctions["GopherMap With Partition"] = func(s world.Statistics) UpdateableRender {
		d := world.NewGopherMapWithParitionGridAndSearch(s)
		return &d
	}
	tileMapFunctions["Cool Black And White Spiral"] = func(s world.Statistics) UpdateableRender {
		d := world.NewSpiralMapController(s)
		return &d
	}

	for k := range tileMapFunctions {
		data.MapData = append(data.MapData, MapData{
			DisplayName: k,
			Value:       k,
		})
	}

	container := container{
		tileMap:  &tileMap,
		tileMaps: tileMapFunctions,
		data:     data,
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&container, data))
	http.HandleFunc("/ProcessWorld", ajaxProcessWorld(&container))
	http.HandleFunc("/ShiftWorldView", ajaxHandleWorldInput(&container))
	http.HandleFunc("/SelectGopher", ajaxSelectGopher(&container))
	http.HandleFunc("/ResetWorld", resetWorld(&container))
	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

func GetFormDataFromStats(stats world.Statistics) []FormData {

	return []FormData{
		FormData{
			DisplayName:        "Width",
			Type:               "Number",
			Name:               "width",
			Value:              strconv.Itoa(stats.Width),
			BootStrapFormWidth: 1,
		},
		FormData{
			DisplayName:        "Height",
			Type:               "Number",
			Name:               "height",
			Value:              strconv.Itoa(stats.Height),
			BootStrapFormWidth: 1,
		},
		FormData{
			DisplayName:        "Initial Population",
			Type:               "Number",
			Name:               "numberOfGophers",
			Value:              strconv.Itoa(stats.NumberOfGophers),
			BootStrapFormWidth: 2,
		},
		FormData{
			DisplayName:        "Max Population",
			Type:               "Number",
			Name:               "maxPopulation",
			Value:              strconv.Itoa(stats.MaximumNumberOfGophers),
			BootStrapFormWidth: 2,
		},
		FormData{
			DisplayName:        "Birth Rate",
			Type:               "Number",
			Name:               "birthRate",
			Value:              strconv.Itoa(stats.GopherBirthRate),
			BootStrapFormWidth: 2,
		},
		FormData{
			DisplayName:        "Food",
			Type:               "Number",
			Name:               "numberOfFood",
			Value:              strconv.Itoa(stats.NumberOfFood),
			BootStrapFormWidth: 2,
		},
	}

}

type WorldPageData struct {
	PageTitle string
	FormData  []FormData
	MapData   []MapData
}

type FormData struct {
	DisplayName        string
	Name               string
	Value              string
	Type               string
	BootStrapFormWidth int
}

type MapData struct {
	DisplayName string
	Value       string
}

func worldToHTML(container *container, data WorldPageData) func(w http.ResponseWriter, r *http.Request) {

	//	renderer := gopherlife.NewRenderer()

	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("static/index.html"))
		err := tmpl.Execute(w, container.data)

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
		container.data.FormData = GetFormDataFromStats(stats)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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

		if err == nil {
			container.tileMap.KeyPress(world.Keys(key))
		}
	}
}
