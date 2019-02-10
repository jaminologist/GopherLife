package world

type Controllable interface {
	Click(x int, y int)
	KeyPress(key int)
}

type Keys int

var leftArrow = "37"
var rightArrow = "39"
var upArrow = "38"
var downArrow = "40"

var pKey = "80"

var qKey = "81"
var wKey = "87"
