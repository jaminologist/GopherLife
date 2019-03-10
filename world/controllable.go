package world

import (
	"encoding/json"
	"fmt"
	"gopherlife/colors"
	"image/color"
	"net/url"
	"strconv"
	"strings"
)

//Controller used to define user controls for a world
type Controller interface {
	Click(x int, y int)
	KeyPress(key Keys)
}

type WorldPageData struct {
	PageTitle   string
	FormData    []FormData
	IsGopherMap bool
}

type FormData struct {
	DisplayName        string
	Name               string
	Value              string
	Type               string
	BootStrapFormWidth int
}

//Keys used to denote correct number for keys on a keyboard when called e.which in js
type Keys int64

//List of keys and their corresponding numbers
const (
	LeftArrow  Keys = 37
	RightArrow Keys = 39
	UpArrow    Keys = 38
	DownArrow  Keys = 40

	PKey Keys = 80
	QKey Keys = 81
	WKey Keys = 87
)

type GopherMapController struct {
	*GopherMap
	*Renderer
	CreateNew func(Statistics) *GopherMap
}

var (
	maleGopherColor         = color.RGBA{90, 218, 255, 1}
	maleGopherSelectedColor = color.RGBA{245, 245, 245, 1}
	youngMaleGopherColor    = color.RGBA{167, 235, 255, 1}

	femaleGopherColor         = color.RGBA{255, 255, 0, 1}
	femaleGopherSelectedColor = color.RGBA{255, 155, 154, 1}
	youngfemaleGopherColor    = color.RGBA{255, 231, 231, 1}

	foodColor          = color.RGBA{204, 112, 0, 1}
	decayedGopherColor = color.RGBA{0, 0, 0, 1}
	grassColor         = color.RGBA{65, 119, 15, 1}
)

func NewGopherMapWithSpiralSearch(stats Statistics) GopherMapController {
	gMap := CreateWorldCustom(stats)
	renderer := NewRenderer(100, 100)
	return GopherMapController{
		GopherMap: gMap,
		Renderer:  &renderer,
		CreateNew: CreateWorldCustom,
	}
}

func NewGopherMapWithParitionGridAndSearch(stats Statistics) GopherMapController {
	gMap := CreatePartitionTileMapCustom(stats)
	renderer := NewRenderer(100, 100)
	return GopherMapController{
		GopherMap: gMap,
		Renderer:  &renderer,
		CreateNew: CreatePartitionTileMapCustom,
	}
}

func (controller *GopherMapController) Start() {
	if controller.GopherMap == nil {
		controller.GopherMap = controller.CreateNew(*controller.Statistics)
	}
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *GopherMapController) Click(x int, y int) {

	action := func() {
		_, ok := controller.SelectEntity(x, y)

		if !ok {
			controller.Renderer.StartX = x - controller.Renderer.RenderSizeX/2
			controller.Renderer.StartY = y - controller.Renderer.RenderSizeY/2
		}
	}

	controller.GopherMap.Add(action)
}

func (controller *GopherMapController) KeyPress(key Keys) {

	//fmt.Println(key)

	switch key {
	case WKey:
		controller.UnSelectGopher()
	case QKey:
		controller.SelectRandomGopher()
	case PKey:
		controller.TogglePause()
	case LeftArrow:
		controller.ShiftRenderer(-1, 0)
	case RightArrow:
		controller.ShiftRenderer(1, 0)
	case UpArrow:
		controller.ShiftRenderer(0, -1)
	case DownArrow:
		controller.ShiftRenderer(0, 1)
	}
}

func TileToColor(tile *Tile, isSelected bool) color.RGBA {

	switch {
	case tile.isEmpty():
		return grassColor
	case tile.Gopher != nil:

		switch tile.Gopher.Gender {
		case Male:
			if isSelected {
				return maleGopherSelectedColor
			} else if !tile.Gopher.IsMature() {
				return youngMaleGopherColor
			} else {
				return maleGopherColor
			}
		case Female:
			if isSelected {
				return femaleGopherSelectedColor
			} else if !tile.Gopher.IsMature() {
				return youngfemaleGopherColor
			} else {
				return femaleGopherColor
			}
		}

		if tile.Gopher.IsDead {
			return decayedGopherColor
		}

	case tile.Food != nil:
		return foodColor
	}

	return color.RGBA{0, 0, 0, 0}

}

