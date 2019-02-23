package world

import (
	"encoding/json"
	"fmt"
	"image/color"
)

type Controllable interface {
	Click(x int, y int)
	KeyPress(key Keys)
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
}

func NewGopherMapWithSpiralSearch(stats Statistics) GopherMapController {
	gMap := CreateWorldCustom(stats)
	renderer := NewRenderer()
	return GopherMapController{
		GopherMap: gMap,
		Renderer:  &renderer,
	}
}

func NewGopherMapWithParitionGridAndSearch(stats Statistics) GopherMapController {
	gMap := CreatePartitionTileMapCustom(stats)
	renderer := NewRenderer()
	return GopherMapController{
		GopherMap: gMap,
		Renderer:  &renderer,
	}
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *GopherMapController) Click(x int, y int) {

	action := func() {
		controller.SelectEntity(x, y)
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

			switch tile.Gopher.Gender {
			case Male:

				if isSelected {
					return maleGopherSelectedColor
				} else {
					return maleGopherColor
				}
			case Female:
				if isSelected {
					return femaleGopherSelectedColor
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

type SpiralMapController struct {
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

	sMap := NewSpiralMap(stats)
	renderer := NewRendererSetUp(stats.Width*2, stats.Height*2)

	//	selectedTile.Gopher.Position.GetX() - renderer.RenderSizeX/2
	//	renderer.StartY = selectedTile.Gopher.Position.GetY() - renderer.RenderSizeY/2

	renderer.ShiftRenderer(stats.Width/2-renderer.RenderSizeX/2, stats.Height/2-renderer.RenderSizeY/2)

	return SpiralMapController{
		SpiralMap: &sMap,
		Renderer:  &renderer,
	}
}

//Click selects the tile on the gopher map and runs the SelectEntity method
func (controller *SpiralMapController) Click(x int, y int) {
}

func (controller *SpiralMapController) KeyPress(key Keys) {
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
