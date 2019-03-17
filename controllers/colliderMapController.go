package controllers

import (
	"encoding/json"
	"gopherlife/renderers"
	"gopherlife/world"
	"image/color"
	"net/url"
	"strconv"
	"strings"
)

type CollisionMapController struct {
	NoPlayerInput
	world.CollisionMapSettings
	*world.CollisionMap
	*renderers.GridRenderer
	CreateNew func(world.CollisionMapSettings) world.CollisionMap
}

func NewCollisionMapController() CollisionMapController {

	settings := world.CollisionMapSettings{
		Dimensions: world.Dimensions{Width: 75, Height: 75},
		Population: world.Population{InitialPopulation: 500},
		IsDiagonal: false,
	}

	renderer := renderers.NewRenderer(100, 100)
	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return CollisionMapController{
		CollisionMapSettings: settings,
		GridRenderer:         &renderer,
		CreateNew:            world.NewCollisionMap,
	}
}

func NewDiagonalCollisionMapController() CollisionMapController {

	settings := world.CollisionMapSettings{
		Dimensions: world.Dimensions{Width: 75, Height: 75},
		Population: world.Population{InitialPopulation: 500},
		IsDiagonal: true,
	}

	renderer := renderers.NewRenderer(100, 100)

	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return CollisionMapController{
		CollisionMapSettings: settings,
		GridRenderer:         &renderer,
		CreateNew:            world.NewCollisionMap,
	}
}

func (controller *CollisionMapController) Start() {
	if controller.CollisionMap == nil {
		sMap := controller.CreateNew(controller.CollisionMapSettings)
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
	settings := controller.CollisionMapSettings

	formdataArray := []FormData{
		FormDataWidth(settings.Width, 2),
		FormDataHeight(settings.Height, 2),
		FormDataInitialPopulation(settings.InitialPopulation, 2),
	}

	return WorldPageData{
		PageTitle: "C O L L I D E R L I F E",
		FormData:  formdataArray,
	}
}

func (controller *CollisionMapController) HandleForm(values url.Values) bool {

	if strings.Contains(values.Encode(), "initialPopulation") {

		width, _ := strconv.ParseInt(values.Get("width"), 10, 64)
		height, _ := strconv.ParseInt(values.Get("height"), 10, 64)
		initialPopulation, _ := strconv.ParseInt(values.Get("initialPopulation"), 10, 64)

		controller.CollisionMapSettings.Width = int(width)
		controller.CollisionMapSettings.Height = int(height)
		controller.CollisionMapSettings.InitialPopulation = int(initialPopulation)

		gmc := world.NewCollisionMap(controller.CollisionMapSettings)
		controller.CollisionMap = &gmc

	}
	return true
}
