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
	return Renderer{RenderSizeX: 30, RenderSizeY: 20, IsPaused: false}
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
		startX = world.SelectedGopher.Position.GetX() - renderer.RenderSizeX/2
		startY = world.SelectedGopher.Position.GetY() - renderer.RenderSizeY/2
	} else {
		startX = renderer.StartX
		startY = renderer.StartY
	}

	for y := startY; y < startY+renderer.RenderSizeY; y++ {

		for x := startX; x < startX+renderer.RenderSizeX; x++ {

			key := math.CoordinateMapKey(x, y)

			if mapPoint, ok := world.world[key]; ok {

				switch {
				case mapPoint.isEmpty():
					renderString += addSpanTagToRender(span{text: "/", class: "grass"})
				case mapPoint.Gopher != nil:

					isSelected := false

					if world.SelectedGopher != nil {
						isSelected = world.SelectedGopher.Position.MapKey() == key
					}
					if isSelected {
						renderString += addSpanTagToRender(span{text: "G", id: key, class: "selected"})
					} else if mapPoint.Gopher.IsDead() {
						renderString += addSpanTagToRender(span{text: "X", id: key, class: "gopher"})
					} else {
						renderString += addSpanTagToRender(span{text: "G", id: key, class: "gopher"})
					}

				case mapPoint.Food != nil:
					renderString += addSpanTagToRender(span{text: "!", id: key, class: "food"})
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
