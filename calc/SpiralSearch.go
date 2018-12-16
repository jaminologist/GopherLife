package calc

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
		maxI:   int(math.Pow(float64(Max(width, height)), 2)),
	}

}

func (s *Spiral) Next() (Coordinates, bool) {

	var coords = Coordinates{-1, -1}
	var foundCoords = false

	for s.i < s.maxI {

		if ((-s.bigX < s.x) && (s.x <= s.bigX)) &&
			((-s.bigY < s.y) && (s.y <= s.bigY)) {
			coords.X, coords.Y = int(s.x), int(s.y)
			foundCoords = true
		}

		if (s.x == s.y) || (s.x < 0 && s.x == -s.y) || (s.x > 0 && s.x == 1-s.y) {
			s.dx, s.dy = -s.dy, s.dx
		}

		s.x, s.y = s.x+s.dx, s.y+s.dy

		s.i = s.i + 1

		if foundCoords {
			break
		}

	}

	if !foundCoords {
		return coords, false
	}

	return coords, true
}
