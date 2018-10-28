package math

func FindNextStep(start Coordinates, end Coordinates) (x int, y int) {

	diffX := start.GetX() - end.GetX()
	diffY := start.GetY() - end.GetY()

	moveX := 0
	moveY := 0

	if diffX > 0 {
		moveX = -1
	} else if diffX < 0 {
		moveX = 1
	}

	if diffY > 0 {
		moveY = -1
	} else if diffY < 0 {
		moveY = 1
	}

	return moveX, moveY

}
