package world

import (
	"gopherlife/geometry"
	"gopherlife/timer"
	"math/rand"
	"time"
)

type SnakeWorldSettings struct {
	Dimensions
	SpeedReduction int
}

type SnakeWorld struct {
	ActionQueuer
	Container

	SnakeWorldSettings

	grid [][]*SnakeWorldTile

	SnakeHead  *SnakePart
	Direction  geometry.Direction
	IsGameOver bool
	Score      int

	FrameTimer timer.StopWatch
}

type SnakeWorldTile struct {
	geometry.Coordinates
	SnakePart *SnakePart
	SnakeWall *SnakeWall
	SnakeFood *SnakeFood
}

func NewSnakeTileGrid(x int, y int, width int, height int) [][]*SnakeWorldTile {
	grid := make([][]*SnakeWorldTile, width)

	for i := 0; i < width; i++ {
		grid[i] = make([]*SnakeWorldTile, height)

		for j := 0; j < height; j++ {
			tile := SnakeWorldTile{
				Coordinates: geometry.NewCoordinate(i, j),
			}
			grid[i][j] = &tile
		}
	}

	return grid
}

//NewEmptySnakeWorld Creates an Empty Snake Map (No Snake, No Wall)
func NewEmptySnakeWorld(settings SnakeWorldSettings) SnakeWorld {
	r := geometry.NewRectangle(0, 0, settings.Width, settings.Height)
	baq := NewFiniteActionQueue(1)

	SnakeWorld := SnakeWorld{
		grid:               NewSnakeTileGrid(0, 0, settings.Width, settings.Height),
		Container:          &r,
		ActionQueuer:       &baq,
		SnakeWorldSettings: settings,
		IsGameOver:         false,
	}

	return SnakeWorld
}

//NewSnakeWorld Creates a SnakeWorld with a Snake and Walls surrounding the edges
func NewSnakeWorld(settings SnakeWorldSettings) SnakeWorld {

	SnakeWorld := NewEmptySnakeWorld(settings)

	snakeHead := SnakePart{}
	startX, startY := settings.Width/2, settings.Height/2-2
	SnakeWorld.InsertSnakePart(startX, startY, &snakeHead)

	snakePartToAttachTo := &snakeHead

	for i := 0; i < 5; i++ {
		snakePartInStomach := SnakePart{}
		x, y := snakePartToAttachTo.GetX(), snakePartToAttachTo.GetY()-1
		SnakeWorld.InsertSnakePart(x, y, &snakePartInStomach)
		snakePartToAttachTo.AttachToBack(&snakePartInStomach)

		snakePartToAttachTo = &snakePartInStomach
	}

	for i := 0; i < settings.Width; i++ {
		SnakeWorld.InsertSnakeWall(i, 0, &SnakeWall{})
		SnakeWorld.InsertSnakeWall(i, settings.Height-1, &SnakeWall{})
	}

	for i := 0; i < settings.Height; i++ {
		SnakeWorld.InsertSnakeWall(0, i, &SnakeWall{})
		SnakeWorld.InsertSnakeWall(settings.Width-1, i, &SnakeWall{})
	}

	SnakeWorld.AddNewSnakeFoodToMap()

	SnakeWorld.SnakeHead = &snakeHead
	SnakeWorld.Direction = geometry.Up

	return SnakeWorld
}

func (sw *SnakeWorld) Update() bool {

	sw.FrameTimer.Start()

	sw.Process()

	if sw.IsGameOver {
		return false
	}

	if !sw.MoveSnake() {
		sw.IsGameOver = true
	}

	for sw.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*time.Duration(sw.SpeedReduction) {
	}

	return true
}

func (sw *SnakeWorld) MoveSnake() bool {

	currentSnakePart := sw.SnakeHead
	currentSnakePart.PassOnFood()

	nextX, nextY := sw.Direction.AddToPoint(currentSnakePart.GetX(), currentSnakePart.GetY())

	hasFood := sw.HasSnakeFood(nextX, nextY)

	//newPartPassedDownThisFrame := false

	for {

		prevX, prevY := currentSnakePart.GetX(), currentSnakePart.GetY()
		sw.RemoveSnakePart(prevX, prevY)

		inserted := sw.InsertSnakePart(nextX, nextY, currentSnakePart)

		if !inserted {
			sw.InsertSnakePart(prevX, prevY, currentSnakePart)
			return false
		}

		if hasFood && currentSnakePart.snakePartInFront == nil {
			if sw.RemoveSnakeFood(nextX, nextY) {
				currentSnakePart.snakePartInStomach = &SnakePart{}
				hasFood = false
				sw.AddNewSnakeFoodToMap()
				sw.Score += 10
			}
		}
		nextX, nextY = prevX, prevY

		if currentSnakePart.snakePartBehind == nil {
			break
		}

		currentSnakePart = currentSnakePart.snakePartBehind
	}

	return true

}

