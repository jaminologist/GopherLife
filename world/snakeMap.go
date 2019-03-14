package world

import (
	"gopherlife/geometry"
	"gopherlife/timer"
	"math/rand"
	"time"
)

type SnakeMap struct {
	ActionQueuer
	Container
	Dimensions

	grid [][]*SnakeMapTile

	SnakeHead  *SnakePart
	Direction  geometry.Direction
	IsGameOver bool
	Score      int

	FrameTimer timer.StopWatch
	FrameSpeed time.Duration
}

type SnakeMapTile struct {
	geometry.Coordinates
	SnakePart *SnakePart
	SnakeWall *SnakeWall
	SnakeFood *SnakeFood
}

func NewSnakeTileGrid(x int, y int, width int, height int) [][]*SnakeMapTile {
	grid := make([][]*SnakeMapTile, width)

	for i := 0; i < width; i++ {
		grid[i] = make([]*SnakeMapTile, height)

		for j := 0; j < height; j++ {
			tile := SnakeMapTile{
				Coordinates: geometry.NewCoordinate(i, j),
			}
			grid[i][j] = &tile
		}
	}

	return grid
}

//NewEmptySnakeMap Creates an Empty Snake Map (No Snake, No Wall)
func NewEmptySnakeMap(d Dimensions, speed int) SnakeMap {
	r := geometry.NewRectangle(0, 0, d.Width, d.Height)
	baq := NewBasicActionQueue(1)

	snakeMap := SnakeMap{
		grid:         NewSnakeTileGrid(0, 0, d.Width, d.Height),
		Container:    &r,
		ActionQueuer: &baq,
		Dimensions:   d,
		IsGameOver:   false,
		FrameSpeed:   time.Duration(speed),
	}

	return snakeMap
}

//NewSnakeMap Creates a SnakeMap with a Snake and Walls surrounding the edges
func NewSnakeMap(d Dimensions, speed int) SnakeMap {

	snakeMap := NewEmptySnakeMap(d, speed)

	snakeHead := SnakePart{}
	startX, startY := d.Width/2, d.Height/2-5
	snakeMap.InsertSnakePart(startX, startY, &snakeHead)

	snakePartToAttachTo := &snakeHead

	for i := 0; i < 5; i++ {
		snakePartInStomach := SnakePart{}
		x, y := snakePartToAttachTo.GetX(), snakePartToAttachTo.GetY()-1
		snakeMap.InsertSnakePart(x, y, &snakePartInStomach)
		snakePartToAttachTo.AttachToBack(&snakePartInStomach)

		snakePartToAttachTo = &snakePartInStomach
	}

	for i := 0; i < d.Width; i++ {
		snakeMap.InsertSnakeWall(i, 0, &SnakeWall{})
		snakeMap.InsertSnakeWall(i, d.Height-1, &SnakeWall{})
	}

	for i := 0; i < d.Height; i++ {
		snakeMap.InsertSnakeWall(0, i, &SnakeWall{})
		snakeMap.InsertSnakeWall(d.Width-1, i, &SnakeWall{})
	}

	snakeMap.AddNewSnakeFoodToMap(0, 0)

	snakeMap.SnakeHead = &snakeHead
	snakeMap.Direction = geometry.Up

	return snakeMap
}

func (sm *SnakeMap) Update() bool {

	sm.FrameTimer.Start()

	sm.Process()

	if sm.IsGameOver {
		return false
	}

	if !sm.MoveSnake() {
		sm.IsGameOver = true
	}

	for sm.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*sm.FrameSpeed {
	}

	return true
}

