package world

import (
	"gopherlife/geometry"
)

//TileContainer contains tiles that can be accessed using an x and y position
type TileContainer interface {
	Tile(x int, y int) (*GopherWorldTile, bool)
}

type GridContainer interface {
	TileContainer
	Grid(x int, y int) (TileContainer, bool)
}

//GopherContainer returns Gopher if the given x and y position has a Gopher
type GopherContainer interface {
	HasGopher(x int, y int) (*Gopher, bool)
}

//FoodContainer returns Food if the given x and y position has food
type FoodContainer interface {
	HasFood(x int, y int) (*Food, bool)
}

//Container returns if the given x and y position is within the container
type Container interface {
	Contains(x int, y int) bool
}

type Basic2DContainer struct {
	grid   [][]*GopherWorldTile
	x      int
	y      int
	width  int
	height int
}

func NewBasic2DContainer(x int, y int, width int, height int) Basic2DContainer {

	container := Basic2DContainer{
		x:      x,
		y:      y,
		width:  width,
		height: height}

	container.grid = make([][]*GopherWorldTile, width)

	for i := 0; i < width; i++ {
		container.grid[i] = make([]*GopherWorldTile, height)

		for j := 0; j < height; j++ {
			tile := GopherWorldTile{}
			container.grid[i][j] = &tile
		}
	}

	return container
}

func (container *Basic2DContainer) Tile(x int, y int) (*GopherWorldTile, bool) {
	if x < container.x || x >= container.width+container.x || y < container.y || y >= container.height+container.y {
		return nil, false
	}

	return container.grid[x-container.x][y-container.y], true
}

type TrackedTileContainer struct {
	b2dc                *Basic2DContainer
	gopherTileLocations map[int]*GopherWorldTile
	foodTileLocations   map[int]*GopherWorldTile
}

func NewTrackedTileContainer(x int, y int, width int, height int) TrackedTileContainer {
	b2dc := NewBasic2DContainer(x, y, width, height)

	return TrackedTileContainer{
		b2dc:                &b2dc,
		gopherTileLocations: make(map[int]*GopherWorldTile),
		foodTileLocations:   make(map[int]*GopherWorldTile),
	}
}

func (container *TrackedTileContainer) ConvertToTrackedTileCoordinates(x int, y int) (gridX int, gridY int) {
	return (x - container.b2dc.x), (y - container.b2dc.y)
}

type BasicGridContainer struct {
	containers [][]*TrackedTileContainer
	gridWidth  int
	gridHeight int
	width      int
	height     int
}

