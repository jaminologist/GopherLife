package world

import (
	"gopherlife/geometry"
	"gopherlife/names"
	"gopherlife/timer"
	"sync"
	"time"
)

type SpiralWorldSettings struct {
	Dimensions
	MaxPopulation int
	WeirdSpiral   bool
}

//SpiralWorld spins right round
type SpiralWorld struct {
	SpiralWorldSettings

	TileContainer
	GopherInserterAndRemover
	ActionQueuer

	ActiveActors chan *SpiralGopher
	*sync.WaitGroup

	nextSpawnCount int

	FrameTimer timer.StopWatch
}

func NewSpiralWorld(settings SpiralWorldSettings) SpiralWorld {

	spiralWorld := SpiralWorld{}

	b2d := NewBasic2DContainer(0, 0, settings.Width, settings.Height)

	qa := NewFiniteActionQueue(settings.MaxPopulation * 2)
	spiralWorld.ActionQueuer = &qa

	spiralWorld.TileContainer = &b2d
	spiralWorld.GopherInserterAndRemover = &b2d

	spiralWorld.SpiralWorldSettings = settings

	spiralWorld.ActiveActors = make(chan *SpiralGopher, settings.MaxPopulation*2)

	var wg sync.WaitGroup
	spiralWorld.WaitGroup = &wg

	spiralWorld.AddNewSpiralGopher()

	return spiralWorld
}

func (spiralWorld *SpiralWorld) Update() bool {

	spiralWorld.FrameTimer.Start()

	numGophers := len(spiralWorld.ActiveActors)
	secondChannel := make(chan *SpiralGopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-spiralWorld.ActiveActors
		spiralWorld.WaitGroup.Add(1)

		go func() {

			gopher.Update()

			if !gopher.IsDead {
				secondChannel <- gopher
			} else {
				gopher.Add(func() {
					spiralWorld.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
				})
			}

			spiralWorld.WaitGroup.Done()

		}()
	}

	spiralWorld.ActiveActors = secondChannel
	spiralWorld.WaitGroup.Wait()

	spiralWorld.nextSpawnCount++

	if spiralWorld.nextSpawnCount > 2 {
		spiralWorld.nextSpawnCount = 0
		spiralWorld.AddNewSpiralGopher()
	}
	spiralWorld.Process()

	for spiralWorld.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*time.Duration(2) {
	}

	return true

}

func (spiralWorld *SpiralWorld) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := geometry.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if spiralWorld.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		spiralWorld.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false

}

func (spiralWorld *SpiralWorld) AddNewSpiralGopher() {

	gopher := NewGopher(names.CuteName(), geometry.Coordinates{0, 0})

	if !spiralWorld.WeirdSpiral {
		spiralWorld.InsertGopher(spiralWorld.Width/2, spiralWorld.Height/2, &gopher)
	}

	spiral := geometry.NewSpiral(spiralWorld.Width, spiralWorld.Height)

	sg := SpiralGopher{
		TileContainer:   spiralWorld,
		ActionQueuer:    spiralWorld.ActionQueuer,
		MoveableGophers: spiralWorld,
		Gopher:          &gopher,
		settings:        &spiralWorld.SpiralWorldSettings,
		Spiral:          &spiral,
	}

	spiralWorld.ActiveActors <- &sg

}

type SpiralGopher struct {
	TileContainer
	ActionQueuer
	MoveableGophers
	*Gopher
	*geometry.Spiral
	settings *SpiralWorldSettings
}

func (gopher *SpiralGopher) Update() {
	if gopher.settings.WeirdSpiral {
		gopher.UpdateWeirdSpiral()
	} else {
		gopher.UpdateNormalSpiral()
	}
}

//UpdateNormalSpiral updates the position of the SpiralGopher to create a normal spiral pattern
func (gopher *SpiralGopher) UpdateNormalSpiral() {

	position, ok := gopher.Spiral.Next()

	x, y := gopher.settings.Width/2+position.GetX(), gopher.settings.Height/2+position.GetY()

	if ok {
		gopher.Add(func() {
			gopher.MoveGopher(gopher.Gopher, x-gopher.Position.GetX(), y-gopher.Position.GetY())
		})
	} else {
		gopher.IsDead = true
	}

}

//UpdateWeirdSpiral updates the position of the SpiralGopher to create a weird spiral pattern
func (gopher *SpiralGopher) UpdateWeirdSpiral() {

	position, ok := gopher.Spiral.Next()

	if ok {
		gopher.Add(func() {
			gopher.MoveGopher(gopher.Gopher, position.GetX(), position.GetY())
		})
	} else {
		gopher.IsDead = true
	}

}
