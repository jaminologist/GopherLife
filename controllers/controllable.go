package controllers

import (
	"encoding/json"
	"fmt"
	"gopherlife/colors"
	"gopherlife/geometry"
	"gopherlife/renderers"
	"gopherlife/world"
	"image/color"
	"net/url"
	"strconv"
	"strings"
)

//UserInputHandler used to handle user inputs for a world
type UserInputHandler interface {
	Click(x int, y int)
	KeyPress(key Keys)
}

type WorldPageData struct {
	PageTitle   string
	FormData    []FormData
	IsGopherMap bool
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
	*world.GopherMap
	*renderers.GridRenderer
	CreateNew func(world.GopherMapSettings) *world.GopherMap
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

//NewGopherMapWithSpiralSearch Returns a Controller with a Gopher Map. Where Gophers search for food using a Spiral To Nearest Search
func NewGopherMapWithSpiralSearch() GopherMapController {

	settings := world.GopherMapSettings{
		Dimensions:      world.Dimensions{Width: 3000, Height: 3000},
		Population:      world.Population{InitialPopulation: 5000, MaxPopulation: 1000000},
		NumberOfFood:    1000000,
		GopherBirthRate: 7,
	}

	gMap := world.CreateWorldCustom(settings)
	renderer := renderers.NewRenderer(100, 100)
	return GopherMapController{
		GopherMap:    gMap,
		GridRenderer: &renderer,
		CreateNew:    world.CreateWorldCustom,
	}
}

//NewGopherMapWithParitionGridAndSearch Returns a Controller with a Gopher Map. Where Gophers search for food using Grid Partition
func NewGopherMapWithParitionGridAndSearch() GopherMapController {

	settings := world.GopherMapSettings{
		Dimensions:      world.Dimensions{Width: 3000, Height: 3000},
		Population:      world.Population{InitialPopulation: 5000, MaxPopulation: 1000000},
		NumberOfFood:    1000000,
		GopherBirthRate: 7,
	}

	gMap := world.CreatePartitionTileMapCustom(settings)
	renderer := renderers.NewRenderer(100, 100)
	return GopherMapController{
		GopherMap:    gMap,
		GridRenderer: &renderer,
		CreateNew:    world.CreatePartitionTileMapCustom,
	}
}

//Start Initiates the controller. If the Map does not exist. The Map will be built
func (controller *GopherMapController) Start() {
	if controller.GopherMap == nil {
		controller.GopherMap = controller.CreateNew(*controller.GopherMapSettings)
	}
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *GopherMapController) Click(x int, y int) {

	action := func() {
		_, ok := controller.SelectEntity(x, y)

		if !ok {
			controller.GridRenderer.StartX = x - controller.GridRenderer.Width/2
			controller.GridRenderer.StartY = y - controller.GridRenderer.Height/2
		}
	}

	controller.GopherMap.Add(action)
}

func (controller *GopherMapController) KeyPress(key Keys) {

	switch key {
	case WKey:
		controller.UnSelectGopher()
	case QKey:
		controller.SelectRandomGopher()
	case PKey:
		controller.TogglePause()
	case LeftArrow:
		controller.Shift(-1, 0)
	case RightArrow:
		controller.Shift(1, 0)
	case UpArrow:
		controller.Shift(0, -1)
	case DownArrow:
		controller.Shift(0, 1)
	}
}

func TileToColor(tile *world.Tile, isSelected bool) color.RGBA {

	switch {
	case tile.IsEmpty():
		return grassColor
	case tile.Gopher != nil:

		switch tile.Gopher.Gender {
		case world.Male:
			if isSelected {
				return maleGopherSelectedColor
			} else if !tile.Gopher.IsMature() {
				return youngMaleGopherColor
			} else {
				return maleGopherColor
			}
		case world.Female:
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
		case tile.IsEmpty():
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

type GopherMapRender struct {
	SelectedGopher *world.Gopher
	renderers.Render
}

func (controller *GopherMapController) MarshalJSON() ([]byte, error) {

	if controller.SelectedGopher != nil {
		controller.GridRenderer.StartX = controller.SelectedGopher.Position.GetX() - controller.GridRenderer.Width/2
		controller.GridRenderer.StartY = controller.SelectedGopher.Position.GetY() - controller.GridRenderer.Height/2
	}

	render := controller.GridRenderer.Draw(controller)

	diagnostics := controller.Diagnostics()

	renderString := ""
	renderString += "<br />"
	renderString += fmt.Sprintf("<span>Number of Gophers: %d </span><br />", controller.NumberOfGophers)
	renderString += fmt.Sprintf("<span>Avg Processing Time (s): %s </span><br />", diagnostics.ProcessStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Avg Gopher Time (s): %s </span><br />", diagnostics.GopherStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span; >Avg Input Time (s): %s </span><br />", diagnostics.InputStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Total Elasped Time (s): %s </span><br />", diagnostics.GlobalStopWatch.GetCurrentElaspedTime().String())

	render.WorldRender = renderString

	gmr := GopherMapRender{
		Render: render,
	}
	if controller.SelectedGopher != nil {
		gmr.SelectedGopher = controller.SelectedGopher
	}

	return json.Marshal(gmr)
}

func (controller *GopherMapController) PageLayout() WorldPageData {

	settings := controller.GopherMapSettings

	formdataArray := []FormData{
		FormDataWidth(settings.Width, 2),
		FormDataHeight(settings.Height, 2),
		FormDataInitialPopulation(settings.InitialPopulation, 2),
		FormDataMaxPopulation(settings.MaxPopulation, 2),
		FormData{
			DisplayName:        "Birth Rate",
			Type:               "Number",
			Name:               "birthRate",
			Value:              strconv.Itoa(settings.GopherBirthRate),
			BootStrapFormWidth: 2,
		},
		FormData{
			DisplayName:        "Food",
			Type:               "Number",
			Name:               "numberOfFood",
			Value:              strconv.Itoa(settings.NumberOfFood),
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

	if strings.Contains(values.Encode(), "birthRate") {

		width, _ := strconv.ParseInt(values.Get("width"), 10, 64)
		height, _ := strconv.ParseInt(values.Get("height"), 10, 64)
		InitialPopulation, _ := strconv.ParseInt(values.Get(FormDataInitialPopulation(0, 0).Name), 10, 64)
		numberOfFood, _ := strconv.ParseInt(values.Get("numberOfFood"), 10, 64)
		birthRate, _ := strconv.ParseInt(values.Get("birthRate"), 10, 64)
		maxPopulation, _ := strconv.ParseInt(values.Get("maxPopulation"), 10, 64)

		settings := world.GopherMapSettings{
			Dimensions:      world.Dimensions{Width: int(width), Height: int(height)},
			Population:      world.Population{InitialPopulation: int(InitialPopulation), MaxPopulation: int(maxPopulation)},
			NumberOfFood:    int(numberOfFood),
			GopherBirthRate: int(birthRate),
		}

		gmc := controller.CreateNew(settings)
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
	world.Statistics
	*world.SpiralMap
	*renderers.GridRenderer
}

func NewSpiralMapController(stats world.Statistics) SpiralMapController {

	stats = world.Statistics{
		Width:                  50,
		Height:                 50,
		NumberOfGophers:        5,
		NumberOfFood:           200,
		MaximumNumberOfGophers: 100000,
		GopherBirthRate:        7,
	}

	renderer := renderers.NewRenderer(75, 75)
	renderer.Shift(stats.Width/2-renderer.Width/2, stats.Height/2-renderer.Height/2)

	return SpiralMapController{
		GridRenderer: &renderer,
		Statistics:   stats,
	}
}

func (controller *SpiralMapController) Start() {
	if controller.SpiralMap == nil {
		sMap := world.NewSpiralMap(controller.Statistics)
		controller.SpiralMap = &sMap
	}
}

func (controller *SpiralMapController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.GridRenderer.Draw(controller))
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
	sm := NewSpiralMapController(world.Statistics{}).SpiralMap
	controller.SpiralMap = sm
	return false
}

type FireWorksController struct {
	NoPlayerInput
	world.GopherMapSettings
	*world.GopherMap
	*renderers.GridRenderer
}

func NewFireWorksController(stats world.Statistics) FireWorksController {

	settings := world.GopherMapSettings{
		Dimensions:      world.Dimensions{Width: 400, Height: 200},
		Population:      world.Population{InitialPopulation: 2000, MaxPopulation: 100000},
		NumberOfFood:    1000,
		GopherBirthRate: 35,
	}

	renderer := renderers.NewRenderer(400, 150)
	renderer.Shift(stats.Width/2-renderer.Width/2, stats.Height/2-renderer.Height/2)

	return FireWorksController{
		GopherMapSettings: settings,
		GridRenderer:      &renderer,
	}
}

func (controller *FireWorksController) Start() {
	if controller.GopherMap == nil {
		controller.GopherMap = world.CreateWorldCustom(controller.GopherMapSettings)
	}
}

func (controller *FireWorksController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.GridRenderer.Draw(controller))
}

func (controller *FireWorksController) RenderTile(x int, y int) color.RGBA {
	if tile, ok := controller.Tile(x, y); ok {

		if tile.HasGopher() {

			if !tile.Gopher.IsMature() {
				if tile.Gopher.Gender == world.Male {
					return colors.Cyan
				} else {
					return colors.Orange
				}
			}

			return colors.Black

		} else if tile.HasFood() {
			return colors.White
		}

		return color.RGBA{0, 0, 0, 1}
	} else {
		return color.RGBA{0, 0, 0, 1}
	}
}

func (controller *FireWorksController) PageLayout() WorldPageData {
	return WorldPageData{}
}

func (controller *FireWorksController) HandleForm(values url.Values) bool {
	sm := NewFireWorksController(world.Statistics{}).GopherMap
	controller.GopherMap = sm
	return false
}

type CollisionMapController struct {
	NoPlayerInput
	world.Statistics
	*world.CollisionMap
	*renderers.GridRenderer
	CreateNew  func(world.Statistics, bool) world.CollisionMap
	IsDiagonal bool
}

func NewCollisionMapController(stats world.Statistics) CollisionMapController {

	stats = world.Statistics{
		Width:                  75,
		Height:                 75,
		NumberOfGophers:        500,
		MaximumNumberOfGophers: 100000,
	}

	renderer := renderers.NewRenderer(100, 100)
	renderer.Shift(stats.Width/2-renderer.Width/2, stats.Height/2-renderer.Height/2)

	return CollisionMapController{
		Statistics:   stats,
		GridRenderer: &renderer,
		CreateNew:    world.NewCollisionMap,
		IsDiagonal:   false,
	}
}

func (controller *CollisionMapController) Start() {
	if controller.CollisionMap == nil {
		sMap := controller.CreateNew(controller.Statistics, controller.IsDiagonal)
		controller.CollisionMap = &sMap
	}
}

func (controller *CollisionMapController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.GridRenderer.Draw(controller))
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
		FormDataHeight(stats.Height, 2),
		FormDataInitialPopulation(stats.NumberOfGophers, 2),
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

		stats := world.Statistics{
			Width:           int(width),
			Height:          int(height),
			NumberOfGophers: int(numberOfGophers),
		}

		gmc := world.NewCollisionMap(stats, controller.IsDiagonal)
		controller.CollisionMap = &gmc
		controller.Statistics = stats

	}
	return true
}

func NewDiagonalCollisionMapController(stats world.Statistics) CollisionMapController {

	stats = world.Statistics{
		Width:                  75,
		Height:                 75,
		NumberOfGophers:        500,
		MaximumNumberOfGophers: 100000,
	}

	renderer := renderers.NewRenderer(100, 100)

	renderer.Shift(stats.Width/2-renderer.Width/2, stats.Height/2-renderer.Height/2)

	return CollisionMapController{
		Statistics:   stats,
		GridRenderer: &renderer,
		CreateNew:    world.NewCollisionMap,
		IsDiagonal:   true,
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
	world.Dimensions
	FrameSpeed   int
	ClickToBegin bool
	*world.SnakeMap
	*renderers.GridRenderer
}

func NewSnakeMapController(stats world.Statistics) SnakeMapController {

	d := world.Dimensions{
		Width:  35,
		Height: 35,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(d.Width/2-renderer.Width/2, d.Height/2-renderer.Height/2)

	return SnakeMapController{
		Dimensions:   d,
		GridRenderer: &renderer,
		FrameSpeed:   10,
	}
}

func (controller *SnakeMapController) Start() {
	if controller.SnakeMap == nil {
		sMap := world.NewSnakeMap(controller.Dimensions, controller.FrameSpeed)
		controller.SnakeMap = &sMap
	}
}

func (controller *SnakeMapController) MarshalJSON() ([]byte, error) {

	render := controller.GridRenderer.Draw(controller)
	render.WorldRender += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	if controller.SnakeMap.IsGameOver {
		render.WorldRender += fmt.Sprintf("<span>Game Over!</span><br />")
	} else if !controller.ClickToBegin {
		render.WorldRender += fmt.Sprintf("<span>Click to Begin")
	}

	return json.Marshal(render)
}

func (controller *SnakeMapController) RenderTile(x int, y int) color.RGBA {

	if sp, ok := controller.Tile(x, y); ok {
		switch {
		case sp.SnakePart != nil:
			if sp.SnakePart.HasPartInStomach() {
				return colors.NokiaFoodGreen
			} else {
				return colors.NokiaBorder
			}
		case sp.SnakeFood != nil:
			return colors.NokiaBorder
		case sp.SnakeWall != nil:
			return colors.NokiaBorder
		default:
			return colors.NokiaGreen
		}
	}
	return colors.White
}

func (controller *SnakeMapController) PageLayout() WorldPageData {
	return WorldPageData{
		PageTitle: "E L O N G A T I N G G O P H E R L I F E",
		FormData: []FormData{
			FormDataSnakeSlowDown(controller.FrameSpeed, 3),
		},
	}
}

func (controller *SnakeMapController) HandleForm(values url.Values) bool {

	fd := FormDataSnakeSlowDown(0, 0)
	if strings.Contains(values.Encode(), fd.Name) {
		speed, _ := strconv.ParseInt(values.Get(fd.Name), 10, 64)
		sm := world.NewSnakeMap(controller.Dimensions, int(speed))
		controller.FrameSpeed = int(speed)
		controller.SnakeMap = &sm

	}

	return true
}

func (controller *SnakeMapController) KeyPress(key Keys) {
	switch key {
	case LeftArrow:
		controller.SnakeMap.ChangeDirection(geometry.Left)
	case RightArrow:
		controller.SnakeMap.ChangeDirection(geometry.Right)
	case UpArrow:
		controller.SnakeMap.ChangeDirection(geometry.Up)
	case DownArrow:
		controller.SnakeMap.ChangeDirection(geometry.Down)
	}
}

func (controller *SnakeMapController) Update() bool {
	if controller.ClickToBegin {
		if controller.IsGameOver {
			controller.ClickToBegin = false
		}
		return controller.SnakeMap.Update()
	}

	return true
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *SnakeMapController) Click(x int, y int) {
	controller.ClickToBegin = true
}

type BlockBlockRevolutionController struct {
	world.Dimensions
	FrameSpeed   int
	ClickToBegin bool
	*world.BlockBlockRevolutionMap
	*renderers.GridRenderer
}

func NewBlockBlockRevolutionController() BlockBlockRevolutionController {

	d := world.Dimensions{
		Width:  10,
		Height: 20,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(d.Width/2-renderer.Width/2, d.Height/2-renderer.Height/2)

	return BlockBlockRevolutionController{
		Dimensions:   d,
		GridRenderer: &renderer,
		FrameSpeed:   10,
	}

}

func (controller *BlockBlockRevolutionController) Start() {
	if controller.BlockBlockRevolutionMap == nil {
		sMap := world.NewBlockBlockRevolutionMap(controller.Dimensions, controller.FrameSpeed)
		controller.BlockBlockRevolutionMap = &sMap
	}
}

func (controller *BlockBlockRevolutionController) MarshalJSON() ([]byte, error) {

	render := controller.GridRenderer.Draw(controller)
	//render.WorldRender += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	/*if controller.SnakeMap.IsGameOver {
		render.WorldRender += fmt.Sprintf("<span>Game Over!</span><br />")
	} else if !controller.ClickToBegin {
		render.WorldRender += fmt.Sprintf("<span>Click to Begin")
	} */

	return json.Marshal(render)
}

func (controller *BlockBlockRevolutionController) RenderTile(x int, y int) color.RGBA {

	if tile, ok := controller.Tile(x, y); ok {
		switch {
		case tile.Block != nil:
			return tile.Block.Color
		default:
			return colors.Black
		}
	}
	return colors.White
}

func (controller *BlockBlockRevolutionController) PageLayout() WorldPageData {
	return WorldPageData{
		PageTitle: "B L O C K B L O C K R E V O L U T I O N",
		FormData: []FormData{
			FormDataSnakeSlowDown(controller.FrameSpeed, 3),
		},
	}
}

func (controller *BlockBlockRevolutionController) HandleForm(values url.Values) bool {

	fd := FormDataSnakeSlowDown(0, 0)
	if strings.Contains(values.Encode(), fd.Name) {
		speed, _ := strconv.ParseInt(values.Get(fd.Name), 10, 64)
		bbrm := world.NewBlockBlockRevolutionMap(controller.Dimensions, int(speed))
		controller.FrameSpeed = int(speed)
		controller.BlockBlockRevolutionMap = &bbrm

	}

	return true
}

func (controller *BlockBlockRevolutionController) KeyPress(key Keys) {
	switch key {
	case LeftArrow:
		controller.Add(func() {
			controller.BlockBlockRevolutionMap.MoveCurrentTetrominoLeft()
		})
	case RightArrow:
		controller.Add(func() {
			controller.BlockBlockRevolutionMap.MoveCurrentTetrominoRight()
		})
	case UpArrow:
		controller.BlockBlockRevolutionMap.RotateTetromino()
	case DownArrow:
		controller.Add(func() {
			controller.BlockBlockRevolutionMap.InstantDown()
		})
	}
}

func (controller *BlockBlockRevolutionController) Update() bool {
	/*if controller.ClickToBegin {
		if controller.IsGameOver {
			controller.ClickToBegin = false
		}
		return controller.SnakeMap.Update()
	}*/

	return controller.BlockBlockRevolutionMap.Update()
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *BlockBlockRevolutionController) Click(x int, y int) {
	controller.ClickToBegin = true
}
