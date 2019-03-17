package controllers

//UserInputHandler used to handle user inputs for a world
type UserInputHandler interface {
	Scroller
	Click(x int, y int)
	KeyPress(key Keys)
}

type Scroller interface {
	Scroll(deltaY int)
}

type NoPlayerInput struct{}

func (controller *NoPlayerInput) Click(x int, y int) {
}

func (controller *NoPlayerInput) KeyPress(key Keys) {
}

type WorldPageData struct {
	PageTitle   string
	FormData    []FormData
	IsGopherWorld bool
}

//Keys used to denote correct number for keys on a keyboard when calling e.which in js
type Keys int64

//List of keys and their corresponding numbers
const (
	LeftArrow  Keys = 37
	RightArrow Keys = 39
	UpArrow    Keys = 38
	DownArrow  Keys = 40

	PKey Keys = 80
	QKey Keys = 81
	WKey Keys = 87
)
