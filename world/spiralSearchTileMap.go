package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

type GopherMap struct {
	Searchable
	TileContainer
	FoodRespawnPickup
	GopherGeneration
	Insertable

	GopherSliceAndChannel
	QueueableActions

	GopherWaitGroup *sync.WaitGroup
	IsPaused        bool
	Moments         int
	SelectedGopher  *Gopher
	diagnostics     Diagnostics

	Statistics
}

func CreateWorldCustom(statistics Statistics) *GopherMap {

	tileMap := GopherMap{}
	tileMap.Statistics = statistics

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)
	tileMap.QueueableActions = &qa

	a := NewBasic2DContainer(0, 0, statistics.Width, statistics.Height)

	tileMap.TileContainer = &a
	tileMap.Insertable = &a

	s := SpiralTileSearch{TileContainer: tileMap.TileContainer}

	tileMap.Searchable = &s

	tileMap.GopherSliceAndChannel = GopherSliceAndChannel{
		ActiveActors: make(chan *GopherActor, statistics.MaximumNumberOfGophers*2),
		ActiveArray:  make([]*GopherActor, statistics.NumberOfGophers),
	}

	frp := FoodRespawnPickup{Insertable: tileMap.Insertable}
	tileMap.FoodRespawnPickup = frp

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

	ag := GopherGeneration{
		Insertable:            tileMap.Insertable,
		maxGenerations:        tileMap.Statistics.MaximumNumberOfGophers,
		GopherSliceAndChannel: &tileMap.GopherSliceAndChannel,
	}

	tileMap.GopherGeneration = ag

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

		var gopherActor = GopherActor{
			Gopher:           &gopher,
			GopherBirthRate:  tileMap.Statistics.GopherBirthRate,
			QueueableActions: tileMap.QueueableActions,
			Searchable:       tileMap.Searchable,
			TileContainer:    tileMap.TileContainer,
			Insertable:       tileMap.Insertable,
			PickableTiles:    tileMap,
			MoveableActors:   tileMap,
			ActorGeneration:  &tileMap.GopherGeneration,
		}

		tileMap.ActiveArray[i] = &gopherActor
		tileMap.ActiveActors <- &gopherActor
		count++
	}

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

	return nil, true
}

func (tileMap *GopherMap) SelectedTile() (*Tile, bool) {

	if tileMap.SelectedGopher != nil {
		if tile, ok := tileMap.Tile(tileMap.SelectedGopher.Position.GetX(), tileMap.SelectedGopher.Position.GetY()); ok {
			return tile, ok
		}
	}
	return nil, false

}

type MoveableActors interface {
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
	tileMap.SelectedGopher = tileMap.ActiveArray[rand.Intn(len(tileMap.ActiveArray))].Gopher
}

func (tileMap *GopherMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *GopherMap) Stats() *Statistics {
	return &tileMap.Statistics
}

func (tileMap *GopherMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

type PickableTiles interface {
	PickUpFood(x int, y int) (*Food, bool)
}

type FoodRespawnPickup struct {
	Insertable
}

func (frp *FoodRespawnPickup) PickUpFood(x int, y int) (*Food, bool) {

	food, ok := frp.RemoveFood(x, y)
	defer func() {

		if ok {

			size := 50
			xrange, yrange := rand.Perm(size), rand.Perm(size)
			food := NewPotato()

		loop:
			for i := 0; i < size; i++ {
				for j := 0; j < size; j++ {
					newX, newY := x+xrange[i]-size/2, y+yrange[j]-size/2
					if frp.InsertFood(newX, newY, &food) {
						break loop
					}
				}
			}

		}

	}()
	return food, ok
}

type ActorGeneration interface {
	AddNewGopher(x int, y int, g *GopherActor) bool
}

type GopherGeneration struct {
	Insertable
	maxGenerations int
	*GopherSliceAndChannel
}

func (gg *GopherGeneration) AddNewGopher(x int, y int, g *GopherActor) bool {

	if len(gg.ActiveArray) <= gg.maxGenerations {
		if gg.InsertGopher(x, y, g.Gopher) {
			gg.ActiveActors <- g
			return true
		}
	}

	return false
}

type GopherSliceAndChannel struct {
	ActiveArray  []*GopherActor
	ActiveActors chan *GopherActor
}

func (tileMap *GopherMap) Act(gopher *GopherActor, channel chan *GopherActor) {
	gopher.Update()
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		gopher.QueueRemoveGopher()
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
	tileMap.GopherSliceAndChannel.ActiveArray = make([]*GopherActor, numGophers)

	secondChannel := make(chan *GopherActor, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveActors
		tileMap.ActiveArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.Act(gopher, secondChannel)

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