func (controller *GopherMapController) RenderTile(x int, y int) color.RGBA {

	if tile, ok := controller.Tile(x, y); ok {

		switch {
		case tile.isEmpty():
			return grassColor
		case tile.Gopher != nil:
			isSelected := false
			if controller.GopherMap.SelectedGopher != nil {
				isSelected = controller.GopherMap.SelectedGopher.Position.GetX() == x &&
					controller.GopherMap.SelectedGopher.Position.GetY() == y
			}
			return TileToColor(tile, isSelected)
		case tile.Food != nil:
			return foodColor
		}
	} else {
		return color.RGBA{137, 207, 240, 1}
	}

	return color.RGBA{137, 207, 240, 1}

}

func (controller *GopherMapController) MarshalJSON() ([]byte, error) {

	if controller.SelectedGopher != nil {
		controller.Renderer.StartX = controller.SelectedGopher.Position.GetX() - controller.Renderer.RenderSizeX/2
		controller.Renderer.StartY = controller.SelectedGopher.Position.GetY() - controller.Renderer.RenderSizeY/2
	}

	render := controller.Renderer.RenderWorld(controller)

	if controller.SelectedGopher != nil {
		render.SelectedGopher = controller.SelectedGopher
	}

	stats := controller.Stats()
	diagnostics := controller.Diagnostics()

	renderString := ""
	renderString += "<br />"
	renderString += fmt.Sprintf("<span>Number of Gophers: %d </span><br />", stats.NumberOfGophers)
	renderString += fmt.Sprintf("<span>Avg Processing Time (s): %s </span><br />", diagnostics.processStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Avg Gopher Time (s): %s </span><br />", diagnostics.gopherStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span; >Avg Input Time (s): %s </span><br />", diagnostics.inputStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Total Elasped Time (s): %s </span><br />", diagnostics.globalStopWatch.GetCurrentElaspedTime().String())

	render.WorldRender = renderString

	return json.Marshal(render)
}

func (controller *GopherMapController) PageLayout() WorldPageData {

	stats := controller.Statistics

	formdataArray := []FormData{
		FormDataWidth(stats.Width, 2),
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

	return WorldPageData{
		PageTitle:   "G O P H E R L I F E <b>2.0</b>",
		FormData:    formdataArray,
		IsGopherMap: true,
	}
}

func (controller *GopherMapController) HandleForm(values url.Values) bool {

	if strings.Contains(values.Encode(), "numberOfGophers") {

		width, _ := strconv.ParseInt(values.Get("width"), 10, 64)
		height, _ := strconv.ParseInt(values.Get("height"), 10, 64)
		numberOfGophers, _ := strconv.ParseInt(values.Get("numberOfGophers"), 10, 64)
		numberOfFood, _ := strconv.ParseInt(values.Get("numberOfFood"), 10, 64)
		birthRate, _ := strconv.ParseInt(values.Get("birthRate"), 10, 64)
		maxPopulation, _ := strconv.ParseInt(values.Get("maxPopulation"), 10, 64)

		stats := Statistics{
			Width:                  int(width),
			Height:                 int(height),
			NumberOfGophers:        int(numberOfGophers),
			NumberOfFood:           int(numberOfFood),
			MaximumNumberOfGophers: int(maxPopulation),
			GopherBirthRate:        int(birthRate),
		}

		gmc := controller.CreateNew(stats)
		controller.GopherMap = gmc

	}

	return true
}

type NoPlayerInput struct{}

func (controller *NoPlayerInput) Click(x int, y int) {
}

func (controller *NoPlayerInput) KeyPress(key Keys) {
}

type SpiralMapController struct {
	NoPlayerInput
	Statistics
	*SpiralMap
	*Renderer
}

func NewSpiralMapController(stats Statistics) SpiralMapController {

	stats = Statistics{
		Width:                  50,
		Height:                 50,
		NumberOfGophers:        5,
		NumberOfFood:           200,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        7,
	}

	renderer := NewRenderer(75, 75)
	renderer.ShiftRenderer(stats.Width/2-renderer.RenderSizeX/2, stats.Height/2-renderer.RenderSizeY/2)

	return SpiralMapController{
		Renderer:   &renderer,
		Statistics: stats,
	}
}

func (controller *SpiralMapController) Start() {
	if controller.SpiralMap == nil {
		sMap := NewSpiralMap(controller.Statistics)
		controller.SpiralMap = &sMap
	}
}

func (controller *SpiralMapController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.Renderer.RenderWorld(controller))
}

func (controller *SpiralMapController) RenderTile(x int, y int) color.RGBA {
	if tile, ok := controller.Tile(x, y); ok {
		if tile.HasGopher() {
			return color.RGBA{255, 255, 255, 1}
		} else {
			return color.RGBA{0, 0, 0, 1}
		}
	} else {
		return color.RGBA{0, 0, 0, 1}
	}
}

func (controller *SpiralMapController) PageLayout() WorldPageData {
	return WorldPageData{}
}

func (controller *SpiralMapController) HandleForm(values url.Values) bool {
	sm := NewSpiralMapController(Statistics{}).SpiralMap
	controller.SpiralMap = sm
	return false
}

type FireWorksController struct {
	NoPlayerInput
	Statistics
	*GopherMap
	*Renderer
}

func NewFireWorksController(stats Statistics) FireWorksController {

	stats = Statistics{
		Width:                  400,
		Height:                 200,
		NumberOfGophers:        1000,
		NumberOfFood:           5000,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        25,
	}

	renderer := NewRenderer(400, 150)
	renderer.ShiftRenderer(stats.Width/2-renderer.RenderSizeX/2, stats.Height/2-renderer.RenderSizeY/2)

	return FireWorksController{
		Statistics: stats,
		Renderer:   &renderer,
	}
}

func (controller *FireWorksController) Start() {
	if controller.GopherMap == nil {
		controller.GopherMap = CreateWorldCustom(controller.Statistics)
	}
}

func (controller *FireWorksController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.Renderer.RenderWorld(controller))
}

