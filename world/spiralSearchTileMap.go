package world

import (
	"gopherlife/geometry"
	"gopherlife/names"
	"math/rand"
	"sync"
)

//GopherMapSettings sets configuration for a Gopher Map
type GopherMapSettings struct {
	Dimensions
	Population

	GopherBirthRate int
	NumberOfFood    int
}

//GopherWorld A map for Gophers!
type GopherWorld struct {
	Searcher
	TileContainer
	GopherContainer
	FoodContainer

	ActionQueuer
	GopherInserterAndRemover
	FoodInserterAndRemover

	*GopherGeneration
	*GopherSliceAndChannel

	Actor *GopherActor

	GopherWaitGroup *sync.WaitGroup
	IsPaused        bool
	SelectedGopher  *Gopher
	diagnostics     Diagnostics

	NumberOfGophers int

	*GopherMapSettings
}

//NewGopherWorld Creates a new GopherWorld a GopherWorld contains food and gophers and can use different actors to update the state of the map
func NewGopherWorld(settings *GopherMapSettings, s Searcher, t TileContainer, g GopherContainer, f FoodContainer, ig GopherInserterAndRemover, iff FoodInserterAndRemover) GopherWorld {

	qa := NewFiniteActionQueue(settings.MaxPopulation * 2)

	gsac := GopherSliceAndChannel{
		ActiveActors: make(chan *Gopher, settings.MaxPopulation*2),
		ActiveArray:  make([]*Gopher, settings.InitialPopulation),
	}

	gg := GopherGeneration{
		GopherInserterAndRemover: ig,
		maxGenerations:           settings.MaxPopulation,
		GopherSliceAndChannel:    &gsac,
	}

	var wg sync.WaitGroup

	return GopherWorld{
		Searcher:                 s,
		TileContainer:            t,
		GopherContainer:          g,
		FoodContainer:            f,
		GopherInserterAndRemover: ig,
		FoodInserterAndRemover:   iff,

		ActionQueuer: &qa,

		GopherGeneration:      &gg,
		GopherSliceAndChannel: &gsac,

		GopherWaitGroup: &wg,

		GopherMapSettings: settings,

		NumberOfGophers: settings.InitialPopulation,
	}

}

func CreateWorldCustom(settings GopherMapSettings) *GopherWorld {

	b2dc := NewBasic2DContainer(0, 0, settings.Width, settings.Height)
	sts := SpiralTileSearch{TileContainer: &b2dc}

	gw := NewGopherWorld(&settings, &sts, &b2dc, &b2dc, &b2dc, &b2dc, &b2dc)

	gw.setUpTiles()
	return &gw

}

func (gw *GopherWorld) setUpTiles() {

	keys := geometry.GenerateRandomizedCoordinateArray(0, 0,
		gw.Width, gw.Height)

	count := 0

	for i := 0; i < gw.InitialPopulation; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.CuteName(), pos)

		gw.InsertGopher(pos.GetX(), pos.GetY(), &gopher)

		if i == 0 {
			gw.SelectedGopher = &gopher
		}

		gw.ActiveArray[i] = &gopher
		gw.ActiveActors <- &gopher
		count++
	}

	actor := GopherActor{
		GopherBirthRate: gw.GopherBirthRate,
		ActionQueuer:    gw.ActionQueuer,
		Searcher:        gw.Searcher,
		GopherContainer: gw.GopherContainer,
		FoodContainer:   gw.FoodContainer,
		FoodPicker:      gw,
		MoveableGophers: gw,
		ActorGeneration: gw.GopherGeneration,
	}

	gw.Actor = &actor

	for i := 0; i < gw.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		gw.InsertFood(pos.GetX(), pos.GetY(), &food)
		count++
	}

}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the GopherWorld
//If there is not a gopher at the give coordinates this function returns zero.
func (gw *GopherWorld) SelectEntity(x int, y int) (*Gopher, bool) {

	if gopher, ok := gw.HasGopher(x, y); ok {
		gw.SelectedGopher = gopher
		return gopher, true
	}

	return nil, false
}

type MoveableGophers interface {
	MoveGopher(gopher *Gopher, moveX int, moveY int) bool
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (gw *GopherWorld) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := geometry.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if gw.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		gw.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false

}

