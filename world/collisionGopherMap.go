package world

import (
	"gopherlife/calc"
	"image/color"
	"math/rand"
	"sync"
)

var collisionMapSpeed = 1

type CollisionMap struct {
	grid [][]*ColliderTile
	Containable
	QueueableActions
	*sync.WaitGroup

	ActiveColliders chan *Collider

	IsDiagonal bool
}

type ColliderTile struct {
	Position
	c *Collider
}

func (tile *ColliderTile) Insert(c *Collider) bool {
	if tile.c == nil {
		c.X = tile.X
		c.Y = tile.Y
		tile.c = c
		return true
	}
	return false
}

func (tile *ColliderTile) HasCollider() bool {
	return tile.c != nil
}

func (tile *ColliderTile) Clear() {
	tile.c = nil
}

func NewCollisionMap(statistics Statistics, isDiagonal bool) CollisionMap {

	qa := NewBasicActionQueue(statistics.NumberOfGophers * 2)
	var wg sync.WaitGroup

	rect := NewRectangle(0, 0, statistics.Width, statistics.Height)

	collisionMap := CollisionMap{
		QueueableActions: &qa,
		WaitGroup:        &wg,
		Containable:      &rect,
		ActiveColliders:  make(chan *Collider, statistics.NumberOfGophers*2),
		IsDiagonal:       isDiagonal,
	}

	collisionMap.grid = make([][]*ColliderTile, statistics.Width)

	for i := 0; i < statistics.Width; i++ {
		collisionMap.grid[i] = make([]*ColliderTile, statistics.Height)

		for j := 0; j < statistics.Height; j++ {
			tile := ColliderTile{
				Position: Position{
					X: i,
					Y: j,
				},
			}
			collisionMap.grid[i][j] = &tile
		}
	}

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		statistics.Width, statistics.Height)

	count := 0

	for i := 0; i < statistics.NumberOfGophers; i++ {

		var velX, velY int

		//collisionMap.IsDiagonal = !collisionMap.IsDiagonal

		if collisionMap.IsDiagonal {
			velX = getNegativeOrPositiveSpeed(collisionMapSpeed)
			velY = getNegativeOrPositiveSpeed(collisionMapSpeed)
		} else {
			if rand.Intn(2) == 0 {
				velX = getNegativeOrPositiveSpeed(collisionMapSpeed)
			} else {
				velY = getNegativeOrPositiveSpeed(collisionMapSpeed)
			}

		}

		pos := keys[count]
		var c = Collider{
			velX:                 velX,
			velY:                 velY,
			ColliderWorldActions: &collisionMap,
			IsDiagonal:           collisionMap.IsDiagonal,
		}

		collisionMap.InsertCollider(pos.GetX(), pos.GetY(), &c)
		collisionMap.ActiveColliders <- &c

		count++
	}

	return collisionMap
}

func (collisionMap *CollisionMap) Update() bool {

	numColliders := len(collisionMap.ActiveColliders)

	for i := 0; i < numColliders; i++ {
		collider := <-collisionMap.ActiveColliders
		collisionMap.WaitGroup.Add(1)
		go func() {
			collider.Update()
			collisionMap.ActiveColliders <- collider
			collisionMap.WaitGroup.Done()
		}()
	}

	collisionMap.WaitGroup.Wait()
	collisionMap.Process()

	return true

}

func getNegativeOrPositiveSpeed(speed int) int {
	if rand.Intn(2) == 0 {
		return speed
	} else {
		return -speed
	}
}

type ColliderWorldActions interface {
	MoveCollider(moveX int, moveY int, c *Collider) bool
}

func (collisionMap *CollisionMap) InsertCollider(x int, y int, c *Collider) bool {

	if collisionMap.Contains(x, y) {
		return collisionMap.grid[x][y].Insert(c)
	}

	return false
}

func (collisionMap *CollisionMap) HasCollider(x int, y int) (*Collider, bool) {

	if collisionMap.Contains(x, y) {
		tile := collisionMap.grid[x][y]

		if tile.HasCollider() {
			return tile.c, true
		}
	}

	return nil, false
}

func (collisionMap *CollisionMap) MoveCollider(moveX int, moveY int, c *Collider) bool {

	newX, newY := c.X+moveX, c.Y+moveY

	if collisionMap.Contains(newX, newY) {

		oldTile := collisionMap.grid[c.X][c.Y]
		newTile := collisionMap.grid[newX][newY]

		if !newTile.HasCollider() {

			collisionMap.QueueableActions.Add(func() {
				if newTile.Insert(c) {
					oldTile.Clear()
				}
			})
			return true
		}
	}

	return false

}

type Collider struct {
	Position
	ColliderWorldActions
	Color color.RGBA

	IsDiagonal bool

	velX           int
	velY           int
	colorSelection int
}

func (collider *Collider) Update() {

	if !collider.MoveCollider(collider.velX, collider.velY, collider) {
		collider.ChangeDirection()
		collider.ChangeColor()
		collider.MoveCollider(collider.velX, collider.velY, collider)
	}
}

func (collider *Collider) ChangeDirection() {

	if collider.IsDiagonal {

		i := rand.Intn(3)

		if i == 0 || i == 2 {
			collider.velX = -1 * collider.velX
		}

		if i == 1 || i == 2 {
			collider.velY = -1 * collider.velY
		}

	} else {

		if calc.Abs(collider.velX) > 0 {
			i := rand.Intn(2)

			if i == 0 {
				collider.velX = -1 * collider.velX
			} else {
				collider.velX = 0
				collider.velY = getNegativeOrPositiveSpeed(collisionMapSpeed)
			}

		} else {

			i := rand.Intn(2)

			if i == 0 {
				collider.velY = -1 * collider.velY
			} else {
				collider.velY = 0
				collider.velX = getNegativeOrPositiveSpeed(collisionMapSpeed)
			}

		}

	}
}

func (collider *Collider) ChangeColor() {

	switch collider.colorSelection {
	case 0:
		collider.Color = color.RGBA{
			255, 0, 0, 1,
		}
	case 1:
		collider.Color = color.RGBA{
			0, 255, 0, 1,
		}
	case 2:
		collider.Color = color.RGBA{
			0, 204, 255, 1,
		}
	case 3:
		collider.Color = color.RGBA{
			255, 0, 255, 1,
		}
	case 4:
		collider.Color = color.RGBA{
			255, 255, 0, 1,
		}
	case 5:
		collider.Color = color.RGBA{
			255, 255, 255, 1,
		}
	default:
		collider.colorSelection = -1
	}

	collider.colorSelection++
}
