package controllers

import (
	"encoding/json"
	"fmt"
	"gopherlife/colors"
	"gopherlife/renderers"
	"gopherlife/world"
	"image/color"
	"net/url"
	"strconv"
	"strings"
)

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

type GopherMapController struct {
	*world.GopherWorld
	*renderers.GridRenderer
	CreateNew func(world.GopherMapSettings) *world.GopherWorld
}

//NewGopherMapWithSpiralSearch Returns a Controller with a Gopher Map. Where Gophers search for food using a Spiral To Nearest Search
func NewGopherMapWithSpiralSearch() GopherMapController {

	settings := world.GopherMapSettings{
		Dimensions:      world.Dimensions{Width: 3000, Height: 3000},
		Population:      world.Population{InitialPopulation: 5000, MaxPopulation: 1000000},
		NumberOfFood:    1000000,
		GopherBirthRate: 7,
	}

	gWorld := world.CreateWorldCustom(settings)
	renderer := renderers.NewRenderer(45, 15)
	return GopherMapController{
		GopherWorld:  gWorld,
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

	gWorld := world.CreatePartitionTileMapCustom(settings)
	renderer := renderers.NewRenderer(100, 100)
	return GopherMapController{
		GopherWorld:  gWorld,
		GridRenderer: &renderer,
		CreateNew:    world.CreatePartitionTileMapCustom,
	}
}

//Start Initiates the controller. If the Map does not exist. The Map will be built
func (controller *GopherMapController) Start() {
	if controller.GopherWorld == nil {
		controller.GopherWorld = controller.CreateNew(*controller.GopherMapSettings)
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

	controller.GopherWorld.Add(action)
}

func (controller *GopherMapController) KeyPress(key Keys) {

	switch key {
	case WKey:
		controller.GopherWorld.Add(func() {
			controller.UnSelectGopher()
		})
	case QKey:
		controller.GopherWorld.Add(func() {
			controller.SelectRandomGopher()
		})
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
			if controller.GopherWorld.SelectedGopher != nil {
				isSelected = controller.GopherWorld.SelectedGopher.Position.GetX() == x &&
					controller.GopherWorld.SelectedGopher.Position.GetY() == y
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

	render.TextBelowCanvas = renderString

	gmr := GopherMapRender{
		Render: render,
	}
	if controller.SelectedGopher != nil {
		gmr.SelectedGopher = controller.SelectedGopher
	} else {
		gmr.SelectedGopher = &world.Gopher{}
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
		PageTitle:     "G O P H E R L I F E <b>2.0</b>",
		FormData:      formdataArray,
		IsGopherWorld: true,
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
		controller.GopherWorld = gmc

	}

	return true
}

type SpiralMapController struct {
	NoPlayerInput
	world.SpiralMapSettings
	*world.SpiralMap
	*renderers.GridRenderer
}

func NewSpiralMapController() SpiralMapController {

	settings := world.SpiralMapSettings{
		Dimensions:    world.Dimensions{Width: 50, Height: 50},
		MaxPopulation: 1000,
		WeirdSpiral:   false,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return SpiralMapController{
		GridRenderer:      &renderer,
		SpiralMapSettings: settings,
	}
}

func NewWeirdSpiralMapController() SpiralMapController {

	settings := world.SpiralMapSettings{
		Dimensions:    world.Dimensions{Width: 50, Height: 50},
		MaxPopulation: 1000,
		WeirdSpiral:   true,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return SpiralMapController{
		GridRenderer:      &renderer,
		SpiralMapSettings: settings,
	}
}

func (controller *SpiralMapController) Start() {
	if controller.SpiralMap == nil {
		sMap := world.NewSpiralMap(controller.SpiralMapSettings)
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
	return WorldPageData{IsGopherWorld: false}
}

func (controller *SpiralMapController) HandleForm(values url.Values) bool {
	controller.SpiralMap = nil
	controller.Start()
	return true
}

type FireWorksController struct {
	NoPlayerInput
	world.GopherMapSettings
	*world.GopherWorld
	*renderers.GridRenderer
}

func NewFireWorksController() FireWorksController {

	settings := world.GopherMapSettings{
		Dimensions:      world.Dimensions{Width: 400, Height: 200},
		Population:      world.Population{InitialPopulation: 2000, MaxPopulation: 100000},
		NumberOfFood:    2500,
		GopherBirthRate: 35,
	}

	renderer := renderers.NewRenderer(400, 150)
	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)
	renderer.TileHeight = 2
	renderer.TileWidth = 2

	return FireWorksController{
		GopherMapSettings: settings,
		GridRenderer:      &renderer,
	}
}

func (controller *FireWorksController) Start() {
	if controller.GopherWorld == nil {
		controller.GopherWorld = world.CreateWorldCustom(controller.GopherMapSettings)
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
	}

	return color.RGBA{0, 0, 0, 1}

}

func (controller *FireWorksController) PageLayout() WorldPageData {
	return WorldPageData{}
}

func (controller *FireWorksController) HandleForm(values url.Values) bool {
	return true
}
