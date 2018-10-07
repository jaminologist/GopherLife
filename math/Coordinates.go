package math

import (
	"fmt"
	"sort"
	"strconv"
	s "strings"
)

type Coordinates struct {
	X int
	Y int
}

type ByNearest []Coordinates

func (a ByNearest) Len() int { return len(a) }
func (a ByNearest) Less(i, j int) bool {

	ix, iy, jx, jy := abs(a[i].X), abs(a[i].Y), abs(a[j].X), abs(a[j].Y)

	return (ix + iy) < (jx + jy)
}
func (a ByNearest) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func NewCoordinate(X int, Y int) Coordinates {
	return Coordinates{X, Y}
}

func RelativeCoordinate(c Coordinates, X int, Y int) Coordinates {
	return Coordinates{c.X + X, c.Y + Y}
}

func StringToCoordinates(coordString string) Coordinates {
	var split = s.Split(coordString, ",")
	var X, _ = strconv.Atoi(split[0])
	var Y, _ = strconv.Atoi(split[1])
	return Coordinates{X, Y}
}

func (c Coordinates) RelativeCoordinate(X int, Y int) Coordinates {
	return Coordinates{c.X + X, c.Y + Y}
}

func (c *Coordinates) Add(X int, Y int) {
	c.X += X
	c.Y += Y
}

func (c *Coordinates) Set(X int, Y int) {
	c.X = X
	c.Y = Y
}

func (c *Coordinates) SetX(X int) {
	c.X = X
}

func (c *Coordinates) GetX() int {
	return c.X
}

func (c *Coordinates) GetY() int {
	return c.Y
}

func (c *Coordinates) SetY(Y int) {
	c.Y = Y
}

func (c *Coordinates) Equals(c2 *Coordinates) bool {
	return c.X == c2.X && c.Y == c2.Y
}

func (c *Coordinates) MapKey() string {
	return CoordinateMapKey(c.X, c.Y)
}

func CoordinateMapKey(X int, Y int) string {
	return fmt.Sprintf("%[1]d,%[2]d", X, Y)
}

func SortCoordinatesUsingCoordinate(coords Coordinates, cs []Coordinates) {

	sort.Slice(cs, func(i, j int) bool {

		ix := abs(coords.X - cs[i].X)
		iy := abs(coords.Y - cs[i].Y)
		jx := abs(coords.X - cs[j].X)
		jy := abs(coords.Y - cs[j].Y)

		return (ix + iy) < (jx + jy)
	})

}
