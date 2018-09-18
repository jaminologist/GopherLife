package world

import (
	"gopherlife/animal"
	"gopherlife/food"
)

type MapPoint struct {
	Gopher *animal.Gopher
	Food   *food.Food
}

func (mp *MapPoint) isEmpty() bool {
	return mp.Gopher == nil && mp.Food == nil
}
