package world

import (
	"fmt"
	"gopherlife/math"
)

type Renderer struct {
	StartX int
	StartY int

	RenderSizeX int
	RenderSizeY int

	IsPaused bool

	World *World
}

type Render struct {
	WorldRender    string
	SelectedGopher *Gopher
}

type span struct {
	class string
	id    string
	text  string
	style string
}

//NewRenderer returns a new Render struct of size 45x and 15 y
func NewRenderer() Renderer {
	return Renderer{RenderSizeX: 45, RenderSizeY: 15, IsPaused: false}
}

func addSpanTagToRender(span span) string {
	return fmt.Sprintf("<a id='%s' class='%s' style='%s'>%s</a>", span.id, span.class, span.style, span.text)
}

func (renderer *Renderer) RenderWorld(world *World) Render {

	render := Render{WorldRender: "", SelectedGopher: &Gopher{}}
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

			key := math.CoordinateMapKey(x, y)

			if mapPoint, ok := world.world[key]; ok {

				switch {
				case mapPoint.isEmpty():
					renderString += addSpanTagToRender(span{text: "/", class: "grass interactable"})
				case mapPoint.Gopher != nil:

					//isSelected := false

					text := "G"
					class := "gopher interactable"
					style := "color:#ffffff"

					if world.SelectedGopher != nil {
						//	isSelected = world.SelectedGopher.Position.MapKey() == key
						class = "selected interactable"
					}

					if mapPoint.Gopher.IsDead() {
						text = "X"
					}

					switch mapPoint.Gopher.Gender {
					case Male:
						style = "color:#A9A9A9"
					case Female:
						style = "color:#FF0080"
					}

					renderString += addSpanTagToRender(span{text: text, id: key, class: class, style: style})

				case mapPoint.Food != nil:
					renderString += addSpanTagToRender(span{text: "!", id: key, class: "food interactable"})
				}
			} else {
				renderString += "<span style='color:#89cff0'; >X</span>"
			}

		}
		renderString += "<br />"

	}
	renderString += "<br />"
	renderString += fmt.Sprintf("<span style='color:#ffffff'; >Number of Gophers: %d </span>", len(world.gopherArray))
	render.WorldRender = renderString

	return render
}

//ShiftRenderer moves the current rendering scope of the renderer by the given x and y values
func (renderer *Renderer) ShiftRenderer(x int, y int) {
	renderer.StartX += x
	renderer.StartY += y
}
