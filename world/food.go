package world

import "gopherlife/geometry"

type Food struct {
	Name     string
	Energy   int
	Position geometry.Coordinates
}

func NewPotato() Food {
	return Food{Name: "Patato", Energy: 50}
}
