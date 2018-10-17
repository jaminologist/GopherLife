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
}

func NewRenderer() Renderer {
	return Renderer{RenderSizeX: 30, RenderSizeY: 25, IsPaused: false}
}

func addSpanTagToRender(span span) string {
	return fmt.Sprintf("<a id='%s' class='%s' >%s</a>", span.id, span.class, span.text)
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

					isSelected := false

					if world.SelectedGopher != nil {
						isSelected = world.SelectedGopher.Position.MapKey() == key
					}
					if isSelected {
						renderString += addSpanTagToRender(span{text: "G", id: key, class: "selected interactable"})
					} else if mapPoint.Gopher.IsDead() {
						renderString += addSpanTagToRender(span{text: "X", id: key, class: "gopher interactable"})
					} else {
						renderString += addSpanTagToRender(span{text: "G", id: key, class: "gopher interactable"})
					}

				case mapPoint.Food != nil:
					renderString += addSpanTagToRender(span{text: "!", id: key, class: "food interactable"})
				}
			} else {
				renderString += "<span style='color:#89cff0'; >X</span>"
			}

		}
		renderString += "<br />"
	}

	render.WorldRender = renderString

	return render
}

func (renderer *Renderer) ShiftRenderer(x int, y int) {
	renderer.StartX += x
	renderer.StartY += y
}