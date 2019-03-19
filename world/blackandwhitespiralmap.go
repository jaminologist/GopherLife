package world

import (
	"gopherlife/geometry"
	"gopherlife/names"
	"gopherlife/timer"
	"sync"
	"time"
)

type SpiralMapSettings struct {
	Dimensions
	MaxPopulation int
	WeirdSpiral   bool
}

//SpiralMap spins right round
type SpiralMap struct {
	SpiralMapSettings

	TileContainer
	GopherInserterAndRemover
	ActionQueuer

	ActiveActors chan *SpiralGopher
	*sync.WaitGroup

	nextSpawnCount int

	FrameTimer timer.StopWatch
}

func NewSpiralMap(settings SpiralMapSettings) SpiralMap {

	spiralMap := SpiralMap{}

	b2d := NewBasic2DContainer(0, 0, settings.Width, settings.Height)

	qa := NewFiniteActionQueue(settings.MaxPopulation * 2)
	spiralMap.ActionQueuer = &qa

	spiralMap.TileContainer = &b2d
	spiralMap.GopherInserterAndRemover = &b2d

	spiralMap.SpiralMapSettings = settings

	spiralMap.ActiveActors = make(chan *SpiralGopher, settings.MaxPopulation*2)

	var wg sync.WaitGroup
	spiralMap.WaitGroup = &wg

	spiralMap.AddNewSpiralGopher()

	return spiralMap
}

func (spiralMap *SpiralMap) Update() bool {

	spiralMap.FrameTimer.Start()

	numGophers := len(spiralMap.ActiveActors)
	secondChannel := make(chan *SpiralGopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-spiralMap.ActiveActors
		spiralMap.WaitGroup.Add(1)

		go func() {

			gopher.Update()

			if !gopher.IsDead {
				secondChannel <- gopher
			} else {
				gopher.Add(func() {
					spiralMap.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
				})
			}

			spiralMap.WaitGroup.Done()

		}()
	}

	spiralMap.ActiveActors = secondChannel
	spiralMap.WaitGroup.Wait()

	spiralMap.nextSpawnCount++

	if spiralMap.nextSpawnCount > 2 {
		spiralMap.nextSpawnCount = 0
		spiralMap.AddNewSpiralGopher()
	}
	spiralMap.Process()

	for spiralMap.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*time.Duration(2) {
	}

	return true

}

func (spiralMap *SpiralMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := geometry.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if spiralMap.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		spiralMap.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false

}

func (spiralMap *SpiralMap) AddNewSpiralGopher() {

	gopher := NewGopher(names.CuteName(), geometry.Coordinates{0, 0})

	//Commented out for cool spiral effect 1
	if !spiralMap.WeirdSpiral {
		spiralMap.InsertGopher(spiralMap.Width/2, spiralMap.Height/2, &gopher)
	}

	spiral := geometry.NewSpiral(spiralMap.Width, spiralMap.Height)

	sg := SpiralGopher{
		TileContainer:   spiralMap,
		ActionQueuer:    spiralMap.ActionQueuer,
		MoveableGophers: spiralMap,
		Gopher:          &gopher,
		settings:        &spiralMap.SpiralMapSettings,
		Spiral:          &spiral,
	}

	spiralMap.ActiveActors <- &sg

}

type SpiralGopher struct {
	TileContainer
	ActionQueuer
	MoveableGophers
	*Gopher
	*geometry.Spiral
	settings *SpiralMapSettings
}

func (gopher *SpiralGopher) Update() {
	if gopher.settings.WeirdSpiral {
		gopher.UpdateWeridSpiral()
	} else {
		gopher.UpdateNormalSpiral()
	}
}

//Cool Effect 2
func (gopher *SpiralGopher) UpdateNormalSpiral() {

	position, ok := gopher.Spiral.Next()
	//position, ok = gopher.Spiral.Next()

	x, y := gopher.settings.Width/2+position.GetX(), gopher.settings.Height/2+position.GetY()

	if ok {
		gopher.Add(func() {
			gopher.MoveGopher(gopher.Gopher, x-gopher.Position.GetX(), y-gopher.Position.GetY())
		})
	} else {
		gopher.IsDead = true
	}

}

//Cool Effect 1
func (gopher *SpiralGopher) UpdateWeridSpiral() {

	position, ok := gopher.Spiral.Next()
	//position, ok = gopher.Spiral.Next()

	//x, y := gopher.Statistics.Width/2+position.GetX(), gopher.Statistics.Height/2+position.GetY()

	if ok {
		gopher.Add(func() {
			gopher.MoveGopher(gopher.Gopher, position.GetX(), position.GetY())
		})
	} else {
		gopher.IsDead = true
	}

}
