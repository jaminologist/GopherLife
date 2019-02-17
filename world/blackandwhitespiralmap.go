package world

import (
	"gopherlife/calc"
	"sync"
)

type SpiralMap struct {
	TileContainer
	Insertable
	QueueableActions

	ActiveArray []*SpiralGopher

	*sync.WaitGroup

	count int
}

func NewSpiralMap(stats Statistics) {

	spiralMap := SpiralMap{}

	b2d := NewBasic2DContainer(0, 0, stats.Width, stats.Height)

	qa := NewBasicActionQueue(stats.MaximumNumberOfGophers * 2)
	spiralMap.QueueableActions = &qa

	spiralMap.TileContainer = &b2d
	spiralMap.Insertable = &b2d

}

func (spiralMap *SpiralMap) Update() {

	spiralMap.Process()

	spiralMap.count++

	if spiralMap.count > 5 {

	}

}

type SpiralGopher struct {
	TileContainer
	QueueableActions
	Insertable
	MoveableActors
	*Gopher
	calc.Spiral
}

func (gopher *SpiralGopher) Update() {

	position, ok := gopher.Spiral.Next()

	if ok {
		gopher.Add(func() {
			gopher.MoveGopher(gopher.Gopher, position.GetX(), position.GetY())
		})
	} else {
		gopher.Add(func() {
			gopher.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
		})
	}

}