func NewBasicGridContainer(width int, height int, gridWidth int, gridHeight int) BasicGridContainer {

	numberOfGridsX, numberOfGridsY := width/gridWidth, height/gridHeight

	if numberOfGridsX*gridWidth < width {
		numberOfGridsX++
	}

	if numberOfGridsY*gridHeight < height {
		numberOfGridsY++
	}

	containers := make([][]*TrackedTileContainer, numberOfGridsX)

	for i := 0; i < numberOfGridsX; i++ {
		containers[i] = make([]*TrackedTileContainer, numberOfGridsY)

		for j := 0; j < numberOfGridsY; j++ {
			ttc := NewTrackedTileContainer(i*gridWidth,
				j*gridHeight,
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

func (container *BasicGridContainer) Tile(x int, y int) (*GopherWorldTile, bool) {

	if grid, ok := container.Grid(x, y); ok {
		if tile, ok := grid.b2dc.Tile(x, y); ok {
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

//GopherInserterAndRemover inserts and removes Gophers from an x and y position
type GopherInserterAndRemover interface {
	InsertGopher(x int, y int, gopher *Gopher) bool
	RemoveGopher(x int, y int) (*Gopher, bool)
}

//FoodInserterAndRemover inserts and removes Food from an x and y position
type FoodInserterAndRemover interface {
	InsertFood(x int, y int, food *Food) bool
	RemoveFood(x int, y int) (*Food, bool)
}

//GopherAndFoodInserterAndRemover inserts and removes Gophers and Food from an x and y position
type GopherAndFoodInserterAndRemover interface {
	GopherInserterAndRemover
	FoodInserterAndRemover
}

//InsertGopher Inserts the given gopher into the tileMap at the specified co-ordinate
func (container *Basic2DContainer) InsertGopher(x int, y int, gopher *Gopher) bool {

	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasGopher() {
			gopher.Position.SetXY(x, y)
			tile.SetGopher(gopher)
			return true
		}
	}

	return false

}

//InsertFood Inserts the given food into the tileMap at the specified co-ordinate
func (container *Basic2DContainer) InsertFood(x int, y int, food *Food) bool {

	if tile, ok := container.Tile(x, y); ok {
		if !tile.HasFood() {
			food.Position.SetXY(x, y)
			tile.SetFood(food)
			return true
		}
	}
	return false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (container *Basic2DContainer) RemoveGopher(x int, y int) (*Gopher, bool) {

	if tile, ok := container.Tile(x, y); ok {
		if tile.HasGopher() {
			gopher := tile.Gopher
			tile.ClearGopher()
			return gopher, true
		}
	}

	return nil, false
}

//RemoveFoodFromWorld Removes food from the given coordinates. Returns the food value.
func (container *Basic2DContainer) RemoveFood(x int, y int) (*Food, bool) {

	if tile, ok := container.Tile(x, y); ok {
		if tile.HasFood() {
			var food = tile.Food
			tile.ClearFood()
			return food, true
		}
	}

	return nil, false
}

func (container *Basic2DContainer) HasGopher(x int, y int) (*Gopher, bool) {

	if tile, ok := container.Tile(x, y); ok {
		if tile.HasGopher() {
			return tile.Gopher, true
		}
	}
	return nil, false
}

func (container *Basic2DContainer) HasFood(x int, y int) (*Food, bool) {

	if tile, ok := container.Tile(x, y); ok {
		if tile.HasFood() {
			return tile.Food, true
		}
	}
	return nil, false
}

func (container *TrackedTileContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if tile, ok := container.b2dc.Tile(x, y); ok {
		container.b2dc.InsertGopher(x, y, gopher)

		x, y = container.ConvertToTrackedTileCoordinates(x, y)
		container.gopherTileLocations[geometry.Hashcode(x, y)] = tile
		return true
	}
	return false
}

func (container *TrackedTileContainer) InsertFood(x int, y int, food *Food) bool {
	if tile, ok := container.b2dc.Tile(x, y); ok {
		container.b2dc.InsertFood(x, y, food)

		x, y = container.ConvertToTrackedTileCoordinates(x, y)
		container.foodTileLocations[geometry.Hashcode(x, y)] = tile
		return true
	}
	return false
}

func (container *TrackedTileContainer) RemoveGopher(x int, y int) (*Gopher, bool) {
	if gopher, ok := container.b2dc.RemoveGopher(x, y); ok {
		x, y = container.ConvertToTrackedTileCoordinates(x, y)
		delete(container.gopherTileLocations, geometry.Hashcode(x, y))
		return gopher, true
	}
	return nil, false
}

func (container *TrackedTileContainer) RemoveFood(x int, y int) (*Food, bool) {
	if food, ok := container.b2dc.RemoveFood(x, y); ok {
		x, y = container.ConvertToTrackedTileCoordinates(x, y)
		delete(container.foodTileLocations, geometry.Hashcode(x, y))
		return food, true
	}
	return nil, false
}

func (container *BasicGridContainer) InsertGopher(x int, y int, gopher *Gopher) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.InsertGopher(x, y, gopher)
	}
	return false
}

func (container *BasicGridContainer) InsertFood(x int, y int, food *Food) bool {
	if grid, ok := container.Grid(x, y); ok {
		return grid.InsertFood(x, y, food)
	}
	return false
}

func (container *BasicGridContainer) RemoveGopher(x int, y int) (*Gopher, bool) {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveGopher(x, y)
	}
	return nil, false
}

func (container *BasicGridContainer) RemoveFood(x int, y int) (*Food, bool) {
	if grid, ok := container.Grid(x, y); ok {
		return grid.RemoveFood(x, y)
	}
	return nil, false
}

func (container *BasicGridContainer) HasGopher(x int, y int) (*Gopher, bool) {
	if grid, ok := container.Grid(x, y); ok {
		return grid.b2dc.HasGopher(x, y)
	}
	return nil, false
}

func (container *BasicGridContainer) HasFood(x int, y int) (*Food, bool) {
	if grid, ok := container.Grid(x, y); ok {
		return grid.b2dc.HasFood(x, y)
	}
	return nil, false
}
