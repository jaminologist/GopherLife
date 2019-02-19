package world

import (
	"image/color"
)

var (
	maleGopherColor           = color.RGBA{90, 218, 255, 1}
	maleGopherSelectedColor   = color.RGBA{245, 245, 245, 1}
	femaleGopherColor         = color.RGBA{255, 255, 0, 1}
	femaleGopherSelectedColor = color.RGBA{255, 155, 154, 1}
	foodColor                 = color.RGBA{204, 112, 0, 1}
	decayedGopherColor        = color.RGBA{0, 0, 0, 1}
	grassColor                = color.RGBA{65, 119, 15, 1}
)

type Renderer struct {
	StartX int
	StartY int

	RenderSizeX int
	RenderSizeY int
}

type Render struct {
	Grid           [][]*RenderTile
	WorldRender    string
	SelectedGopher *Gopher
}

type RenderTile struct {
	color.RGBA
}

type Renderable interface {
	RenderTile(x int, y int) color.RGBA
}

//NewRenderer returns a new Render struct of size 45x and 15 y
func NewRenderer() Renderer {
	s := color.RGBA{0, 0, 0, 0}
	_ = s
	return Renderer{RenderSizeX: 100, RenderSizeY: 100}
}

func NewRendererSetUp(width int, height int) Renderer {
	return Renderer{RenderSizeX: width, RenderSizeY: height}
}

func (renderer *Renderer) RenderWorld(tileMap Renderable) Render {

	render := Render{WorldRender: "", SelectedGopher: &Gopher{}}

	render.Grid = make([][]*RenderTile, renderer.RenderSizeX)

	for i := 0; i < renderer.RenderSizeX; i++ {
		render.Grid[i] = make([]*RenderTile, renderer.RenderSizeY)

		for j := 0; j < renderer.RenderSizeY; j++ {
			renderTile := RenderTile{}
			render.Grid[i][j] = &renderTile
		}
	}

	startX := renderer.StartX
	startY := renderer.StartY

	for y := startY; y < startY+renderer.RenderSizeY; y++ {
		for x := startX; x < startX+renderer.RenderSizeX; x++ {
			render.Grid[x-startX][y-startY].RGBA = tileMap.RenderTile(x, y)
		}
	}

	return render
}

//ShiftRenderer moves the current rendering scope of the renderer by the given x and y values
func (renderer *Renderer) ShiftRenderer(x int, y int) {
	renderer.StartX += x
	renderer.StartY += y
}
