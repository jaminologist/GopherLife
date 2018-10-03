package world

import (
	"gopherlife/food"
)

type MapPoint struct {
	Gopher *Gopher
	Food   *food.Food
}

func (mp *MapPoint) isEmpty() bool {
	return mp.Gopher == nil && mp.Food == nil
}
