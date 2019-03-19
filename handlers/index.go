package handlers

import (
	"encoding/json"
	"fmt"
	"gopherlife/controllers"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
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

func (c *ControllerContainer) Add(rc RenderController, key string) {
	c.RenderControllers[key] = rc
}

func (c *ControllerContainer) AddSelected(rc RenderController, key string) {
	c.Add(rc, key)
	c.SelectedKey = key
}

func (c *ControllerContainer) Selected() RenderController {

	if c.SelectedKey == "" {
		for k := range c.RenderControllers {
			c.SelectedKey = k
			break
		}
	}

	return c.RenderControllers[c.SelectedKey]
}

func (c *ControllerContainer) PopulatePageData() {

	data := PageData{}
	data.WorldPageData = c.Selected().PageLayout()
	data.Selected = c.SelectedKey

	keys := make([]string, 0)
	for key := range c.RenderControllers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		data.WorldSelectFormInput = append(data.WorldSelectFormInput, WorldSelectFormInput{
			DisplayName: key,
			Value:       key,
		})
	}

	c.pageData = data

}

type PageData struct {
	controllers.WorldPageData
	WorldSelectFormInput []WorldSelectFormInput
	Selected             string
}

type WorldSelectFormInput struct {
	DisplayName string
	Value       string
}

func SetUpPage() {

	ControllerContainer := NewControllerContainer()

	ss := controllers.NewGopherWorldWithSpiralSearch()
	ControllerContainer.AddSelected(&ss, "GopherWorld With Spiral Search")

	ps := controllers.NewGopherWorldWithParitionGridAndSearch()
	ControllerContainer.Add(&ps, "GopherWorld With Partition")

	sm := controllers.NewSpiralWorldController()
	ControllerContainer.Add(&sm, "Black and White Spiral World")

	wsm := controllers.NewWeirdSpiralWorldController()
	ControllerContainer.Add(&wsm, "Black and White Spiral World (Weird)")

	fireworks := controllers.NewFireWorksController()
	ControllerContainer.Add(&fireworks, "Fireworks!")

	collision := controllers.NewCollisionWorldController()
	ControllerContainer.Add(&collision, "Collision World")

	diagonalCollision := controllers.NewDiagonalCollisionWorldController()
	ControllerContainer.Add(&diagonalCollision, "Collision World (Diagonal)")

	SnakeWorld := controllers.NewSnakeWorldController()
	ControllerContainer.Add(&SnakeWorld, "Elongating Gopher")

	blockblockRevolution := controllers.NewBlockBlockRevolutionController()
	ControllerContainer.Add(&blockblockRevolution, "Block Block Revolution")

	ControllerContainer.Selected().Start()
	ControllerContainer.PopulatePageData()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", worldToHTML(&ControllerContainer))
	http.HandleFunc("/Update", Update(&ControllerContainer))
	http.HandleFunc("/Click", HandleClick(&ControllerContainer))
	http.HandleFunc("/KeyPress", HandleKeyPress(&ControllerContainer))
	http.HandleFunc("/Scroll", HandleScroll(&ControllerContainer))
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

		worldSelection := r.FormValue("worldSelection")
		ControllerContainer.SelectedKey = worldSelection
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
			w.WriteHeader(200)
		}
	}

}

func HandleScroll(ControllerContainer *ControllerContainer) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		deltaY := r.FormValue("deltaY")
		deltaYNum, err := strconv.Atoi(deltaY)

		if err == nil {
			ControllerContainer.Selected().Scroll(deltaYNum)
		}
		w.WriteHeader(200)
	}
}
