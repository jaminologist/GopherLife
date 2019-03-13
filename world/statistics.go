package world

import "gopherlife/timer"

type Dimensions struct {
	Width  int
	Height int
}

//Statistics is used primarily by the 'GopherMap' struct and details
//all editable variables of the map
type Statistics struct {
	Width                  int
	Height                 int
	NumberOfGophers        int
	MaximumNumberOfGophers int
	GopherBirthRate        int
	NumberOfFood           int
}

//Diagnostics is used primarily by the 'GopherMap' struct and is used to track
//how long different parts of the 'Update' method take
type Diagnostics struct {
	GlobalStopWatch  timer.StopWatch
	InputStopWatch   timer.StopWatch
	GopherStopWatch  timer.StopWatch
	ProcessStopWatch timer.StopWatch
}
