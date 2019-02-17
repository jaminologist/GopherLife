package world

import "encoding/json"

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

func (controller *GopherMapController) MarshalJSON() ([]byte, error) {
	return json.Marshal(controller.Renderer.RenderWorld(controller))
}
