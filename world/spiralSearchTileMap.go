package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

//GopherMap A map for Gophers!
type GopherMap struct {
	Searchable

	TileContainer
	GopherContainer
	FoodContainer

	QueueableActions
	InsertableGophers
	InsertableFood

	*GopherGeneration
	*GopherSliceAndChannel

	Actor *GopherActor

	GopherWaitGroup *sync.WaitGroup
	IsPaused        bool
	Moments         int
	SelectedGopher  *Gopher
	diagnostics     Diagnostics

	*Statistics
}

//NewGopherMap Creates a new GopherMap a GopherMap contains food and gophers and can use different actors to update the state of the map
func NewGopherMap(statistics *Statistics, s Searchable, t TileContainer, g GopherContainer, f FoodContainer, ig InsertableGophers, iff InsertableFood) GopherMap {

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)

	gsac := GopherSliceAndChannel{
		ActiveActors: make(chan *Gopher, statistics.MaximumNumberOfGophers*2),
		ActiveArray:  make([]*Gopher, statistics.NumberOfGophers),
	}

	gg := GopherGeneration{
		InsertableGophers:     ig,
		maxGenerations:        statistics.MaximumNumberOfGophers,
		GopherSliceAndChannel: &gsac,
	}

	var wg sync.WaitGroup

	return GopherMap{
		Searchable:        s,
		TileContainer:     t,
		GopherContainer:   g,
		FoodContainer:     f,
		InsertableGophers: ig,
		InsertableFood:    iff,

		QueueableActions: &qa,

		GopherGeneration:      &gg,
		GopherSliceAndChannel: &gsac,

		GopherWaitGroup: &wg,

		Statistics: statistics,
	}

}

func CreateWorldCustom(statistics Statistics) *GopherMap {

	b2dc := NewBasic2DContainer(0, 0, statistics.Width, statistics.Height)
	sts := SpiralTileSearch{TileContainer: &b2dc}

	tileMap := NewGopherMap(&statistics, &sts, &b2dc, &b2dc, &b2dc, &b2dc, &b2dc)

	tileMap.setUpTiles()
	return &tileMap

}

func (tileMap *GopherMap) setUpTiles() {

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		tileMap.Statistics.Width, tileMap.Statistics.Height)

	count := 0

	for i := 0; i < tileMap.Statistics.NumberOfGophers; i++ {

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
		GopherBirthRate:  tileMap.Statistics.GopherBirthRate,
		QueueableActions: tileMap.QueueableActions,
		Searchable:       tileMap.Searchable,
		GopherContainer:  tileMap.GopherContainer,
		FoodContainer:    tileMap.FoodContainer,
		PickableTiles:    tileMap,
		MoveableGophers:  tileMap,
		ActorGeneration:  tileMap.GopherGeneration,
	}

	tileMap.Actor = &actor

	for i := 0; i < tileMap.Statistics.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		tileMap.InsertFood(pos.GetX(), pos.GetY(), &food)
		count++
	}

}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the tileMap
//If there is not a gopher at the give coordinates this function returns zero.
func (tileMap *GopherMap) SelectEntity(x int, y int) (*Gopher, bool) {

	tileMap.SelectedGopher = nil

	if mapPoint, ok := tileMap.Tile(x, y); ok {
		if mapPoint.Gopher != nil {
			tileMap.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, false
}

type MoveableGophers interface {
	MoveGopher(gopher *Gopher, moveX int, moveY int) bool
}

//MoveGopher Handles the movement of a give gopher, Attempts to move a gopher by moveX and moveY.
func (tileMap *GopherMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := calc.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
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

func (tileMap *GopherMap) Stats() *Statistics {
	return tileMap.Statistics
}

func (tileMap *GopherMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

type PickableTiles interface {
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
	InsertableGophers
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

	if !tileMap.diagnostics.globalStopWatch.IsStarted() {
		tileMap.diagnostics.globalStopWatch.Start()
	}

	tileMap.diagnostics.processStopWatch.Start()
	tileMap.processGophers()
	tileMap.processQueuedTasks()
	tileMap.Statistics.NumberOfGophers = len(tileMap.ActiveActors)
	if tileMap.Statistics.NumberOfGophers > 0 {
		tileMap.Moments++
	}

	tileMap.diagnostics.processStopWatch.Stop()

	return true

}

func (tileMap *GopherMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

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

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *GopherMap) processQueuedTasks() {
	tileMap.diagnostics.inputStopWatch.Start()
	tileMap.QueueableActions.Process()
	tileMap.diagnostics.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *GopherMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

type SpiralTileSearch struct {
	TileContainer
}

func (spiralTileSearch *SpiralTileSearch) Search(position calc.Coordinates, width int, height int, max int, searchType SearchType) []calc.Coordinates {

	var coordsArray = []calc.Coordinates{}

	spiral := calc.NewSpiral(width, height)

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

	calc.SortByNearestFromCoordinate(position, coordsArray)

	return coordsArray
}
