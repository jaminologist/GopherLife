package world

import (
	"gopherlife/math"
)

type Renderer struct {
	StartX int
	StartY int

	RenderSize int

	IsPaused bool

	World *World
}

func NewRenderer() Renderer {
	return Renderer{RenderSize: 25, IsPaused: false}
}

func (renderer *Renderer) RenderWorld(world *World) string {

	renderString := ""

	startX := 0
	startY := 0

	if world.SelectedGopher != nil {
		startX = world.SelectedGopher.Position.GetX() - renderer.RenderSize/2
		startY = world.SelectedGopher.Position.GetY() - renderer.RenderSize/2
	} else {
		startX = renderer.StartX
		startY = renderer.StartY
	}

	renderString += "<h1>"
	for y := startY; y < startY+renderer.RenderSize; y++ {

		for x := startX; x < startX+renderer.RenderSize; x++ {

			key := math.CoordinateMapKey(x, y)

			if mapPoint, ok := world.world[key]; ok {

				switch {
				case mapPoint.isEmpty():
					renderString += "<span class='grass'>/</span>"
				case mapPoint.Gopher != nil:

					isSelected := false

					if world.SelectedGopher != nil {
						isSelected = world.SelectedGopher.Position.MapKey() == key
					}
					if isSelected {
						renderString += "<a id='" + key + "' style='color:yellow;'>G</a>"
					} else if mapPoint.Gopher.IsDead() {
						renderString += "<span id='" + key + "' class='gopher'>X</span>"
					} else {
						renderString += "<span id='" + key + "' class='gopher'>G</span>"
					}

				case mapPoint.Food != nil:
					renderString += "<span class='food'>!</span>"
				}
			} else {
				renderString += "<span style='color:#89cff0'; >X</span>"
			}

		}
		renderString += "<br />"
	}

	renderString += "</h1>"

	return renderString
}
