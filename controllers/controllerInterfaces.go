package controllers

//UserInputHandler handles the 'Scroll', 'Click' and 'KeyPress' user inputs
type UserInputHandler interface {
	Scroller
	Clicker
	KeyPresser
}

//Scroller handles the 'Scroll' user input
type Scroller interface {
	Scroll(deltaY int)
}

//Clicker handles the 'Click' user input
type Clicker interface {
	Click(x int, y int)
}

//KeyPresser handles the 'KeyPress' user input
type KeyPresser interface {
	KeyPress(key Keys)
}

//Keys the number assigned to a keyboard 'key' when calling e.which in js
type Keys int64

//List of used keys and their corresponding numbers
const (
	LeftArrow  Keys = 37
	RightArrow Keys = 39
	UpArrow    Keys = 38
	DownArrow  Keys = 40

	PKey Keys = 80
	QKey Keys = 81
	WKey Keys = 87
)

type NoPlayerInput struct{}

func (controller *NoPlayerInput) Click(x int, y int) {
}

func (controller *NoPlayerInput) KeyPress(key Keys) {
}

type WorldPageData struct {
	PageTitle     string
	FormData      []FormData
	IsGopherWorld bool
}
