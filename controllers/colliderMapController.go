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

type CollisionWorldController struct {
	NoPlayerInput
	world.CollisionWorldSettings
	*world.CollisionWorld
	*renderers.GridRenderer
	CreateNew func(world.CollisionWorldSettings) world.CollisionWorld
}

func NewCollisionWorldController() CollisionWorldController {

	settings := world.CollisionWorldSettings{
		Dimensions: world.Dimensions{Width: 75, Height: 75},
		Population: world.Population{InitialPopulation: 500},
		IsDiagonal: false,
	}

	renderer := renderers.NewRenderer(100, 100)
	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return CollisionWorldController{
		CollisionWorldSettings: settings,
		GridRenderer:         &renderer,
		CreateNew:            world.NewCollisionWorld,
	}
}

func NewDiagonalCollisionWorldController() CollisionWorldController {

	settings := world.CollisionWorldSettings{
		Dimensions: world.Dimensions{Width: 75, Height: 75},
		Population: world.Population{InitialPopulation: 500},
		IsDiagonal: true,
	}

	renderer := renderers.NewRenderer(100, 100)

	renderer.Shift(settings.Width/2-renderer.Width/2, settings.Height/2-renderer.Height/2)

	return CollisionWorldController{
		CollisionWorldSettings: settings,
		GridRenderer:         &renderer,
		CreateNew:            world.NewCollisionWorld,
	}
}

func (controller *CollisionWorldController) Start() {
	if controller.CollisionWorld == nil {
		sMap := controller.CreateNew(controller.CollisionWorldSettings)
		controller.CollisionWorld = &sMap
	}
}

func (controller *CollisionWorldController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.GridRenderer.Draw(controller))
}

func (controller *CollisionWorldController) RenderTile(x int, y int) color.RGBA {

	if controller.Contains(x, y) {
		if c, ok := controller.HasCollider(x, y); ok {
			return c.Color
		} else {
			return color.RGBA{0, 0, 0, 1}
		}
	}

	return color.RGBA{255, 255, 255, 1}
}

func (controller *CollisionWorldController) PageLayout() WorldPageData {
	settings := controller.CollisionWorldSettings

	formdataArray := []FormData{
		FormDataWidth(settings.Width, 2),
		FormDataHeight(settings.Height, 2),
		FormDataInitialPopulation(settings.InitialPopulation, 2),
	}

	return WorldPageData{
		FormData: formdataArray,
	}
}

func (controller *CollisionWorldController) HandleForm(values url.Values) bool {

	if strings.Contains(values.Encode(), "initialPopulation") {

		width, _ := strconv.ParseInt(values.Get("width"), 10, 64)
		height, _ := strconv.ParseInt(values.Get("height"), 10, 64)
		initialPopulation, _ := strconv.ParseInt(values.Get("initialPopulation"), 10, 64)

		controller.CollisionWorldSettings.Width = int(width)
		controller.CollisionWorldSettings.Height = int(height)
		controller.CollisionWorldSettings.InitialPopulation = int(initialPopulation)

		gmc := world.NewCollisionWorld(controller.CollisionWorldSettings)
		controller.CollisionWorld = &gmc

	}
	return true
}
