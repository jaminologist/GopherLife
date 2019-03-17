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

type SnakeMapController struct {
	world.SnakeMapSettings
	ClickToBegin bool
	*world.SnakeMap
	*renderers.GridRenderer
}

func NewSnakeMapController() SnakeMapController {

	d := world.Dimensions{
		Width:  35,
		Height: 35,
	}

	renderer := renderers.NewRenderer(50, 50)
	renderer.Shift(d.Width/2-renderer.Width/2, d.Height/2-renderer.Height/2)
	renderer.TileWidth = 10
	renderer.TileHeight = 10

	return SnakeMapController{
		SnakeMapSettings: world.SnakeMapSettings{d, 5},
		GridRenderer:     &renderer,
	}
}

func (controller *SnakeMapController) Start() {
	if controller.SnakeMap == nil {
		sMap := world.NewSnakeMap(controller.SnakeMapSettings)
		controller.SnakeMap = &sMap
	}
}

func (controller *SnakeMapController) MarshalJSON() ([]byte, error) {

	render := controller.GridRenderer.Draw(controller)
	render.TextBelowCanvas += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	if controller.SnakeMap.IsGameOver {
		render.TextBelowCanvas += fmt.Sprintf("<span>Game Over!</span><br />")
	} else if !controller.ClickToBegin {
		render.TextBelowCanvas += fmt.Sprintf("<span>Click to Begin")
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
			FormDataSnakeSlowDown(controller.SnakeMapSettings.SpeedReduction, 3),
		},
	}
}

func (controller *SnakeMapController) HandleForm(values url.Values) bool {

	fd := FormDataSnakeSlowDown(0, 0)
	if strings.Contains(values.Encode(), fd.Name) {
		speedReduction, _ := strconv.ParseInt(values.Get(fd.Name), 10, 64)
		controller.SnakeMapSettings.SpeedReduction = int(speedReduction)
		controller.SnakeMap = nil
		controller.Start()
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
