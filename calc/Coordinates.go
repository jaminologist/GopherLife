package calc

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	s "strings"
)

type Coordinates struct {
	X int
	Y int
}

//Abs returns the absolute value of an integer
func Abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func NewCoordinate(X int, Y int) Coordinates {
	return Coordinates{X, Y}
}

func StringToCoordinates(coordString string) Coordinates {
	var split = s.Split(coordString, ",")
	var X, _ = strconv.Atoi(split[0])
	var Y, _ = strconv.Atoi(split[1])
	return Coordinates{X, Y}
}

//RelativeCoordinate returns a new Coordinate that is of distance +x and +y
//from a Coordinate
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

func (c *Coordinates) GetX() int {
	return c.X
}

func (c *Coordinates) GetY() int {
	return c.Y
}

//Difference returns the difference between the x and y co-ordinates of two points
func (c *Coordinates) Difference(c2 Coordinates) (int, int) {

	diffX := c.GetX() - c2.GetX()
	diffY := c.GetY() - c2.GetY()

	return diffX, diffY
}

//IsInRange Checks if one coordinate is in range of another co-ordinate, using minimum x and y distances
func (c *Coordinates) IsInRange(c2 Coordinates, minX int, minY int) bool {
	diffX, diffY := c.Difference(c2)
	return Abs(diffX) <= minX && Abs(diffY) <= minY
}

//Equals checks the if two co-ordinates share the same x and y value
func (c *Coordinates) Equals(c2 *Coordinates) bool {
	return c.X == c2.X && c.Y == c2.Y
}

//MapKey creates a string 'key' for a Coordinate struct
func (c *Coordinates) MapKey() string {
	return CoordinateMapKey(c.X, c.Y)
}

//CoordinateMapKey converts an x and y value into a map key
func CoordinateMapKey(X int, Y int) string {
	return fmt.Sprintf("%[1]d,%[2]d", X, Y)
}

func Hashcode(x int, y int) int {
	return (31 * x) + y
}

//SortByNearestFromCoordinate Sorts an array of coordinates by nearest to the given coordinate
func SortByNearestFromCoordinate(coords Coordinates, cs []Coordinates) {

	sort.Slice(cs, func(i, j int) bool {

		ix := Abs(coords.X - cs[i].X)
		iy := Abs(coords.Y - cs[i].Y)
		jx := Abs(coords.X - cs[j].X)
		jy := Abs(coords.Y - cs[j].Y)

		return (ix + iy) < (jx + jy)
	})

}

func GenerateCoordinateArray(startX int, startY int, endX int, endY int) []Coordinates {

	slice := make([]Coordinates, Abs(startX-endX)*Abs(startX-startY))

	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			slice = append(slice, NewCoordinate(i, j))
		}
	}

	return slice
}

func GenerateRandomizedCoordinateArray(startX int, startY int, endX int, endY int) []Coordinates {

	slice := GenerateCoordinateArray(0, 0, endX, endY)

	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})

	return slice
}
