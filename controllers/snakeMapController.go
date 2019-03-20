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

type SnakeWorldController struct {
	world.SnakeWorldSettings
	ClickToBegin bool
	*world.SnakeWorld
	*renderers.GridRenderer
}

func NewSnakeWorldController() SnakeWorldController {

	d := world.Dimensions{
		Width:  35,
		Height: 35,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(d.Width/2-renderer.Width/2, d.Height/2-renderer.Height/2)
	renderer.TileWidth = 10
	renderer.TileHeight = 10

	return SnakeWorldController{
		SnakeWorldSettings: world.SnakeWorldSettings{d, 5},
		GridRenderer:     &renderer,
	}
}

func (controller *SnakeWorldController) Start() {
	if controller.SnakeWorld == nil {
		sMap := world.NewSnakeWorld(controller.SnakeWorldSettings)
		controller.SnakeWorld = &sMap
	}
}

func (controller *SnakeWorldController) MarshalJSON() ([]byte, error) {

	render := controller.GridRenderer.Draw(controller)
	render.TextBelowCanvas += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	if controller.SnakeWorld.IsGameOver {
		render.TextBelowCanvas += fmt.Sprintf("<span>Game Over!</span><br />")
	} else if !controller.ClickToBegin {
		render.TextBelowCanvas += fmt.Sprintf("<span>Click to Begin")
	}

	return json.Marshal(render)
}

func (controller *SnakeWorldController) RenderTile(x int, y int) color.RGBA {

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

func (controller *SnakeWorldController) PageLayout() WorldPageData {
	return WorldPageData{
		PageTitle: "E L O N G A T I N G G O P H E R L I F E",
		FormData: []FormData{
			FormDataSnakeSlowDown(controller.SnakeWorldSettings.SpeedReduction, 3),
		},
	}
}

func (controller *SnakeWorldController) HandleForm(values url.Values) bool {

	fd := FormDataSnakeSlowDown(0, 0)
	if strings.Contains(values.Encode(), fd.Name) {
		speedReduction, _ := strconv.ParseInt(values.Get(fd.Name), 10, 64)
		controller.SnakeWorldSettings.SpeedReduction = int(speedReduction)
		controller.SnakeWorld = nil
		controller.Start()
	}

	return true
}

func (controller *SnakeWorldController) KeyPress(key Keys) {
	switch key {
	case LeftArrow:
		controller.SnakeWorld.ChangeDirection(geometry.Left)
	case RightArrow:
		controller.SnakeWorld.ChangeDirection(geometry.Right)
	case UpArrow:
		controller.SnakeWorld.ChangeDirection(geometry.Up)
	case DownArrow:
		controller.SnakeWorld.ChangeDirection(geometry.Down)
	}
}

func (controller *SnakeWorldController) Update() bool {
	if controller.ClickToBegin {
		if controller.IsGameOver {
			controller.ClickToBegin = false
		}
		return controller.SnakeWorld.Update()
	}

	return true
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *SnakeWorldController) Click(x int, y int) {
	controller.ClickToBegin = true
}