func (sm *SnakeMap) MoveSnake() bool {

	currentSnakePart := sm.SnakeHead

	nextX, nextY := sm.Direction.AddToPoint(currentSnakePart.GetX(), currentSnakePart.GetY())

	hasFood := sm.HasSnakeFood(nextX, nextY)

	newPartPassedDownThisFrame := false

	for {

		prevX, prevY := currentSnakePart.GetX(), currentSnakePart.GetY()
		sm.RemoveSnakePart(prevX, prevY)

		inserted := sm.InsertSnakePart(nextX, nextY, currentSnakePart)

		if !inserted {
			sm.InsertSnakePart(prevX, prevY, currentSnakePart)
			return false
		}

		if hasFood && currentSnakePart.snakePartInFront == nil {
			if sm.RemoveSnakeFood(nextX, nextY) {
				currentSnakePart.snakePartInStomach = &SnakePart{}
				hasFood = false
				sm.AddNewSnakeFoodToMap(nextX, nextY)
				sm.Score += 10
			}
		}

		//Fun litte bug, due to the boolean if two things are swallowed at the same time on one will stay still on the screen.
		//You can decide whether to fix this or not
		if currentSnakePart.snakePartInStomach != nil && currentSnakePart.snakePartBehind != nil && !newPartPassedDownThisFrame {
			currentSnakePart.snakePartBehind.snakePartInStomach = currentSnakePart.snakePartInStomach
			currentSnakePart.snakePartInStomach = nil
			newPartPassedDownThisFrame = true
		}

		nextX, nextY = prevX, prevY

		if currentSnakePart.snakePartBehind == nil {

			if currentSnakePart.snakePartInStomach != nil {
				currentSnakePart.AttachToBack(currentSnakePart.snakePartInStomach)
				currentSnakePart.snakePartInStomach = nil
			}

			break
		}

		currentSnakePart = currentSnakePart.snakePartBehind
	}

	return true

}

func (smt *SnakeMap) ChangeDirection(d geometry.Direction) {

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

func (sm *SnakeMap) AddNewSnakeFoodToMap(oldX int, oldY int) bool {

	xrange, yrange := rand.Perm(sm.Width), rand.Perm(sm.Height)

	for i := 0; i < sm.Width; i++ {
		for j := 0; j < sm.Height; j++ {
			newX, newY := xrange[i], yrange[j]
			if sm.InsertSnakeFood(newX, newY, &SnakeFood{}) {
				return true
			}
		}
	}

	return false
}

func (smt *SnakeMap) Tile(x int, y int) (*SnakeMapTile, bool) {
	if smt.Contains(x, y) {
		return smt.grid[x][y], true
	}
	return nil, false
}

func (smt *SnakeMap) InsertSnakePart(x int, y int, sp *SnakePart) bool {
	if smt.Contains(x, y) {
		return smt.grid[x][y].InsertSnakePart(sp)
	}
	return false
}

func (smt *SnakeMap) InsertSnakeFood(x int, y int, sf *SnakeFood) bool {
	if smt.Contains(x, y) {
		return smt.grid[x][y].InsertSnakeFood(sf)
	}
	return false
}

func (smt *SnakeMap) InsertSnakeWall(x int, y int, w *SnakeWall) bool {
	if smt.Contains(x, y) {
		return smt.grid[x][y].InsertWall(w)
	}
	return false
}

func (smt *SnakeMap) RemoveSnakePart(x int, y int) bool {
	if smt.Contains(x, y) {
		smt.grid[x][y].RemoveSnakePart()
		return true
	}
	return false
}

func (smt *SnakeMap) RemoveSnakeFood(x int, y int) bool {
	if smt.Contains(x, y) {
		smt.grid[x][y].RemoveSnakeFood()
		return true
	}
	return false
}

func (smt *SnakeMap) HasSnakeFood(x int, y int) bool {
	if tile, ok := smt.Tile(x, y); ok {
		return tile.SnakeFood != nil
	}
	return false
}

func (smt *SnakeMapTile) InsertWall(w *SnakeWall) bool {
	if smt.SnakeWall == nil && smt.SnakeFood == nil && smt.SnakePart == nil {
		w.SetPosition(smt.GetX(), smt.GetY())
		smt.SnakeWall = w
		return true
	}
	return false
}

func (smt *SnakeMapTile) InsertSnakeFood(sf *SnakeFood) bool {
	if smt.SnakeFood == nil && smt.SnakeWall == nil && smt.SnakePart == nil {
		sf.SetPosition(smt.GetX(), smt.GetY())
		smt.SnakeFood = sf
		return true
	}
	return false
}

func (smt *SnakeMapTile) InsertSnakePart(sp *SnakePart) bool {
	if smt.SnakePart == nil && smt.SnakeWall == nil {
		sp.SetPosition(smt.GetX(), smt.GetY())
		smt.SnakePart = sp
		return true
	}

	return false
}

func (smt *SnakeMapTile) RemoveSnakePart() {
	smt.SnakePart = nil
}

func (smt *SnakeMapTile) RemoveSnakeFood() {
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