func (gw *GopherWorld) SelectRandomGopher() {
	gw.SelectedGopher = gw.ActiveArray[rand.Intn(len(gw.ActiveArray))]
}

func (gw *GopherWorld) UnSelectGopher() {
	gw.SelectedGopher = nil
}

func (gw *GopherWorld) Diagnostics() *Diagnostics {
	return &gw.diagnostics
}

type FoodPicker interface {
	PickUpFood(x int, y int) (*Food, bool)
}

func (gw *GopherWorld) PickUpFood(x int, y int) (*Food, bool) {

	food, ok := gw.RemoveFood(x, y)
	defer func() {

		if ok {

			size := 50
			xrange, yrange := rand.Perm(size), rand.Perm(size)
			food := NewPotato()

		loop:
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					newX, newY := x+xrange[i]-size/2, y+yrange[j]-size/2
					if gw.InsertFood(newX, newY, &food) {
						break loop
					}
				}
			}

		}

	}()
	return food, ok
}

type ActorGeneration interface {
	AddNewGopher(x int, y int, g *Gopher) bool
}

type GopherGeneration struct {
	GopherInserterAndRemover
	maxGenerations int
	*GopherSliceAndChannel
}

func (gg *GopherGeneration) AddNewGopher(x int, y int, gopher *Gopher) bool {

	if len(gg.ActiveArray) <= gg.maxGenerations {
		if gg.InsertGopher(x, y, gopher) {
			gg.ActiveActors <- gopher
			return true
		}
	}

	return false
}

type GopherSliceAndChannel struct {
	ActiveArray  []*Gopher
	ActiveActors chan *Gopher
}

func (gw *GopherWorld) Act(actor *GopherActor, gopher *Gopher, channel chan *Gopher) {
	actor.Update(gopher)
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		gw.Add(func() {
			gw.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
		})
	}
	gw.GopherWaitGroup.Done()
}

func (gw *GopherWorld) Update() bool {

	if gw.IsPaused {
		return false
	}

	if gw.SelectedGopher != nil && gw.SelectedGopher.IsDecayed() {
		gw.SelectRandomGopher()
	}

	if !gw.diagnostics.GlobalStopWatch.IsStarted() {
		gw.diagnostics.GlobalStopWatch.Start()
	}

	gw.diagnostics.ProcessStopWatch.Start()
	gw.processGophers()

	gw.processQueuedTasks()

	gw.NumberOfGophers = len(gw.ActiveActors)

	gw.diagnostics.ProcessStopWatch.Stop()

	return true

}

func (gw *GopherWorld) processGophers() {

	gw.diagnostics.GopherStopWatch.Start()

	numGophers := len(gw.ActiveActors)
	gw.GopherSliceAndChannel.ActiveArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-gw.ActiveActors
		gw.ActiveArray[i] = gopher
		gw.GopherWaitGroup.Add(1)
		go gw.Act(gw.Actor, gopher, secondChannel)

	}
	gw.ActiveActors = secondChannel
	gw.GopherWaitGroup.Wait()

	gw.diagnostics.GopherStopWatch.Stop()
}

func (gw *GopherWorld) processQueuedTasks() {
	gw.diagnostics.InputStopWatch.Start()
	gw.ActionQueuer.Process()
	gw.diagnostics.InputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (gw *GopherWorld) TogglePause() {
	gw.IsPaused = !gw.IsPaused
}

type SpiralTileSearch struct {
	TileContainer
}

func (spiralTileSearch *SpiralTileSearch) Search(position geometry.Coordinates, width int, height int, max int, searchType SearchType) []geometry.Coordinates {

	var coordsArray = []geometry.Coordinates{}

	spiral := geometry.NewSpiral(width, height)

	var query TileQuery

	switch searchType {
	case SearchForFood:
		query = CheckMapPointForFood
	case SearchForEmptySpace:
		query = CheckMapPointForEmptySpace
	case SearchForFemaleGopher:
		query = CheckMapPointForFemaleGopher
	}

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > max {
			break
		}

		relativeCoords := position.RelativeCoordinate(coordinates.X, coordinates.Y)

		if tile, ok := spiralTileSearch.Tile(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if query(tile) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	geometry.SortByNearestFromCoordinate(position, coordsArray)

	return coordsArray
}
