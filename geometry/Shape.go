package geometry

type Rectangle struct {
	x      int
	y      int
	width  int
	height int
}

//NewRectangle Returns a new Rectangle of given width. height and position
func NewRectangle(x int, y int, width int, height int) Rectangle {

	return Rectangle{
		x:      x,
		y:      y,
		width:  width,
		height: height,
	}
}

//Contains Returns true if the Rectangle contains x and y
func (r *Rectangle) Contains(x int, y int) bool {
	if x < r.x || x >= r.width+r.x || y < r.y || y >= r.height+r.y {
		return false
	}
	return true
}
