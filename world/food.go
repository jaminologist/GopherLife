package world

import "gopherlife/geometry"

type Food struct {
	Name     string
	Energy   int
	Position geometry.Coordinates
}

func NewCarrot() Food {
	return Food{Name: "Carrot", Energy: 5}
}

func NewPotato() Food {
	return Food{Name: "Patato", Energy: 50}
}

func New() Food {
	return Food{Name: "Bean", Energy: 5}
}
