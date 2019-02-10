package world

import (
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

type SpiralSearchTileMap struct {
	Searchable
	TileContainer
	FoodRespawnPickup
	GopherGeneration
	Insertable
	GopherMapUpdater
}

func CreateWorldCustom(statistics Statistics) *SpiralSearchTileMap {

	tileMap := SpiralSearchTileMap{}

	a := NewBasic2DContainer(0, 0, statistics.Width, statistics.Height)
	tileMap.TileContainer = &a

	s := SpiralTileSearch{TileContainer: tileMap.TileContainer}
	tileMap.Searchable = &s

	qa := NewBasicActionQueue(statistics.MaximumNumberOfGophers * 2)
	tileMap.QueueableActions = &qa

	frp := FoodRespawnPickup{Insertable: &tileMap}
	tileMap.FoodRespawnPickup = frp

	tileMap.Insertable = &a

	tileMap.Statistics = statistics

	var wg sync.WaitGroup

	tileMap.GopherMapUpdater.GopherWaitGroup = &wg

	tileMap.GopherMapUpdater.GopherSliceAndChannel = GopherSliceAndChannel{
		ActiveActors: make(chan *GopherActor, statistics.MaximumNumberOfGophers*2),
		ActiveArray:  make([]*GopherActor, statistics.NumberOfGophers),
	}

	ag := GopherGeneration{
		Insertable:            tileMap.Insertable,
		maxGenerations:        tileMap.Statistics.MaximumNumberOfGophers,
		GopherSliceAndChannel: &tileMap.GopherMapUpdater.GopherSliceAndChannel,
	}

	tileMap.GopherGeneration = ag

	tileMap.setUpMapPoints()
	return &tileMap

}

func CreateTileMap() *SpiralSearchTileMap {
	tileMap := CreateWorldCustom(
		Statistics{
			Width:                  3000,
			Height:                 3000,
			NumberOfGophers:        5000,
			NumberOfFood:           1000000,
			MaximumNumberOfGophers: 100000,
			GopherBirthRate:        7,
		},
	)
	return tileMap
}

func (tileMap *SpiralSearchTileMap) setUpMapPoints() {

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

		tileMap.GopherMapUpdater.ActiveArray[i] = &gopherActor
		tileMap.GopherMapUpdater.ActiveActors <- &gopherActor
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
func (tileMap *SpiralSearchTileMap) SelectEntity(x int, y int) (*Gopher, bool) {

	tileMap.SelectedGopher = nil

	if mapPoint, ok := tileMap.Tile(x, y); ok {
		if mapPoint.Gopher != nil {
			tileMap.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, true
}

func (tileMap *SpiralSearchTileMap) SelectedTile() (*Tile, bool) {

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
func (tileMap *SpiralSearchTileMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := calc.Coordinates{X: gopher.Position.X, Y: gopher.Position.Y}
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	if tileMap.InsertGopher(targetPosition.GetX(), targetPosition.GetY(), gopher) {
		tileMap.RemoveGopher(currentPosition.GetX(), currentPosition.GetY())
		return true
	}
	return false

}

func (tileMap *GopherMapUpdater) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.ActiveArray[0].Gopher
}

func (tileMap *SpiralSearchTileMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *SpiralSearchTileMap) Stats() *Statistics {
	return &tileMap.Statistics
}

func (tileMap *SpiralSearchTileMap) Diagnostics() *Diagnostics {
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

type GopherMapUpdater struct {
	GopherSliceAndChannel
	QueueableActions

	GopherWaitGroup *sync.WaitGroup
	IsPaused        bool
	Moments         int
	SelectedGopher  *Gopher
	diagnostics     Diagnostics

	Statistics
}

func (tileMap *GopherMapUpdater) Act(gopher *GopherActor, channel chan *GopherActor) {
	gopher.Update()
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		gopher.QueueRemoveGopher()
	}
	tileMap.GopherWaitGroup.Done()
}

func (tileMap *GopherMapUpdater) Update() bool {

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

func (tileMap *GopherMapUpdater) processGophers() {

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

func (tileMap *GopherMapUpdater) processQueuedTasks() {
	tileMap.diagnostics.inputStopWatch.Start()
	tileMap.QueueableActions.Process()
	tileMap.diagnostics.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *GopherMapUpdater) TogglePause() {
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
