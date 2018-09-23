package math

import "fmt"
import s "strings"
import "strconv"

type Coordinates struct {
	x int
	y int
}

type ByNearest []Coordinates

func (a ByNearest) Len() int { return len(a) }
func (a ByNearest) Less(i, j int) bool {

	ix, iy, jx, jy := abs(a[i].x), abs(a[i].y), abs(a[j].x), abs(a[j].y)

	return (ix + iy) < (jx + jy)
}
func (a ByNearest) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func NewCoordinate(x int, y int) Coordinates {
	return Coordinates{x, y}
}

func RelativeCoordinate(c Coordinates, x int, y int) Coordinates {
	return Coordinates{c.x + x, c.y + y}
}

func StringToCoordinates(coordString string) Coordinates {
	var split = s.Split(coordString, ",")
	var x, _ = strconv.Atoi(split[0])
	var y, _ = strconv.Atoi(split[1])
	return Coordinates{x, y}
}

func (c Coordinates) RelativeCoordinate(x int, y int) Coordinates {
	return Coordinates{c.x + x, c.y + y}
}

func (c *Coordinates) Add(x int, y int) {
	c.x += x
	c.y += y
}

func (c *Coordinates) Set(x int, y int) {
	c.x = x
	c.y = y
}

func (c *Coordinates) SetX(x int) {
	c.x = x
}

func (c *Coordinates) GetX() int {
	return c.x
}

func (c *Coordinates) GetY() int {
	return c.y
}

func (c *Coordinates) SetY(y int) {
	c.y = y
}

func (c *Coordinates) Equals(c2 *Coordinates) bool {
	return c.x == c2.x && c.y == c2.y
}

func (c *Coordinates) MapKey() string {
	return CoordinateMapKey(c.x, c.y)
}

func CoordinateMapKey(x int, y int) string {
	return fmt.Sprintf("%[1]d,%[2]d", x, y)
}