func (smt *SnakeWorld) ChangeDirection(d geometry.Direction) {

	setDirection := func(d geometry.Direction) {
		smt.Add(func() {
			smt.Direction = d
		})
	}

	switch d {
	case geometry.Left:
		fallthrough
	case geometry.Right:
		if smt.Direction == geometry.Up || smt.Direction == geometry.Down {
			setDirection(d)
		}
	case geometry.Up:
		fallthrough
	case geometry.Down:
		if smt.Direction == geometry.Left || smt.Direction == geometry.Right {
			setDirection(d)
		}
	}

}

func (sw *SnakeWorld) AddNewSnakeFoodToMap() bool {

	xrange, yrange := rand.Perm(sw.Width), rand.Perm(sw.Height)

	for i := 0; i < sw.Width; i++ {
		for j := 0; j < sw.Height; j++ {
			newX, newY := xrange[i], yrange[j]
			if sw.InsertSnakeFood(newX, newY, &SnakeFood{}) {
				return true
			}
		}
	}

	return false
}

func (smt *SnakeWorld) Tile(x int, y int) (*SnakeWorldTile, bool) {
	if smt.Contains(x, y) {
		return smt.grid[x][y], true
	}
	return nil, false
}

func (smt *SnakeWorld) InsertSnakePart(x int, y int, sp *SnakePart) bool {
	if tile, ok := smt.Tile(x, y); ok {
		return tile.InsertSnakePart(sp)
	}
	return false
}

func (smt *SnakeWorld) InsertSnakeFood(x int, y int, sf *SnakeFood) bool {
	if tile, ok := smt.Tile(x, y); ok {
		return tile.InsertSnakeFood(sf)
	}
	return false
}

func (smt *SnakeWorld) InsertSnakeWall(x int, y int, w *SnakeWall) bool {
	if tile, ok := smt.Tile(x, y); ok {
		return tile.InsertWall(w)
	}
	return false
}

func (smt *SnakeWorld) RemoveSnakePart(x int, y int) bool {
	if tile, ok := smt.Tile(x, y); ok {
		tile.RemoveSnakePart()
		return true
	}
	return false
}

func (smt *SnakeWorld) RemoveSnakeFood(x int, y int) bool {
	if tile, ok := smt.Tile(x, y); ok {
		tile.RemoveSnakeFood()
		return true
	}
	return false
}

func (smt *SnakeWorld) HasSnakeFood(x int, y int) bool {
	if tile, ok := smt.Tile(x, y); ok {
		return tile.SnakeFood != nil
	}
	return false
}

func (smt *SnakeWorldTile) InsertWall(w *SnakeWall) bool {
	if smt.SnakeWall == nil && smt.SnakeFood == nil && smt.SnakePart == nil {
		w.SetXY(smt.GetX(), smt.GetY())
		smt.SnakeWall = w
		return true
	}
	return false
}

func (smt *SnakeWorldTile) InsertSnakeFood(sf *SnakeFood) bool {
	if smt.SnakeFood == nil && smt.SnakeWall == nil && smt.SnakePart == nil {
		sf.SetXY(smt.GetX(), smt.GetY())
		smt.SnakeFood = sf
		return true
	}
	return false
}

func (smt *SnakeWorldTile) InsertSnakePart(sp *SnakePart) bool {
	if smt.SnakePart == nil && smt.SnakeWall == nil {
		sp.SetXY(smt.GetX(), smt.GetY())
		smt.SnakePart = sp
		return true
	}

	return false
}

func (smt *SnakeWorldTile) RemoveSnakePart() {
	smt.SnakePart = nil
}

func (smt *SnakeWorldTile) RemoveSnakeFood() {
	smt.SnakeFood = nil
}

type SnakePart struct {
	geometry.Coordinates
	snakePartInFront   *SnakePart
	snakePartBehind    *SnakePart
	snakePartInStomach *SnakePart
}

//AttachToBack Adds a SnakePart to the back of another SnakePart.
func (sp *SnakePart) AttachToBack(partToAttach *SnakePart) {
	sp.snakePartBehind = partToAttach
	partToAttach.snakePartInFront = sp
}

//PassOnFood Walks through a Snake and moves any food in the stomach to the snack part behind
func (sp *SnakePart) PassOnFood() {

	if sp.snakePartBehind != nil {
		sp.snakePartBehind.PassOnFood()
	}

	if sp.HasPartInStomach() {

		if sp.snakePartBehind == nil {
			sp.AttachToBack(sp.snakePartInStomach)
			sp.snakePartInStomach = nil
		} else {
			sp.snakePartBehind.snakePartInStomach = sp.snakePartInStomach
			sp.snakePartInStomach = nil
		}
	}
}

//HasPartInStomach Return true is there is a SnakePart inside the Stomach
func (sp *SnakePart) HasPartInStomach() bool {
	return sp.snakePartInStomach != nil
}

//SnakeWall Used in the Snake Game has an X and Y Coordinate
type SnakeWall struct {
	geometry.Coordinates
}

//SnakeFood Used in the Snake Game has an X and Y Coordinate
type SnakeFood struct {
	geometry.Coordinates
}
