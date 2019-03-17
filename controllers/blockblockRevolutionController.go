package controllers

import (
	"encoding/json"
	"gopherlife/colors"
	"gopherlife/renderers"
	"gopherlife/world"
	"image/color"
	"net/url"
	"strconv"
	"strings"
)

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
	renderer.TileWidth = 15
	renderer.TileHeight = 15

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
	//render.TextBelowCanvas += fmt.Sprintf("<span>Score: %d </span><br />", controller.Score)

	/*if controller.SnakeMap.IsGameOver {
		render.TextBelowCanvas += fmt.Sprintf("<span>Game Over!</span><br />")
	} else if !controller.ClickToBegin {
		render.TextBelowCanvas += fmt.Sprintf("<span>Click to Begin")
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
