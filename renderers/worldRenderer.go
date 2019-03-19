package renderers

import (
	"image/color"
)

type GridRenderer struct {
	StartX int
	StartY int

	Width  int
	Height int

	TileWidth  int
	TileHeight int
}

type Render struct {
	Grid            [][]*RenderTile
	TextBelowCanvas string
	StartX          int
	StartY          int
	TileWidth       int
	TileHeight      int
}

type RenderTile struct {
	color.RGBA
}

//RenderTileContainer when given an x and y value a color should be returned
type RenderTileContainer interface {
	RenderTile(x int, y int) color.RGBA
}

//NewRenderer returns a new Render struct of size 45x and 15 y
func NewRenderer(width int, height int) GridRenderer {
	return GridRenderer{Width: width, Height: height, TileWidth: 5, TileHeight: 5}
}

//Draw returns a struct containing all colors found within the dimensions of the Renderer
func (renderer *GridRenderer) Draw(container RenderTileContainer) Render {

	render := Render{
		TextBelowCanvas: "",
		StartX:          renderer.StartX,
		StartY:          renderer.StartY,
		TileWidth:       renderer.TileWidth,
		TileHeight:      renderer.TileHeight,
	}

	render.Grid = make([][]*RenderTile, renderer.Width)

	for i := 0; i < renderer.Width; i++ {
		render.Grid[i] = make([]*RenderTile, renderer.Height)

		for j := 0; j < renderer.Height; j++ {
			renderTile := RenderTile{}
			render.Grid[i][j] = &renderTile
		}
	}

	startX := renderer.StartX
	startY := renderer.StartY

	for y := startY; y < startY+renderer.Height; y++ {
		for x := startX; x < startX+renderer.Width; x++ {
			render.Grid[x-startX][y-startY].RGBA = container.RenderTile(x, y)
		}
	}

	return render
}

//Shift Moves the Starting X and Y of the renderer
func (gr *GridRenderer) Shift(x int, y int) {
	gr.StartX += x
	gr.StartY += y
}

func (gr *GridRenderer) Scroll(deltaY int) {

	if deltaY < 0 { //Wheel Up
		gr.TileWidth = gr.TileWidth + 1
		gr.TileHeight = gr.TileHeight + 1
	} else { //Wheel Down
		gr.TileWidth = gr.TileWidth - 1
		gr.TileHeight = gr.TileHeight - 1

		if gr.TileWidth <= 1 || gr.TileHeight <= 1 {
			gr.TileWidth = 1
			gr.TileHeight = 1
		}
	}

}
