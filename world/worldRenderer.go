package world

import (
	"fmt"
)

type Renderer struct {
	StartX int
	StartY int

	RenderSizeX int
	RenderSizeY int

	World *World
}

type Render struct {
	Grid [][]*Tile

	WorldRender    string
	SelectedGopher *Gopher
}

type Tile struct {
	Color string
}

type span struct {
	class string
	id    string
	text  string
	style string
}

var maleGopherColor = "#5adaff"
var maleGopherSelectedColor = "#f5f5f5"
var femaleGopherColor = "#FFFF00"
var femaleGopherSelectedColor = "#ff9b9a"
var foodColor = "#cc7000"
var decayedGopherColor = "#000000"

//NewRenderer returns a new Render struct of size 45x and 15 y
func NewRenderer() Renderer {
	return Renderer{RenderSizeX: 250, RenderSizeY: 200}
}

func addSpanTagToRender(span span) string {
	return fmt.Sprintf("<a id='%s' class='%s' style='%s'>%s</a>", span.id, span.class, span.style, span.text)
}

func (renderer *Renderer) RenderWorld(world *World) Render {

	render := Render{WorldRender: "", SelectedGopher: &Gopher{}}

	render.Grid = make([][]*Tile, renderer.RenderSizeX)

	for i := 0; i < renderer.RenderSizeX; i++ {
		render.Grid[i] = make([]*Tile, renderer.RenderSizeY)

		for j := 0; j < renderer.RenderSizeY; j++ {
			tile := Tile{Color: "Hello"}
			render.Grid[i][j] = &tile
		}
	}

	renderString := ""

	startX := 0
	startY := 0

	if world.SelectedGopher != nil {
		render.SelectedGopher = world.SelectedGopher
		renderer.StartX = world.SelectedGopher.Position.GetX() - renderer.RenderSizeX/2
		renderer.StartY = world.SelectedGopher.Position.GetY() - renderer.RenderSizeY/2
	}

	startX = renderer.StartX
	startY = renderer.StartY

	for y := startY; y < startY+renderer.RenderSizeY; y++ {

		for x := startX; x < startX+renderer.RenderSizeX; x++ {

			tile := render.Grid[x-startX][y-startY]
			tile.Color = "#41770f"

			if mapPoint, ok := world.GetMapPoint(x, y); ok {

				switch {
				case mapPoint.isEmpty():
					tile.Color = "#41770f"
				case mapPoint.Gopher != nil:
					isSelected := false
					if world.SelectedGopher != nil {
						isSelected = world.SelectedGopher.Position.GetX() == x && world.SelectedGopher.Position.GetY() == y
					}

					switch mapPoint.Gopher.Gender {
					case Male:

						if isSelected {
							tile.Color = maleGopherSelectedColor
						} else {
							tile.Color = maleGopherColor
						}
					case Female:
						if isSelected {
							tile.Color = femaleGopherSelectedColor
						} else {
							tile.Color = femaleGopherColor
						}
					}

					if mapPoint.Gopher.IsDead {
						tile.Color = decayedGopherColor
					}

					if !mapPoint.Gopher.IsMature() {
						//text = strings.ToLower(text)
					}
				case mapPoint.Food != nil:
					tile.Color = foodColor
				}
			} else {
				tile.Color = "#89cff0"
			}

		}
	}
	renderString += "<br />"
	renderString += fmt.Sprintf("<span>Number of Gophers: %d </span><br />", len(world.gopherArray))
	renderString += fmt.Sprintf("<span>Avg Processing Time (s): %s </span><br />", world.processStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Avg Gopher Time (s): %s </span><br />", world.gopherStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span; >Avg Input Time (s): %s </span><br />", world.inputStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Total Elasped Time (s): %s </span><br />", world.globalStopWatch.GetCurrentElaspedTime().String())

	render.WorldRender = renderString

	return render
}

//ShiftRenderer moves the current rendering scope of the renderer by the given x and y values
func (renderer *Renderer) ShiftRenderer(x int, y int) {
	renderer.StartX += x
	renderer.StartY += y
}
