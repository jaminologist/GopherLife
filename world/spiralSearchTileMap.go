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
	QueueableActions
	FoodRespawnPickup
	GopherGeneration
	Insertable

	ActiveActors chan *GopherActor

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

	Statistics
	diagnostics Diagnostics
}

type SpiralTileSearch struct {
	TileContainer
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

	tileMap.ActiveActors = make(chan *GopherActor, statistics.MaximumNumberOfGophers*2)

	tileMap.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

	ag := GopherGeneration{
		Insertable:     tileMap.Insertable,
		maxGenerations: tileMap.Statistics.MaximumNumberOfGophers,
		ActiveGophers:  tileMap.ActiveActors,
		gopherArray:    tileMap.gopherArray,
	}

	tileMap.GopherGeneration = ag

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

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

		tileMap.gopherArray[i] = &gopher
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

	currentMapPoint, exists := tileMap.Tile(gopher.Position.GetX(), gopher.Position.GetY())

	if !exists {
		return false
	}

	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)
	tarGetTile, exists := tileMap.Tile(targetPosition.GetX(), targetPosition.GetY())

	if exists && tarGetTile.Gopher == nil {

		tarGetTile.Gopher = gopher
		currentMapPoint.Gopher = nil

		gopher.Position = targetPosition

		return true
	}

	return false
}

func (tileMap *SpiralSearchTileMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.gopherArray[rand.Intn(len(tileMap.gopherArray))]
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
	ActiveGophers  chan *GopherActor
	gopherArray    []*Gopher
}

func (gg *GopherGeneration) AddNewGopher(x int, y int, g *GopherActor) bool {

	if len(gg.gopherArray) <= gg.maxGenerations {
		if gg.InsertGopher(x, y, g.Gopher) {
			gg.ActiveGophers <- g
			return true
		}
	}

	return false
}

func (tileMap *SpiralSearchTileMap) onFoodPickUp(location calc.Coordinates) {

	size := 50
	xrange, yrange := rand.Perm(size), rand.Perm(size)
	food := NewPotato()

loop:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			newX, newY := location.GetX()+xrange[i]-size/2, location.GetY()+yrange[j]-size/2
			if tileMap.InsertFood(newX, newY, &food) {
				break loop
			}
		}
	}
}

func (tileMap *SpiralSearchTileMap) PerformEntityAction(gopher *Gopher, wg *sync.WaitGroup, channel chan *Gopher) {

	gopher.PerformMoment(tileMap)

	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		tileMap.QueueRemoveGopher(gopher)
	}

	wg.Done()
}

func (tileMap *SpiralSearchTileMap) Act(gopher *GopherActor, wg *sync.WaitGroup, channel chan *GopherActor) {
	gopher.Update()
	if !gopher.IsDecayed() {
		channel <- gopher
	} else {
		gopher.QueueRemoveGopher()
	}
	wg.Done()
}

func (tileMap *SpiralSearchTileMap) Update() bool {

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

func (tileMap *SpiralSearchTileMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

	numGophers := len(tileMap.ActiveActors)
	tileMap.gopherArray = make([]*Gopher, numGophers)
	tileMap.GopherGeneration.gopherArray = tileMap.gopherArray

	secondChannel := make(chan *GopherActor, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveActors
		tileMap.gopherArray[i] = gopher.Gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.Act(gopher, tileMap.GopherWaitGroup, secondChannel)

	}
	tileMap.ActiveActors = secondChannel
	tileMap.GopherGeneration.ActiveGophers = tileMap.ActiveActors
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *SpiralSearchTileMap) processQueuedTasks() {
	tileMap.diagnostics.inputStopWatch.Start()
	tileMap.QueueableActions.Process()
	tileMap.diagnostics.inputStopWatch.Stop()
}

//TogglePause Toggles the pause
func (tileMap *SpiralSearchTileMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (tileMap *SpiralSearchTileMap) QueueRemoveGopher(gopher *Gopher) {

	tileMap.Add(func() {
		//gopher = nil
		if mapPoint, ok := tileMap.Tile(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			mapPoint.Gopher = nil
		}
	})
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (tileMap *SpiralSearchTileMap) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	tileMap.Add(func() {
		success := tileMap.MoveGopher(gopher, moveX, moveY)
		_ = success
	})

}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (tileMap *SpiralSearchTileMap) QueuePickUpFood(gopher *Gopher) {

	tileMap.Add(func() {
		food, ok := tileMap.RemoveFood(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			tileMap.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})
}

func (tileMap *SpiralSearchTileMap) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	tileMap.Add(func() {

		if mapPoint, ok := tileMap.Tile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(tileMap.Statistics.GopherBirthRate)

			emptySpaces := tileMap.Search(gopher.Position, 10, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])

						var gopherActor = GopherActor{
							Gopher:           &newborn,
							GopherBirthRate:  tileMap.GopherBirthRate,
							QueueableActions: tileMap,
							Searchable:       tileMap.Searchable,
							TileContainer:    tileMap.TileContainer,
							Insertable:       tileMap,
							PickableTiles:    &tileMap.FoodRespawnPickup,
							MoveableActors:   tileMap,
							ActorGeneration:  &tileMap.GopherGeneration,
						}

						tileMap.GopherGeneration.AddNewGopher(pos.GetX(), pos.GetY(), &gopherActor)
					}

				}

			}

		}
	})

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

		/*if coordinates.X == 0 && coordinates.Y == 0 {
			continue
		}*/

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
