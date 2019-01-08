package world

import (
	"fmt"
)

type Renderer struct {
	StartX int
	StartY int

	RenderSizeX int
	RenderSizeY int

	TileMap *TileMap
}

type Render struct {
	Grid [][]*RenderTile

	WorldRender    string
	SelectedGopher *Gopher
}

type RenderTile struct {
	Color string
}

const (
	maleGopherColor           = "#5adaff"
	maleGopherSelectedColor   = "#f5f5f5"
	femaleGopherColor         = "#FFFF00"
	femaleGopherSelectedColor = "#ff9b9a"
	foodColor                 = "#cc7000"
	decayedGopherColor        = "#000000"
)

//NewRenderer returns a new Render struct of size 45x and 15 y
func NewRenderer() Renderer {
	return Renderer{RenderSizeX: 250, RenderSizeY: 200}
}

func (renderer *Renderer) RenderWorld(tileMap *TileMap) Render {

	render := Render{WorldRender: "", SelectedGopher: &Gopher{}}

	render.Grid = make([][]*RenderTile, renderer.RenderSizeX)

	for i := 0; i < renderer.RenderSizeX; i++ {
		render.Grid[i] = make([]*RenderTile, renderer.RenderSizeY)

		for j := 0; j < renderer.RenderSizeY; j++ {
			renderTile := RenderTile{Color: "Hello"}
			render.Grid[i][j] = &renderTile
		}
	}

	renderString := ""

	startX := 0
	startY := 0

	if tileMap.SelectedGopher != nil {
		render.SelectedGopher = tileMap.SelectedGopher
		renderer.StartX = tileMap.SelectedGopher.Position.GetX() - renderer.RenderSizeX/2
		renderer.StartY = tileMap.SelectedGopher.Position.GetY() - renderer.RenderSizeY/2
	}

	startX = renderer.StartX
	startY = renderer.StartY

	for y := startY; y < startY+renderer.RenderSizeY; y++ {

		for x := startX; x < startX+renderer.RenderSizeX; x++ {

			renderTile := render.Grid[x-startX][y-startY]
			renderTile.Color = "#41770f"

			if mapPoint, ok := tileMap.GetTile(x, y); ok {

				switch {
				case mapPoint.isEmpty():
					renderTile.Color = "#41770f"
				case mapPoint.Gopher != nil:
					isSelected := false
					if tileMap.SelectedGopher != nil {
						isSelected = tileMap.SelectedGopher.Position.GetX() == x && tileMap.SelectedGopher.Position.GetY() == y
					}

					switch mapPoint.Gopher.Gender {
					case Male:

						if isSelected {
							renderTile.Color = maleGopherSelectedColor
						} else {
							renderTile.Color = maleGopherColor
						}
					case Female:
						if isSelected {
							renderTile.Color = femaleGopherSelectedColor
						} else {
							renderTile.Color = femaleGopherColor
						}
					}

					if mapPoint.Gopher.IsDead {
						renderTile.Color = decayedGopherColor
					}

					if !mapPoint.Gopher.IsMature() {
						//text = strings.ToLower(text)
					}
				case mapPoint.Food != nil:
					renderTile.Color = foodColor
				}
			} else {
				renderTile.Color = "#89cff0"
			}

		}
	}
	renderString += "<br />"
	renderString += fmt.Sprintf("<span>Number of Gophers: %d </span><br />", len(tileMap.gopherArray))
	renderString += fmt.Sprintf("<span>Avg Processing Time (s): %s </span><br />", tileMap.processStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Avg Gopher Time (s): %s </span><br />", tileMap.gopherStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span; >Avg Input Time (s): %s </span><br />", tileMap.inputStopWatch.GetAverage().String())
	renderString += fmt.Sprintf("<span>Total Elasped Time (s): %s </span><br />", tileMap.globalStopWatch.GetCurrentElaspedTime().String())

	render.WorldRender = renderString

	return render
}

//ShiftRenderer moves the current rendering scope of the renderer by the given x and y values
func (renderer *Renderer) ShiftRenderer(x int, y int) {
	renderer.StartX += x
	renderer.StartY += y
}
