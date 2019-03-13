package handler

import (
	"encoding/json"
	"fmt"
	"gopherlife/controllers"
	"gopherlife/world"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type RenderController interface {
	Update() bool
	Start()
	PageLayout() controllers.WorldPageData
	HandleForm(url.Values) bool
	json.Marshaler
	controllers.UserInputHandler
}

type ControllerContainer struct {
	SelectedKey string

	RenderControllers map[string]RenderController
	pageData          PageData
}

func NewControllerContainer() ControllerContainer {

	return ControllerContainer{
		RenderControllers: make(map[string]RenderController),
	}
}

func (c *ControllerContainer) Add(rc RenderController, key string, selected bool) {
	c.RenderControllers[key] = rc

	if selected {
		c.SelectedKey = key
	}
}

func (c *ControllerContainer) Selected() RenderController {
	return c.RenderControllers[c.SelectedKey]
}

func (c *ControllerContainer) PopulatePageData() {

	data := PageData{}
	data.WorldPageData = c.Selected().PageLayout()
	data.Selected = c.SelectedKey

	for k := range c.RenderControllers {
		data.MapData = append(data.MapData, MapData{
			DisplayName: k,
			Value:       k,
		})
	}

	c.pageData = data

}

type PageData struct {
	controllers.WorldPageData
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

	ControllerContainer := NewControllerContainer()

	ss := controllers.NewGopherMapWithSpiralSearch(stats)
	ControllerContainer.Add(&ss, "GopherMap With Spiral Search", false)

	ps := controllers.NewGopherMapWithParitionGridAndSearch(stats)
	ControllerContainer.Add(&ps, "GopherMap With Partition", false)

	cbws := controllers.NewSpiralMapController(stats)
	ControllerContainer.Add(&cbws, "Cool Black And White Spiral", false)

	fireworks := controllers.NewFireWorksController(stats)
	ControllerContainer.Add(&fireworks, "Fireworks!", false)

	collision := controllers.NewCollisionMapController(stats)
	ControllerContainer.Add(&collision, "Collision Map", false)

	diagonalCollision := controllers.NewDiagonalCollisionMapController(stats)
	ControllerContainer.Add(&diagonalCollision, "Diagonal Collision Map", false)

	snakeMap := controllers.NewSnakeMapController(stats)
	ControllerContainer.Add(&snakeMap, "Elongateing Gopher", false)

	blockblockRevolution := controllers.NewBlockBlockRevolutionController()
	ControllerContainer.Add(&blockblockRevolution, "Block Block Revolution", true)

	ControllerContainer.Selected().Start()
	ControllerContainer.PopulatePageData()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", worldToHTML(&ControllerContainer))
	http.HandleFunc("/Update", Update(&ControllerContainer))
	http.HandleFunc("/ShiftWorldView", HandleKeyPress(&ControllerContainer))
	http.HandleFunc("/Click", HandleClick(&ControllerContainer))
	http.HandleFunc("/ResetWorld", ResetWorld(&ControllerContainer))
	http.HandleFunc("/SwitchWorld", SwitchWorld(&ControllerContainer))
	fmt.Println("Listening...")
	http.ListenAndServe(":8080", nil)

}

func worldToHTML(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("static/index.html"))
		err := tmpl.Execute(w, ControllerContainer.pageData)

		if err != nil {
			log.Printf("Template executing error: ", err)
		}

	}

}

//Update Runs the Update function of the selected RenderController. Returns JSON
func Update(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		ControllerContainer.Selected().Update()
		jsonData, err := ControllerContainer.Selected().MarshalJSON()

		if err == nil {
			w.Write(jsonData)
		} else {
			//w.Write(err.Error())
			w.WriteHeader(500)
		}

	}
}

func ResetWorld(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ControllerContainer.Selected().HandleForm(r.Form)
		ControllerContainer.pageData.WorldPageData = ControllerContainer.Selected().PageLayout()
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func SwitchWorld(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		mapSelection := r.FormValue("mapSelection")
		ControllerContainer.SelectedKey = mapSelection
		ControllerContainer.Selected().Start()

		ControllerContainer.pageData.WorldPageData = ControllerContainer.Selected().PageLayout()
		ControllerContainer.pageData.Selected = ControllerContainer.SelectedKey
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func HandleClick(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		x, y := r.FormValue("x"), r.FormValue("y")

		xNum, _ := strconv.Atoi(x)
		yNum, _ := strconv.Atoi(y)

		ControllerContainer.Selected().Click(xNum, yNum)
		w.WriteHeader(200)
	}
}

func HandleKeyPress(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		keydown := r.FormValue("keydown")
		key, err := strconv.ParseInt(keydown, 10, 64)

		if err == nil {
			ControllerContainer.Selected().KeyPress(controllers.Keys(key))
		}
		w.WriteHeader(200)
	}

}
