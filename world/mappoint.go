package world

type MapPoint struct {
	Gopher *Gopher
	Food   *Food
}

func (mp *MapPoint) isEmpty() bool {
	return mp.Gopher == nil && mp.Food == nil
}

func (mp *MapPoint) HasGopher() bool {
	return mp.Gopher != nil
}

func (mp *MapPoint) HasFood() bool {
	return mp.Food != nil
}

func (mp *MapPoint) SetGopher(g *Gopher) {
	mp.Gopher = g
}

func (mp *MapPoint) SetFood(f *Food) {
	mp.Food = f
}

func (mp *MapPoint) ClearGopher() {
	mp.Gopher = nil
}

func (mp *MapPoint) ClearFood() {
	mp.Food = nil
}
