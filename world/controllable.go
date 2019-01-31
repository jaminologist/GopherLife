package world

type Controllable interface {
	Click(x int, y int)
	KeyPress(key int)
}
