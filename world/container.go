package world

import (
	"fmt"
	"gopherlife/calc"
)

type TileContainer interface {
	Tile(x int, y int) (*Tile, bool)
}

type Basic2DContainer struct {
	grid   [][]*Tile
	x      int
	y      int
	width  int
	height int
}

func NewBasic2DContainer(x int, y int, width int, height int) Basic2DContainer {

	container := Basic2DContainer{width: width, height: height}

	container.grid = make([][]*Tile, width)

	for i := 0; i < width; i++ {
		container.grid[i] = make([]*Tile, height)

		for j := 0; j < height; j++ {
			tile := Tile{nil, nil}
			container.grid[i][j] = &tile
		}
	}

	return container
}

func (container *Basic2DContainer) Tile(x int, y int) (*Tile, bool) {

	if x < container.x || x >= container.width || y < container.y || y >= container.height {
		return nil, false
	}

	return container.grid[x][y], true
}

type TrackedTileContainer struct {
	x      int
	y      int
	width  int
	height int
	TileContainer
	gopherTileLocations map[int]*Tile
	foodTileLocations   map[int]*Tile
	Insertable
}

func NewTrackedTileContainer(x int, y int, width int, height int) TrackedTileContainer {
	b2dc := NewBasic2DContainer(x, y, width, height)

	return TrackedTileContainer{
		x:                   x,
		y:                   y,
		width:               width,
		height:              height,
		TileContainer:       &b2dc,
		gopherTileLocations: make(map[int]*Tile),
		foodTileLocations:   make(map[int]*Tile),
	}
}

func (container *TrackedTileContainer) Tile(x int, y int) (*Tile, bool) {
	return container.TileContainer.Tile(x, y)
}

func (container *TrackedTileContainer) ConvertToTrackedTileCoordinates(x int, y int) (gridX int, gridY int) {
	return (x - container.x), (y - container.y)
}

func (container *TrackedTileContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasGopher() {
			tile.SetGopher(gopher)
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			container.gopherTileLocations[calc.Hashcode(x, y)] = tile
			return true
		}
	}

	return false

}

func (container *TrackedTileContainer) InsertFood(x int, y int, food *Food) bool {
	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasFood() {
			tile.SetFood(food)
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			container.foodTileLocations[calc.Hashcode(x, y)] = tile
			return true
		}
	}

	return false

}

func (container *TrackedTileContainer) RemoveGopher(x int, y int, gopher *Gopher) bool {
	if tile, ok := container.Tile(x, y); ok {
		if tile.HasGopher() {
			tile.ClearGopher()
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			delete(container.foodTileLocations, calc.Hashcode(x, y))
			return true
		}
	}
	return false
}

func (container *TrackedTileContainer) RemoveFood(x int, y int, food *Food) bool {
	if tile, ok := container.Tile(x, y); ok {
		if tile.HasFood() {
			tile.ClearFood()
			x, y = container.ConvertToTrackedTileCoordinates(x, y)
			delete(container.foodTileLocations, calc.Hashcode(x, y))
			return true
		}
	}
	return false
}

type GridContainer interface {
	TileContainer
	Grid(x int, y int) (TileContainer, bool)
}

type BasicGridContainer struct {
	containers [][]*TrackedTileContainer
	gridWidth  int
	gridHeight int
	width      int
	height     int
}

func NewBasicGridContainer(width int, height int, gridWidth int, gridHeight int) BasicGridContainer {

	numberOfGridsX := width / gridWidth

	if numberOfGridsX*gridWidth < width {
		numberOfGridsX++
	}

	numberOfGridsY := height / gridHeight

	if numberOfGridsY*gridHeight < height {
		numberOfGridsY++
	}

	containers := make([][]*TrackedTileContainer, numberOfGridsX)

	for i := 0; i < numberOfGridsX; i++ {
		containers[i] = make([]*TrackedTileContainer, numberOfGridsY)

		for j := 0; j < numberOfGridsY; j++ {

			fmt.Println(i, j)
			ttc := NewTrackedTileContainer(i*numberOfGridsX,
				j*numberOfGridsY,
				gridWidth,
				gridHeight)
			containers[i][j] = &ttc
		}
	}

	return BasicGridContainer{
		containers: containers,
		gridWidth:  gridWidth,
		gridHeight: gridHeight,
		width:      width,
		height:     height,
	}
}

func (container *BasicGridContainer) Tile(x int, y int) (*Tile, bool) {
	if grid, ok := container.Grid(x, y); ok {
		if tile, ok := grid.Tile(x, y); ok {
			return tile, ok
		}
	}
	return nil, false
}

//Takes an X and Y Position, and finds which grid it should be in
func (container *BasicGridContainer) Grid(x int, y int) (*TrackedTileContainer, bool) {

	if x < 0 || x >= container.width || y < 0 || y >= container.height {
		return nil, false
	}

	x, y = container.convertToGridCoordinates(x, y)

	val := container.containers[x][y]

	if val != nil {
		return val, true
	}

	return nil, false
}

func (container *BasicGridContainer) convertToGridCoordinates(x int, y int) (int, int) {
	gridX, gridY := x/container.gridWidth, y/container.gridHeight
	return gridX, gridY
}

type Insertable interface {
	InsertGopher(x int, y int, gopher *Gopher) bool
	InsertFood(x int, y int, food *Food) bool
	RemoveGopher(x int, y int, gopher *Gopher) bool
	RemoveFood(x int, y int, food *Food) (Food bool)
}

type GridInsertable struct {
	Grid [][]GridContainer
}

func (container *BasicGridContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if grid, ok := container.Grid(x, y); ok {
		gopher.Position.X = x
		gopher.Position.Y = y
		return grid.InsertGopher(x, y, gopher)
	}
	return false
}

func (container *BasicGridContainer) InsertFood(x int, y int, food *Food) bool {
	if grid, ok := container.Grid(x, y); ok {
		food.Position.X = x
		food.Position.Y = y
		return grid.InsertFood(x, y, food)
	}
	return false
}

func (container *BasicGridContainer) RemoveGopher(x int, y int, gopher *Gopher) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveGopher(x, y, gopher)
	}
	return false
}

func (container *BasicGridContainer) RemoveFood(x int, y int, food *Food) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveFood(x, y, food)
	}
	return false
}
