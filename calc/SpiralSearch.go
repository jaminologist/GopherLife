package calc

import (
	"fmt"
	"math"
)

func SpiralCoordinates(n int) {

	x, y, d, m := 0, 0, 1, 1

	for i := 0; i < n; i++ {
		for (2 * x * d) < m {
			fmt.Println(x, y)
			x = x + d
		}
		for (2 * y * d) < m {
			fmt.Println(x, y)
			y = y + d
		}

		d = -1 * d
		m = m + 1
	}

}

type Spiral struct {
	width  int
	height int

	x int
	y int

	dx int
	dy int

	i int
}

func NewSpiral(width int, height int) Spiral {

	return Spiral{
		width:  width,
		height: height,
		x:      0,
		y:      0,
		dx:     0,
		dy:     -1,
		i:      0,
	}

}

func (s *Spiral) Next() (Coordinates, bool) {

	X := s.width
	Y := s.height

	var coords = Coordinates{}

	for s.i < int(math.Pow(float64(Max(X, Y)), 2)) {

		if (-X/2 < s.x) && (s.x <= X/2) && (-Y/2 < s.y) && (s.y <= Y/2) {
			coords.X, coords.Y = s.x, s.y
		}

		if (s.x == s.y) || (s.x < 0 && s.x == -s.y) || (s.x > 0 && s.x == 1-s.y) {
			s.dx, s.dy = -s.dy, s.dx
		}

		s.x, s.y = s.x+s.dx, s.y+s.dy

		s.i = s.i + 1

		if coords != (Coordinates{}) {
			break
		}

	}

	if coords == (Coordinates{}) {
		return coords, false
	}

	return coords, true
}

func SpiralCoordinates2(n int) {

	X, Y := 20, 5

	x, y := 0, 0

	dx := 0
	dy := -1

	for i := 0; i < int(math.Pow(float64(Max(X, Y)), 2)); i++ {

		if (-X/2 < x) && (x <= X/2) && (-Y/2 < y) && (y <= Y/2) {
			fmt.Println(x, y)
		}

		if (x == y) || (x < 0 && x == -y) || (x > 0 && x == 1-y) {
			dx, dy = -dy, dx
		}

		x, y = x+dx, y+dy
	}

	//return x, y

}
