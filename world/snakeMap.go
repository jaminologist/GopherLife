package world

import "sync"

type SnakeMap struct {
	grid [][]*SnakeMapTile
	Containable
	ActionQueuer
	*sync.WaitGroup

	ActiveColliders chan *Collider

	IsDiagonal bool
}

type SnakeMapTile struct {
	Position
	SnakePart
	*Wall
}

func (smt *SnakeMapTile) InsertWall(w *Wall) bool {
	if smt.Wall == nil {
		w.SetPosition(smt.GetX(), smt.GetY())
		smt.Wall = w
		return true
	}
	return false
}

type SnakePart struct {
	Position
	prev *SnakePart
}

type Wall struct {
	Position
}