func (controller *FireWorksController) RenderTile(x int, y int) color.RGBA {
	if tile, ok := controller.Tile(x, y); ok {
		return TileToColor(tile, false)
	} else {
		return color.RGBA{0, 0, 0, 1}
	}
}

func (controller *FireWorksController) PageLayout() WorldPageData {
	return WorldPageData{}
}

func (controller *FireWorksController) HandleForm(values url.Values) bool {
	sm := NewFireWorksController(Statistics{}).GopherMap
	controller.GopherMap = sm
	return false
}

type CollisionMapController struct {
	NoPlayerInput
	Statistics
	*CollisionMap
	*Renderer
	CreateNew  func(Statistics, bool) CollisionMap
	IsDiagonal bool
}

func NewCollisionMapController(stats Statistics) CollisionMapController {

	stats = Statistics{
		Width:                  75,
		Height:                 75,
		NumberOfGophers:        500,
		MaximumNumberOfGophers: 100000,
	}

	renderer := NewRenderer(100, 100)
	renderer.ShiftRenderer(stats.Width/2-renderer.RenderSizeX/2, stats.Height/2-renderer.RenderSizeY/2)

	return CollisionMapController{
		Statistics: stats,
		Renderer:   &renderer,
		CreateNew:  NewCollisionMap,
		IsDiagonal: false,
	}
}

func (controller *CollisionMapController) Start() {
	if controller.CollisionMap == nil {
		sMap := controller.CreateNew(controller.Statistics, controller.IsDiagonal)
		controller.CollisionMap = &sMap
	}
}

func (controller *CollisionMapController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.Renderer.RenderWorld(controller))
}

func (controller *CollisionMapController) RenderTile(x int, y int) color.RGBA {

	if controller.Contains(x, y) {
		if c, ok := controller.HasCollider(x, y); ok {
			return c.Color
		} else {
			return color.RGBA{0, 0, 0, 1}
		}
	}

	return color.RGBA{255, 255, 255, 1}
}

func (controller *CollisionMapController) PageLayout() WorldPageData {
	stats := controller.Statistics

	formdataArray := []FormData{
		FormDataWidth(stats.Width, 2),
		FormData{
			DisplayName:        "Height",
			Type:               "Number",
			Name:               "height",
			Value:              strconv.Itoa(stats.Height),
			BootStrapFormWidth: 2,
		},
		FormData{
			DisplayName:        "Initial Population",
			Type:               "Number",
			Name:               "numberOfGophers",
			Value:              strconv.Itoa(stats.NumberOfGophers),
			BootStrapFormWidth: 2,
		},
	}

	return WorldPageData{
		PageTitle: "C O L L I D E R L I F E",
		FormData:  formdataArray,
	}
}

