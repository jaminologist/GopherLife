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

//GopherMap A map for Gophers!
type GopherMap struct {
	Searchable

	TileContainer
	GopherContainer
	FoodContainer

	ActionQueuer
	GopherInserter
	FoodInserter

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

//NewGopherMap Creates a new GopherMap a GopherMap contains food and gophers and can use different actors to update the state of the map
func NewGopherMap(settings *GopherMapSettings, s Searchable, t TileContainer, g GopherContainer, f FoodContainer, ig GopherInserter, iff FoodInserter) GopherMap {

	qa := NewBasicActionQueue(settings.MaxPopulation * 2)

	gsac := GopherSliceAndChannel{
		ActiveActors: make(chan *Gopher, settings.MaxPopulation*2),
		ActiveArray:  make([]*Gopher, settings.InitialPopulation),
	}

	gg := GopherGeneration{
		GopherInserter:        ig,
		maxGenerations:        settings.MaxPopulation,
		GopherSliceAndChannel: &gsac,
	}

	var wg sync.WaitGroup

	return GopherMap{
		Searchable:      s,
		TileContainer:   t,
		GopherContainer: g,
		FoodContainer:   f,
		GopherInserter:  ig,
		FoodInserter:    iff,

		ActionQueuer: &qa,

		GopherGeneration:      &gg,
		GopherSliceAndChannel: &gsac,

		GopherWaitGroup: &wg,

		GopherMapSettings: settings,

		NumberOfGophers: settings.InitialPopulation,
	}

}

func CreateWorldCustom(settings GopherMapSettings) *GopherMap {

	b2dc := NewBasic2DContainer(0, 0, settings.Width, settings.Height)
	sts := SpiralTileSearch{TileContainer: &b2dc}

	tileMap := NewGopherMap(&settings, &sts, &b2dc, &b2dc, &b2dc, &b2dc, &b2dc)

	tileMap.setUpTiles()
	return &tileMap

}

func (tileMap *GopherMap) setUpTiles() {

	keys := geometry.GenerateRandomizedCoordinateArray(0, 0,
		tileMap.Width, tileMap.Height)

	count := 0

	for i := 0; i < tileMap.InitialPopulation; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.CuteName(), pos)

		tileMap.InsertGopher(pos.GetX(), pos.GetY(), &gopher)

		if i == 0 {
			tileMap.SelectedGopher = &gopher
		}

		tileMap.ActiveArray[i] = &gopher
		tileMap.ActiveActors <- &gopher
		count++
	}

	actor := GopherActor{
		GopherBirthRate: tileMap.GopherBirthRate,
		ActionQueuer:    tileMap.ActionQueuer,
		Searchable:      tileMap.Searchable,
		GopherContainer: tileMap.GopherContainer,
		FoodContainer:   tileMap.FoodContainer,
		FoodPicker:      tileMap,
		MoveableGophers: tileMap,
		ActorGeneration: tileMap.GopherGeneration,
	}

	tileMap.Actor = &actor

	for i := 0; i < tileMap.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		tileMap.InsertFood(pos.GetX(), pos.GetY(), &food)
		count++
	}

}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the tileMap
//If there is not a gopher at the give coordinates this function returns zero.
func (tileMap *GopherMap) SelectEntity(x int, y int) (*Gopher, bool) {

	if gopher, ok := tileMap.HasGopher(x, y); ok {
		tileMap.SelectedGopher = gopher
		return gopher, true
	}

	return nil, false
}

type MoveableGophers interface {
	MoveGopher(gopher *Gopher, moveX int, moveY int) bool
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (tileMap *GopherMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := geometry.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if tileMap.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		tileMap.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false

}

func (tileMap *GopherMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.ActiveArray[rand.Intn(len(tileMap.ActiveArray))]
}

func (tileMap *GopherMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *GopherMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

type FoodPicker interface {
	PickUpFood(x int, y int) (*Food, bool)
}

func (tileMap *GopherMap) PickUpFood(x int, y int) (*Food, bool) {

	food, ok := tileMap.RemoveFood(x, y)
	defer func() {

		if ok {

			size := 50
			xrange, yrange := rand.Perm(size), rand.Perm(size)
			food := NewPotato()

		loop:
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					newX, newY := x+xrange[i]-size/2, y+yrange[j]-size/2
					if tileMap.InsertFood(newX, newY, &food) {
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
	GopherInserter
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

func (tileMap *GopherMap) Act(actor *GopherActor, gopher *Gopher, channel chan *Gopher) {
	actor.Update(gopher)
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		tileMap.Add(func() {
			tileMap.RemoveGopher(gopher.Position.GetX(), gopher.Position.GetY())
		})
	}
	tileMap.GopherWaitGroup.Done()
}

func (tileMap *GopherMap) Update() bool {

	if tileMap.IsPaused {
		return false
	}

	if tileMap.SelectedGopher != nil && tileMap.SelectedGopher.IsDecayed() {
		tileMap.SelectRandomGopher()
	}

	if !tileMap.diagnostics.GlobalStopWatch.IsStarted() {
		tileMap.diagnostics.GlobalStopWatch.Start()
	}

	tileMap.diagnostics.ProcessStopWatch.Start()
	tileMap.processGophers()
	tileMap.processQueuedTasks()
	tileMap.NumberOfGophers = len(tileMap.ActiveActors)

	tileMap.diagnostics.ProcessStopWatch.Stop()

	return true

}

func (tileMap *GopherMap) processGophers() {

	tileMap.diagnostics.GopherStopWatch.Start()

	numGophers := len(tileMap.ActiveActors)
	tileMap.GopherSliceAndChannel.ActiveArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveActors
		tileMap.ActiveArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.Act(tileMap.Actor, gopher, secondChannel)

	}
	tileMap.ActiveActors = secondChannel
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.GopherStopWatch.Stop()
}

func (tileMap *GopherMap) processQueuedTasks() {
	tileMap.diagnostics.InputStopWatch.Start()
	tileMap.ActionQueuer.Process()
	tileMap.diagnostics.InputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *GopherMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
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
