package world

import (
	"gopherlife/colors"
	"gopherlife/geometry"
	"image/color"
	"math/rand"
	"sync"
)

var collisionMapSpeed = 1

var colliderColors = []color.RGBA{colors.Red, colors.Blue, colors.Cyan, colors.Pink, colors.Yellow, colors.White}

type CollisionMapSettings struct {
	Dimensions
	Population
	IsDiagonal bool
}

type CollisionMap struct {
	grid [][]*ColliderTile
	Container
	ActionQueuer
	CollisionMapSettings
	*sync.WaitGroup

	ActiveColliders chan *Collider
}

type ColliderTile struct {
	geometry.Coordinates
	c *Collider
}

//Insert Adds the Collider to the Tile, if Empty
func (tile *ColliderTile) Insert(c *Collider) bool {
	if tile.c == nil {
		c.X = tile.X
		c.Y = tile.Y
		tile.c = c
		return true
	}
	return false
}

//HasCollider Return true if Collider exists
func (tile *ColliderTile) HasCollider() bool {
	return tile.c != nil
}

//Clear Empties the ColliderTile
func (tile *ColliderTile) Clear() {
	tile.c = nil
}

//NewEmptyCollisionMap Creates an Empty Collision Map
func NewEmptyCollisionMap(settings CollisionMapSettings) CollisionMap {

	qa := NewBasicActionQueue(settings.InitialPopulation * 2)
	var wg sync.WaitGroup

	rect := geometry.NewRectangle(0, 0, settings.Width, settings.Height)

	collisionMap := CollisionMap{
		ActionQueuer:         &qa,
		WaitGroup:            &wg,
		Container:            &rect,
		ActiveColliders:      make(chan *Collider, settings.InitialPopulation*2),
		CollisionMapSettings: settings,
	}

	collisionMap.grid = make([][]*ColliderTile, settings.Width)

	for i := 0; i < settings.Width; i++ {
		collisionMap.grid[i] = make([]*ColliderTile, settings.Height)

		for j := 0; j < settings.Height; j++ {
			tile := ColliderTile{
				Coordinates: geometry.Coordinates{
					X: i,
					Y: j,
				},
			}
			collisionMap.grid[i][j] = &tile
		}
	}

	return collisionMap
}

//NewCollisionMap Creates a Populated Collision Map
func NewCollisionMap(settings CollisionMapSettings) CollisionMap {

	collisionMap := NewEmptyCollisionMap(settings)

	keys := geometry.GenerateRandomizedCoordinateArray(0, 0,
		settings.Width, settings.Height)

	count := 0

	for i := 0; i < settings.InitialPopulation; i++ {

		var velX, velY int

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

//Update all Active Colliders inside the Map
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

//InsertCollider Sets x and y of Collider and places it into map
func (collisionMap *CollisionMap) InsertCollider(x int, y int, c *Collider) bool {

	if collisionMap.Contains(x, y) {
		return collisionMap.grid[x][y].Insert(c)
	}

	return false
}

//HasCollider Checks if a colliders exists at X and Y and returns the Collider
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

			collisionMap.ActionQueuer.Add(func() {
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
	geometry.Coordinates
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

		if geometry.Abs(collider.velX) > 0 {
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

//ChangeColor Changes a Collider's color and increments ColorSelection
func (collider *Collider) ChangeColor() {

	collider.colorSelection++

	if collider.colorSelection > len(colliderColors)-1 || collider.colorSelection < 0 {
		collider.colorSelection = 0
	}

	collider.Color = colliderColors[collider.colorSelection]
}
