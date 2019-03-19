package geometry

import (
	"math"
)

type Spiral struct {
	width  int
	height int

	x float64
	y float64

	bigX float64
	bigY float64

	dx float64
	dy float64

	i    int
	maxI int
}

//NewSpiral Creates a new spiral that can be used to step through coordinates of the given width and height. Coordinates always start from (0,0)
func NewSpiral(width int, height int) Spiral {

	return Spiral{
		width:  width,
		height: height,
		bigX:   float64(width) / 2,
		bigY:   float64(height) / 2,
		x:      0,
		y:      0,
		dx:     0,
		dy:     -1,
		i:      0,
		maxI:   int(math.Pow(math.Max(float64(width), float64(height)), 2)),
	}

}

//Next Gets the next step in the spiral path. If there is no next step, returns false
func (s *Spiral) Next() (Coordinates, bool) {

	var c = Coordinates{}
	var nextPositionFound = false

	for s.i < s.maxI {

		if ((-s.bigX < s.x) && (s.x <= s.bigX)) &&
			((-s.bigY < s.y) && (s.y <= s.bigY)) {
			c.X, c.Y = int(s.x), int(s.y)
			nextPositionFound = true
		}

		if (s.x == s.y) || (s.x < 0 && s.x == -s.y) || (s.x > 0 && s.x == 1-s.y) {
			s.dx, s.dy = -s.dy, s.dx //Change direction of the spiral
		}

		s.x, s.y = s.x+s.dx, s.y+s.dy

		s.i = s.i + 1

		if nextPositionFound {
			break
		}

	}

	if !nextPositionFound {
		return c, false
	}

	return c, true
}
