package world

import (
	"gopherlife/calc"
	"image/color"
	"sync"
)

type CollisionMap struct {
	grid [][]*ColliderTile
	Containable
	QueueableActions
	*sync.WaitGroup

	ActiveColliders chan *Collider
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
		//fmt.Println("Success!")
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

func NewCollisionMap(statistics Statistics) CollisionMap {

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)
	var wg sync.WaitGroup

	rect := NewRectangle(0, 0, statistics.Width, statistics.Height)

	collsionMap := CollisionMap{
		QueueableActions: &qa,
		WaitGroup:        &wg,
		Containable:      &rect,
		ActiveColliders:  make(chan *Collider, statistics.MaximumNumberOfGophers*2),
	}

	collsionMap.grid = make([][]*ColliderTile, statistics.Width)

	for i := 0; i < statistics.Width; i++ {
		collsionMap.grid[i] = make([]*ColliderTile, statistics.Height)

		for j := 0; j < statistics.Height; j++ {
			tile := ColliderTile{
				Position: Position{
					X: i,
					Y: j,
				},
			}
			collsionMap.grid[i][j] = &tile
		}
	}

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		statistics.Width, statistics.Height)

	count := 0

	for i := 0; i < statistics.NumberOfGophers; i++ {

		pos := keys[count]
		var c = Collider{
			velX:                 1,
			ColliderWorldActions: &collsionMap,
		}

		collsionMap.InsertCollider(pos.GetX(), pos.GetX(), &c)
		collsionMap.ActiveColliders <- &c

		count++
	}

	return collsionMap
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
		//fmt.Println(newX, newY)
		newTile := collisionMap.grid[newX][newY]

		if !newTile.HasCollider() {

			collisionMap.QueueableActions.Add(func() {
				if newTile.Insert(c) {
					oldTile.Clear()
				}
			})
		}
		return true
	}

	return false

}

type Collider struct {
	Position
	ColliderWorldActions
	color.RGBA

	velX int
	velY int
}

func (collider *Collider) Update() {

	if !collider.MoveCollider(collider.velX, collider.velY, collider) {
		//fmt.Println("Ow")
		collider.ChangeDirection()
	}
}

func (collider *Collider) ChangeDirection() {
	collider.velX = -1 * collider.velX
}
