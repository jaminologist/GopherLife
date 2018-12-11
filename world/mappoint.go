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

func (mp *MapPoint) HasGopher() bool {
	return mp.Gopher == nil
}

func (mp *MapPoint) HasFood() bool {
	return mp.Food == nil
}

func (mp *MapPoint) SetGopher(g *Gopher) {
	mp.Gopher = g
}

func (mp *MapPoint) SetFood(f *Food) {
	mp.Food = f
}
