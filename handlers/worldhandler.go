package handler

import (
	"encoding/json"
	"fmt"
	"gopherlife/world"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type UpdateableRender interface {
	Update() bool
	Start()
	PageLayout() world.WorldPageData
	HandleForm(url.Values) bool
	json.Marshaler
	world.Controller
}

type container struct {
	tileMap  UpdateableRender
	tileMaps map[string]UpdateableRender
	pageData PageData
}

type PageData struct {
	world.WorldPageData
	MapData  []MapData
	Selected string
}

type MapData struct {
	DisplayName string
	Value       string
}

func SetUpPage() {

	/*stats := world.Statistics{
		Width:                  3000,
		Height:                 3000,
		NumberOfGophers:        5000,
		NumberOfFood:           1000000,
		MaximumNumberOfGophers: 1000000,
		GopherBirthRate:        7,
	}*/

	stats := world.Statistics{
		Width:                  100,
		Height:                 100,
		NumberOfGophers:        250,
		NumberOfFood:           1000,
		MaximumNumberOfGophers: 10000,
		GopherBirthRate:        7,
	}

	tileMapFunctions := make(map[string]UpdateableRender)
	ss := world.NewGopherMapWithSpiralSearch(stats)
	tileMapFunctions["GopherMap With Spiral Search"] = &ss
	ps := world.NewGopherMapWithParitionGridAndSearch(stats)
	tileMapFunctions["GopherMap With Partition"] = &ps
	cbws := world.NewSpiralMapController(stats)
	tileMapFunctions["Cool Black And White Spiral"] = &cbws
	fireworks := world.NewFireWorksController(stats)
	tileMapFunctions["Fireworks!"] = &fireworks
	collision := world.NewCollisionMapController(stats)
	tileMapFunctions["Collision Map"] = &collision
	diagonalCollision := world.NewDiagonalCollisionMapController(stats)
	tileMapFunctions["Diagonal Collision Map"] = &diagonalCollision
	snakeMap := world.NewSnakeMapController(stats)
	tileMapFunctions["Snake!!!!"] = &snakeMap
	blockblockRevolution := world.NewBlockBlockRevolutionController()
	tileMapFunctions["blockblockRevolution"] = &blockblockRevolution

	var selected = "blockblockRevolution"
	var tileMap = tileMapFunctions[selected]

	tileMap.Start()
	data := PageData{}
	data.WorldPageData = tileMap.PageLayout()
	data.Selected = selected

	for k := range tileMapFunctions {
		data.MapData = append(data.MapData, MapData{
			DisplayName: k,
			Value:       k,
		})
	}

	container := container{
		tileMap:  tileMap,
		tileMaps: tileMapFunctions,
		pageData: data,
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", worldToHTML(&container))
	http.HandleFunc("/ProcessWorld", ajaxProcessWorld(&container))
	http.HandleFunc("/ShiftWorldView", HandleKeyPress(&container))
	http.HandleFunc("/Click", HandleClick(&container))
	http.HandleFunc("/ResetWorld", ResetWorld(&container))
	http.HandleFunc("/SwitchWorld", SwitchWorld(&container))
	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

func worldToHTML(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("static/index.html"))
		err := tmpl.Execute(w, container.pageData)

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

func ResetWorld(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var tileMap UpdateableRender

		if tileMapFunc, ok := container.tileMaps[container.pageData.Selected]; ok {
			tileMapFunc.HandleForm(r.Form)
			tileMap = tileMapFunc
		} else {
			adr := world.NewGopherMapWithSpiralSearch(world.Statistics{})
			tileMap = &adr
		}

		container.tileMap = tileMap
		container.pageData.WorldPageData = container.tileMap.PageLayout()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func SwitchWorld(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		var stats world.Statistics
		mapSelection := r.FormValue("mapSelection")

		var tileMap UpdateableRender

		if tileMapFunc, ok := container.tileMaps[mapSelection]; ok {
			tileMap = tileMapFunc
			container.pageData.Selected = mapSelection
		} else {
			adr := world.NewGopherMapWithSpiralSearch(stats)
			tileMap = &adr
		}

		tileMap.Start()
		container.tileMap = tileMap
		container.pageData.WorldPageData = container.tileMap.PageLayout()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleClick(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		x, y := r.FormValue("x"), r.FormValue("y")

		xNum, _ := strconv.Atoi(x)
		yNum, _ := strconv.Atoi(y)

		container.tileMap.Click(xNum, yNum)
		w.WriteHeader(200)
	}
}

func HandleKeyPress(container *container) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		keydown := r.FormValue("keydown")
		key, err := strconv.ParseInt(keydown, 10, 64)

		if err == nil {
			container.tileMap.KeyPress(world.Keys(key))
		}
		w.WriteHeader(200)
	}

}
