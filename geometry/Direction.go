package geometry

type Direction int

const (
	Up    Direction = 1
	Left  Direction = 2
	Down  Direction = 3
	Right Direction = 4
)

//TurnClockWise90 Returns a Direction 90 degrees clockwise from the given Direction
func (d Direction) TurnClockWise90() Direction {

	switch d {

	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up

	}
	panic("Direction not covered")
}

//TurnAntiClockWise90 Returns a Direction 90 degrees anticlockwise from the given Direction
func (d Direction) TurnAntiClockWise90() Direction {

	switch d {
	case Up:
		return Left
	case Left:
		return Down
	case Down:
		return Right
	case Right:
		return Up
	}
	panic("Direction not covered")
}

//AddToPoint Adds the Direction to the X and Y value. The amount added is of distance 1
func (d Direction) AddToPoint(x int, y int) (int, int) {

	switch d {
	case Up:
		return x, y + 1
	case Right:
		return x + 1, y
	case Down:
		return x, y - 1
	case Left:
		return x - 1, y
	}
	panic("Direction not covered")
}
