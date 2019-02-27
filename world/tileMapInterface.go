package world

import "gopherlife/calc"

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
	globalStopWatch  calc.StopWatch
	inputStopWatch   calc.StopWatch
	gopherStopWatch  calc.StopWatch
	processStopWatch calc.StopWatch
}
