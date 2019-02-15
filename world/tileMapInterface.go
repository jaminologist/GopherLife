package world

import "gopherlife/calc"

//TileMap is cool
type TileMap interface {
	Update() bool

	TogglePause()

	SelectedTile() (*Tile, bool)

	SelectEntity(x int, y int) (*Gopher, bool)
	SelectRandomGopher()
	UnSelectGopher()

	Stats() *Statistics
	Diagnostics() *Diagnostics

	TileContainer
	Searchable
}

type Statistics struct {
	Width                  int
	Height                 int
	NumberOfGophers        int
	MaximumNumberOfGophers int
	GopherBirthRate        int
	NumberOfFood           int
}

type Diagnostics struct {
	globalStopWatch  calc.StopWatch
	inputStopWatch   calc.StopWatch
	gopherStopWatch  calc.StopWatch
	processStopWatch calc.StopWatch
}
