package world

import (
	"gopherlife/calc"
	"math/rand"
	"time"
)

type SnakeMap struct {
	ActionQueuer
	Containable
	Dimensions

	grid [][]*SnakeMapTile

	SnakeHead  *SnakePart
	Direction  Direction
	IsGameOver bool
	Score      int

	FrameTimer calc.StopWatch
	FrameSpeed time.Duration
}

type SnakeMapTile struct {
	Position
	SnakePart *SnakePart
	Wall      *Wall
	SnakeFood *SnakeFood
}

type Direction int

const (
	Left  Direction = 1
	Right Direction = 2
	Up    Direction = 3
	Down  Direction = 4
)

func NewSnakeMap(d Dimensions, speed int) SnakeMap {

	r := NewRectangle(0, 0, d.Width, d.Height)
	grid := make([][]*SnakeMapTile, d.Width)

	for i := 0; i < d.Width; i++ {
		grid[i] = make([]*SnakeMapTile, d.Height)

		for j := 0; j < d.Height; j++ {
			tile := SnakeMapTile{
				Position: Position{
					X: i,
					Y: j,
				},
			}
			grid[i][j] = &tile
		}
	}

	baq := NewBasicActionQueue(1)

	snakeMap := SnakeMap{
		grid:         grid,
		Containable:  &r,
		ActionQueuer: &baq,
		Dimensions:   d,
		IsGameOver:   false,
		FrameSpeed:   time.Duration(speed),
	}

	snakeHead := SnakePart{}
	startX, startY := d.Width/2, d.Height/2-5
	snakeMap.InsertSnakePart(startX, startY, &snakeHead)

	snakePartToAttachTo := &snakeHead

	for i := 0; i < 5; i++ {
		snakePartInStomach := SnakePart{}
		x, y := snakePartToAttachTo.GetX(), snakePartToAttachTo.GetY()-1
		snakeMap.InsertSnakePart(x, y, &snakePartInStomach)
		snakePartToAttachTo.Attach(&snakePartInStomach)

		snakePartToAttachTo = &snakePartInStomach
	}

	for i := 0; i < d.Width; i++ {
		snakeMap.InsertSnakeWall(i, 0, &Wall{})
		snakeMap.InsertSnakeWall(i, d.Height-1, &Wall{})
	}

	for i := 0; i < d.Height; i++ {
		snakeMap.InsertSnakeWall(0, i, &Wall{})
		snakeMap.InsertSnakeWall(d.Width-1, i, &Wall{})
	}

	snakeMap.AddNewSnakeFoodToMap(0, 0)

	snakeMap.SnakeHead = &snakeHead
	snakeMap.Direction = Up

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

	x, y := 0, 0

	switch sm.Direction {
	case Left:
		x = -1
	case Right:
		x = 1
	case Up:
		y = 1
	case Down:
		y = -1
	}

	currentSnakePart := sm.SnakeHead

	nextX, nextY := currentSnakePart.GetX()+x, currentSnakePart.GetY()+y

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
				currentSnakePart.Attach(currentSnakePart.snakePartInStomach)
				currentSnakePart.snakePartInStomach = nil
			}

			break
		}

		currentSnakePart = currentSnakePart.snakePartBehind
	}

	return true

}

func (smt *SnakeMap) ChangeDirection(d Direction) {

	setDirection := func(d Direction) {
		smt.Add(func() {
			smt.Direction = d
		})
	}

	switch d {
	case Left:
		fallthrough
	case Right:
		if smt.Direction == Up || smt.Direction == Down {
			setDirection(d)
		}
	case Up:
		fallthrough
	case Down:
		if smt.Direction == Left || smt.Direction == Right {
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

func (smt *SnakeMap) InsertSnakeWall(x int, y int, w *Wall) bool {
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

func (smt *SnakeMapTile) InsertWall(w *Wall) bool {
	if smt.Wall == nil && smt.SnakeFood == nil && smt.SnakePart == nil {
		w.SetPosition(smt.GetX(), smt.GetY())
		smt.Wall = w
		return true
	}
	return false
}

func (smt *SnakeMapTile) InsertSnakeFood(sf *SnakeFood) bool {
	if smt.SnakeFood == nil && smt.Wall == nil && smt.SnakePart == nil {
		sf.SetPosition(smt.GetX(), smt.GetY())
		smt.SnakeFood = sf
		return true
	}
	return false
}

func (smt *SnakeMapTile) InsertSnakePart(sp *SnakePart) bool {
	if smt.SnakePart == nil && smt.Wall == nil {
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
	Position
	snakePartInFront   *SnakePart
	snakePartBehind    *SnakePart
	snakePartInStomach *SnakePart
}

func (sp *SnakePart) Attach(partToAttach *SnakePart) {
	sp.snakePartBehind = partToAttach
	partToAttach.snakePartInFront = sp
}

type Wall struct {
	Position
}

type SnakeFood struct {
	Position
}
