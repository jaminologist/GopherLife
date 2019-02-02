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

	QueueGopherMove(gopher *Gopher, x int, y int)
	QueuePickUpFood(gopher *Gopher)
	QueueMating(mate *Gopher, coords calc.Coordinates)
	QueueRemoveGopher(gopher *Gopher)

	TileContainer
	Searchable
}

type Searchable interface {
	Search(position calc.Coordinates, width int, height int, max int, searchType SearchType) []calc.Coordinates
}

type SearchType int

const (
	SearchForFood SearchType = iota
	SearchForEmptySpace
	SearchForFemaleGopher
	FemaleGopher
)

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