func (controller *CollisionMapController) HandleForm(values url.Values) bool {

	if strings.Contains(values.Encode(), "numberOfGophers") {

		width, _ := strconv.ParseInt(values.Get("width"), 10, 64)
		height, _ := strconv.ParseInt(values.Get("height"), 10, 64)
		numberOfGophers, _ := strconv.ParseInt(values.Get("numberOfGophers"), 10, 64)

		stats := Statistics{
			Width:           int(width),
			Height:          int(height),
			NumberOfGophers: int(numberOfGophers),
		}

		gmc := NewCollisionMap(stats, controller.IsDiagonal)
		controller.CollisionMap = &gmc
		controller.Statistics = stats

	}
	return true
}

func NewDiagonalCollisionMapController(stats Statistics) CollisionMapController {

	stats = Statistics{
		Width:                  75,
		Height:                 75,
		NumberOfGophers:        500,
		MaximumNumberOfGophers: 100000,
	}

	renderer := NewRenderer(100, 100)

	renderer.ShiftRenderer(stats.Width/2-renderer.RenderSizeX/2, stats.Height/2-renderer.RenderSizeY/2)

	return CollisionMapController{
		Statistics: stats,
		Renderer:   &renderer,
		CreateNew:  NewCollisionMap,
		IsDiagonal: true,
	}
}

func FormDataWidth(width int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Width",
		Type:               "Number",
		Name:               "width",
		Value:              strconv.Itoa(width),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

func FormDataSnakeSlowDown(slowdown int, bootstrapColumnWidth int) FormData {
	return FormData{
		DisplayName:        "Snake SlowDown",
		Type:               "Number",
		Name:               "SnakeSlowDown",
		Value:              strconv.Itoa(slowdown),
		BootStrapFormWidth: bootstrapColumnWidth,
	}
}

type SnakeMapController struct {
	Dimensions
	FrameSpeed int
	*SnakeMap
	*Renderer
}

func NewSnakeMapController(stats Statistics) SnakeMapController {

	d := Dimensions{
		Width:  50,
		Height: 50,
	}

	renderer := NewRenderer(50, 50)
	renderer.ShiftRenderer(d.Width/2-renderer.RenderSizeX/2, d.Height/2-renderer.RenderSizeY/2)

	return SnakeMapController{
		Dimensions: d,
		Renderer:   &renderer,
	}
}

func (controller *SnakeMapController) Start() {
	if controller.SnakeMap == nil {
		sMap := NewSnakeMap(controller.Dimensions, controller.FrameSpeed)
		controller.SnakeMap = &sMap
	}
}

func (controller *SnakeMapController) MarshalJSON() ([]byte, error) {

	render := controller.Renderer.RenderWorld(controller)

	render.WorldRender += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	if controller.SnakeMap.IsGameOver {
		render.WorldRender += fmt.Sprintf("<span>Game Over!</span><br />")
	}

	return json.Marshal(render)
}

func (controller *SnakeMapController) RenderTile(x int, y int) color.RGBA {

	if sp, ok := controller.Tile(x, y); ok {
		switch {
		case sp.SnakePart != nil:
			return colors.NokiaBorder
		case sp.SnakeFood != nil:
			return colors.NokiaBorder
		case sp.Wall != nil:
			return colors.NokiaBorder
		default:
			return colors.NokiaGreen
		}
	}
	return colors.White
}

func (controller *SnakeMapController) PageLayout() WorldPageData {
	return WorldPageData{
		PageTitle: "S N A K E L I F E",
		FormData: []FormData{
			FormDataSnakeSlowDown(1, 3),
		},
	}
}

func (controller *SnakeMapController) HandleForm(values url.Values) bool {

	fd := FormDataSnakeSlowDown(0, 0)
	if strings.Contains(values.Encode(), fd.Name) {
		speed, _ := strconv.ParseInt(values.Get(fd.Name), 10, 64)
		sm := NewSnakeMap(controller.Dimensions, int(speed))
		controller.SnakeMap = &sm

	}

	return true
}

func (controller *SnakeMapController) KeyPress(key Keys) {
	switch key {
	case LeftArrow:
		controller.SnakeMap.ChangeDirection(Left)
	case RightArrow:
		controller.SnakeMap.ChangeDirection(Right)
	case UpArrow:
		controller.SnakeMap.ChangeDirection(Up)
	case DownArrow:
		controller.SnakeMap.ChangeDirection(Down)
	}
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *SnakeMapController) Click(x int, y int) {
}
