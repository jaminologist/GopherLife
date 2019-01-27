package world

import (
	"errors"
	"gopherlife/calc"
	"gopherlife/names"
	"math/rand"
	"sync"
)

//PartitionTileMap is cool

const (
	gridWidth  = 5
	gridHeight = 5
)

type PartitionTileMap struct {
	grid [][]*GridSection

	gridWidth  int
	gridHeight int

	numberOfGridsWide   int
	numberOfGridsHeight int

	actionQueue   chan func()
	ActiveGophers chan *Gopher

	GopherWaitGroup *sync.WaitGroup
	SelectedGopher  *Gopher
	gopherArray     []*Gopher
	Moments         int
	IsPaused        bool

	Statistics
	diagnostics Diagnostics
}

func CreatePartitionTileMapCustom(statistics Statistics) *PartitionTileMap {

	tileMap := PartitionTileMap{}
	tileMap.Statistics = statistics
	tileMap.actionQueue = make(chan func(), statistics.MaximumNumberOfGophers*2)

	tileMap.numberOfGridsWide = statistics.Width / gridWidth

	if tileMap.numberOfGridsWide*gridWidth < statistics.Width {
		tileMap.numberOfGridsWide++
	}

	tileMap.numberOfGridsHeight = statistics.Height / gridHeight

	if tileMap.numberOfGridsWide*gridHeight < statistics.Height {
		tileMap.numberOfGridsHeight++
	}

	tileMap.grid = make([][]*GridSection, tileMap.numberOfGridsWide)
	tileMap.gridWidth = gridWidth
	tileMap.gridHeight = gridHeight

	for i := 0; i < tileMap.numberOfGridsWide; i++ {
		tileMap.grid[i] = make([]*GridSection, tileMap.numberOfGridsHeight)

		for j := 0; j < tileMap.numberOfGridsHeight; j++ {
			gridSection := NewGridSection(i*gridWidth, j*gridHeight, gridWidth, gridHeight)
			tileMap.grid[i][j] = &gridSection
		}
	}

	tileMap.ActiveGophers = make(chan *Gopher, statistics.NumberOfGophers)
	tileMap.gopherArray = make([]*Gopher, statistics.NumberOfGophers)

	var wg sync.WaitGroup
	tileMap.GopherWaitGroup = &wg

	tileMap.setUpTiles()
	return &tileMap
}

func CreatePartitionTileMap() *PartitionTileMap {
	tileMap := CreatePartitionTileMapCustom(
		Statistics{
			Width:                  3000,
			Height:                 3000,
			NumberOfGophers:        5000,
			NumberOfFood:           50000,
			MaximumNumberOfGophers: 100000,
			GopherBirthRate:        7,
		},
	)
	return tileMap
}

func (tileMap *PartitionTileMap) setUpTiles() {

	keys := calc.GenerateRandomizedCoordinateArray(0, 0,
		tileMap.Statistics.Width, tileMap.Statistics.Height)

	count := 0

	for i := 0; i < tileMap.Statistics.NumberOfGophers; i++ {

		pos := keys[count]
		var gopher = NewGopher(names.CuteName(), pos)

		tileMap.InsertGopher(&gopher, pos.GetX(), pos.GetY())

		if i == 0 {
			tileMap.SelectedGopher = &gopher
		}

		tileMap.gopherArray[i] = &gopher
		tileMap.ActiveGophers <- &gopher
		count++
	}

	for i := 0; i < tileMap.Statistics.NumberOfFood; i++ {
		pos := keys[count]
		var food = NewPotato()
		tileMap.InsertFood(&food, pos.GetX(), pos.GetY())
		count++
	}
}

//TogglePause Toggles the pause
func (tileMap *PartitionTileMap) TogglePause() {
	tileMap.IsPaused = !tileMap.IsPaused
}

func (tileMap *PartitionTileMap) InsertGopher(gopher *Gopher, x int, y int) bool {
	if gridSection, ok := tileMap.getGrid(x, y); ok {
		gopher.Position.X = x
		gopher.Position.Y = y
		err := gridSection.InsertGopher(x, y, gopher)

		if err == nil {
			return true
		}
	}
	return false
}

func (tileMap *PartitionTileMap) InsertFood(food *Food, x int, y int) {
	if gridSection, ok := tileMap.getGrid(x, y); ok {
		food.Position.X = x
		food.Position.Y = y
		gridSection.InsertFood(x, y, food)
	}
}

func (tileMap *PartitionTileMap) Update() bool {

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
	tileMap.Statistics.NumberOfGophers = len(tileMap.ActiveGophers)
	if tileMap.Statistics.NumberOfGophers > 0 {
		tileMap.Moments++
	}

	tileMap.diagnostics.processStopWatch.Stop()

	return true

}

func (tileMap *PartitionTileMap) processGophers() {

	tileMap.diagnostics.gopherStopWatch.Start()

	numGophers := len(tileMap.ActiveGophers)
	tileMap.gopherArray = make([]*Gopher, numGophers)

	secondChannel := make(chan *Gopher, numGophers*2)
	for i := 0; i < numGophers; i++ {
		gopher := <-tileMap.ActiveGophers
		tileMap.gopherArray[i] = gopher
		tileMap.GopherWaitGroup.Add(1)
		go tileMap.performEntityAction(gopher, secondChannel)

	}
	tileMap.ActiveGophers = secondChannel
	tileMap.GopherWaitGroup.Wait()

	tileMap.diagnostics.gopherStopWatch.Stop()
}

func (tileMap *PartitionTileMap) performEntityAction(gopher *Gopher, channel chan *Gopher) {

	gopher.PerformMoment(tileMap)

	if !gopher.IsDecayed() {

		wait := true

		for wait {
			select {
			case channel <- gopher:
				wait = false
			default:
				//	fmt.Println("Can't Write")
			}
		}

	} else {
		tileMap.QueueRemoveGopher(gopher)
	}

	tileMap.GopherWaitGroup.Done()

}

func (tileMap *PartitionTileMap) processQueuedTasks() {

	tileMap.diagnostics.inputStopWatch.Start()

	wait := true
	for wait {
		select {
		case action := <-tileMap.actionQueue:
			action()
		default:
			wait = false
		}
	}

	tileMap.diagnostics.inputStopWatch.Stop()
}

func (tileMap *PartitionTileMap) Tile(x int, y int) (*Tile, bool) {
	if x < 0 || x >= tileMap.Statistics.Width || y < 0 || y >= tileMap.Statistics.Height {
		return nil, false
	}

	if gridSection, ok := tileMap.getGrid(x, y); ok {

		if tile, ok := gridSection.GetMapPointFromGrid(x, y); ok {
			return tile, ok
		}
		tilevar := Tile{nil, nil}
		return &tilevar, true
	}
	return nil, false
}

func (tileMap *PartitionTileMap) SelectedTile() (*Tile, bool) {

	if tileMap.SelectedGopher != nil {
		if tile, ok := tileMap.Tile(tileMap.SelectedGopher.Position.GetX(), tileMap.SelectedGopher.Position.GetY()); ok {
			return tile, ok
		}
	}
	return nil, false

}

//SelectEntity Uses the given co-ordinates to select and return a gopher in the tileMap
//If there is not a gopher at the give coordinates this function returns zero.
func (tileMap *PartitionTileMap) SelectEntity(x int, y int) (*Gopher, bool) {

	tileMap.SelectedGopher = nil

	if mapPoint, ok := tileMap.Tile(x, y); ok {
		if mapPoint.Gopher != nil {
			tileMap.SelectedGopher = mapPoint.Gopher
			return mapPoint.Gopher, true
		}
	}

	return nil, true
}

func (tileMap *PartitionTileMap) SelectRandomGopher() {
	tileMap.SelectedGopher = tileMap.gopherArray[rand.Intn(len(tileMap.gopherArray))]
}

func (tileMap *PartitionTileMap) UnSelectGopher() {
	tileMap.SelectedGopher = nil
}

func (tileMap *PartitionTileMap) Stats() *Statistics {
	return &tileMap.Statistics
}

func (tileMap *PartitionTileMap) Diagnostics() *Diagnostics {
	return &tileMap.diagnostics
}

//QueueGopherMove Adds the Move Gopher Method to the Input Queue.
func (tileMap *PartitionTileMap) QueueGopherMove(gopher *Gopher, moveX int, moveY int) {

	tileMap.addFunctionToWorldInputActions(func() {
		success := tileMap.MoveGopher(gopher, moveX, moveY)
		_ = success
	})
}

func (tileMap *PartitionTileMap) MoveGopher(gopher *Gopher, moveX int, moveY int) bool {

	currentPosition := gopher.Position
	targetPosition := gopher.Position.RelativeCoordinate(moveX, moveY)

	grid, exists := tileMap.getGrid(currentPosition.GetX(), currentPosition.GetY())

	if !exists {
		return false
	}

	var targetGrid *GridSection
	targetX, targetY := targetPosition.GetX(), targetPosition.GetY()

	if grid.Contains(targetX, targetY) {
		targetGrid = grid
	} else {
		if val, ok := tileMap.getGrid(targetX, targetY); ok {
			targetGrid = val
		} else {
			return false
		}
	}

	err := targetGrid.InsertGopher(targetX, targetY, gopher)

	if err == nil {
		grid.RemoveGopher(gopher)
		gopher.Position.Set(targetX, targetY)
		return true
	}
	return false
}

//QueuePickUpFood Adds the PickUp Food Method to the Input Queue. If food is at the give position it is added to the Gopher's
//held food variable
func (tileMap *PartitionTileMap) QueuePickUpFood(gopher *Gopher) {
	tileMap.addFunctionToWorldInputActions(func() {
		food, ok := tileMap.removeFoodFromWorld(gopher.Position.GetX(), gopher.Position.GetY())
		if ok {
			gopher.HeldFood = food
			tileMap.onFoodPickUp(gopher.Position)
			gopher.ClearFoodTargets()
		}
	})
}

func (tileMap *PartitionTileMap) removeFoodFromWorld(x int, y int) (*Food, bool) {

	if grid, mapPoint, ok := tileMap.getGridAndMapPoint(x, y); ok {
		if mapPoint.Food != nil {
			food := mapPoint.Food
			grid.RemoveFood(mapPoint.Food)
			return food, true
		}
	}

	return nil, false
}

func (tileMap *PartitionTileMap) onFoodPickUp(location calc.Coordinates) {

	size := 50

	xrange := rand.Perm(size)
	yrange := rand.Perm(size)

loop:
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {

			newX := location.GetX() + xrange[i] - size/2
			newY := location.GetY() + yrange[j] - size/2

			if grid, ok := tileMap.getGrid(newX, newY); ok {
				food := NewPotato()
				mp, exists := grid.GetMapPointFromGrid(newX, newY)

				if !exists || mp.Food == nil {
					grid.InsertFood(newX, newY, &food)
					break loop
				}
			}
		}
	}

}

func (tileMap *PartitionTileMap) QueueMating(gopher *Gopher, matePosition calc.Coordinates) {

	tileMap.addFunctionToWorldInputActions(func() {

		if mapPoint, ok := tileMap.Tile(matePosition.GetX(), matePosition.GetY()); ok && mapPoint.HasGopher() {

			mate := mapPoint.Gopher
			litterNumber := rand.Intn(tileMap.Statistics.GopherBirthRate)

			emptySpaces := tileMap.Search(gopher.Position, 10, litterNumber, SearchForEmptySpace)

			if mate.Gender == Female && len(emptySpaces) > 0 {
				mate.IsMated = true
				mate.CounterTillReadyToFindLove = 0

				for i := 0; i < litterNumber; i++ {

					if i < len(emptySpaces) {
						pos := emptySpaces[i]
						newborn := NewGopher(names.CuteName(), emptySpaces[i])

						if len(tileMap.gopherArray) <= tileMap.Statistics.MaximumNumberOfGophers {
							if tileMap.InsertGopher(&newborn, pos.GetX(), pos.GetY()) {
								tileMap.ActiveGophers <- &newborn
							}
						}
					}

				}

			}

		}
	})

}

//QueueRemoveGopher Adds the Remove Gopher Method to the Input Queue.
func (tileMap *PartitionTileMap) QueueRemoveGopher(gopher *Gopher) {

	tileMap.addFunctionToWorldInputActions(func() {
		if gridSection, ok := tileMap.getGrid(gopher.Position.GetX(), gopher.Position.GetY()); ok {
			gridSection.RemoveGopher(gopher)
		}
	})
}

func (tileMap *PartitionTileMap) Search(startPosition calc.Coordinates, radius int, maximumFind int, searchType SearchType) []calc.Coordinates {
	var coordsArray = []calc.Coordinates{}

	spiral := calc.NewSpiral(radius, radius)

	var query TileQuery

	switch searchType {
	case SearchForFood:
		locations := queryForFood(tileMap, radius, startPosition.GetX(), startPosition.GetY())
		calc.SortByNearestFromCoordinate(startPosition, locations)
		return locations
	case SearchForEmptySpace:
		query = CheckMapPointForEmptySpace
	case SearchForFemaleGopher:
		locations := queryForFemalePartner(tileMap, radius, startPosition.GetX(), startPosition.GetY())
		calc.SortByNearestFromCoordinate(startPosition, locations)
		return locations
	}

	for {

		coordinates, hasNext := spiral.Next()

		if hasNext == false || len(coordsArray) > maximumFind {
			break
		}

		relativeCoords := startPosition.RelativeCoordinate(coordinates.X, coordinates.Y)

		if tile, ok := tileMap.Tile(relativeCoords.GetX(), relativeCoords.GetY()); ok {
			if query(tile) {
				coordsArray = append(coordsArray, relativeCoords)
			}
		}
	}

	calc.SortByNearestFromCoordinate(startPosition, coordsArray)

	return coordsArray
}

func queryForFood(tileMap *PartitionTileMap, size int, x int, y int) []calc.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-size, y-size, x+size, y+size

	startX, startY := tileMap.convertToGridCoordinates(x-size, y-size)
	endX, endY := tileMap.convertToGridCoordinates(x+size, y+size)

	foodLocations := make([]calc.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {

			if grid, ok := tileMap.getGridConvertedCoordinates(x, y); ok {
				for key := range grid.foodTileLocations {

					tile := grid.foodTileLocations[key]

					i, j := tile.Food.Position.GetX(), tile.Food.Position.GetY()
					if i >= worldStartX &&
						i < worldEndX &&
						j >= worldStartY &&
						j < worldEndY {
						foodLocations = append(foodLocations, tile.Food.Position)
					}
				}
			}

		}
	}

	return foodLocations
}

func queryForFemalePartner(tileMap *PartitionTileMap, size int, x int, y int) []calc.Coordinates {

	worldStartX, worldStartY, worldEndX, worldEndY := x-size, y-size, x+size, y+size

	startX, startY := tileMap.convertToGridCoordinates(x-size, y-size)
	endX, endY := tileMap.convertToGridCoordinates(x+size, y+size)

	locations := make([]calc.Coordinates, 0)

	for x := startX; x <= endX; x++ {
		for y := startY; y <= endY; y++ {

			if grid, ok := tileMap.getGridConvertedCoordinates(x, y); ok {
				for key := range grid.gopherTileLocations {

					tile := grid.gopherTileLocations[key]

					i, j := tile.Gopher.Position.GetX(), tile.Gopher.Position.GetY()
					if i >= worldStartX &&
						i < worldEndX &&
						j >= worldStartY &&
						j < worldEndY && tile.Gopher.Gender == Female && tile.Gopher.IsLookingForLove() {
						locations = append(locations, tile.Gopher.Position)
					}
				}
			}

		}
	}

	return locations
}

//AddFunctionToWorldInputActions is used to store functions that write data to the tileMap.
func (tileMap *PartitionTileMap) addFunctionToWorldInputActions(inputFunction func()) {
	tileMap.actionQueue <- inputFunction
}

func (tileMap *PartitionTileMap) getGrid(x int, y int) (*GridSection, bool) {
	gridX, gridY := tileMap.convertToGridCoordinates(x, y)
	return tileMap.getGridConvertedCoordinates(gridX, gridY)
}

func (tileMap *PartitionTileMap) getGridAndMapPoint(x int, y int) (*GridSection, *Tile, bool) {

	if gridSection, ok := tileMap.getGrid(x, y); ok {
		mp, ok := gridSection.GetMapPointFromGrid(x, y)
		return gridSection, mp, ok
	}

	return nil, nil, false
}

func (tileMap *PartitionTileMap) getGridConvertedCoordinates(x int, y int) (*GridSection, bool) {

	if x < 0 || y < 0 || x >= tileMap.numberOfGridsWide || y >= tileMap.numberOfGridsHeight {
		return nil, false
	}

	val := tileMap.grid[x][y]

	if val != nil {
		return val, true
	}

	return nil, false

}

func (tileMap *PartitionTileMap) convertToGridCoordinates(x int, y int) (int, int) {
	gridX, gridY := x/tileMap.gridWidth, y/tileMap.gridHeight
	return gridX, gridY
}

type GridSection struct {
	x int
	y int

	width  int
	height int

	grid [][]*Tile

	gopherTileLocations map[int]*Tile
	foodTileLocations   map[int]*Tile

	InputActions chan func()
}

//NewGridSection Creates a new Grid Section with the given position
//and size
func NewGridSection(x int, y int, width int, height int) GridSection {

	twoD := make([][]*Tile, width)

	for i := 0; i < width; i++ {
		twoD[i] = make([]*Tile, height)
	}

	return GridSection{
		x:                   x,
		y:                   y,
		width:               width,
		height:              height,
		grid:                twoD,
		gopherTileLocations: make(map[int]*Tile),
		foodTileLocations:   make(map[int]*Tile),
		InputActions:        make(chan func(), width*height*2)}
}

func (gridSection *GridSection) GetPosition() (int, int) {
	return gridSection.x, gridSection.y
}

func (gridSection *GridSection) Tile(x int, y int) (*Tile, bool) {

	gridX, gridY := (x - gridSection.x), (y - gridSection.y)

	if gridX < 0 || gridY < 0 || gridX > gridSection.width-1 || gridY > gridSection.height-1 {
		return nil, false
	}

	val := gridSection.grid[gridX][gridY]

	if val == nil {
		return nil, false
	}

	return val, true
}

func (gridSection *GridSection) SetMapPoint(x int, y int, mp *Tile) {

	x, y = (x - gridSection.x), (y - gridSection.y)

	gridSection.grid[x][y] = mp
}

func (gridSection *GridSection) convertToGridCoordinates(x int, y int) (gridX int, gridY int) {
	return (x - gridSection.x), (y - gridSection.y)
}

func (gridSection *GridSection) InsertGopher(x int, y int, gopher *Gopher) error {

	if gridSection.Contains(x, y) {

		tile, ok := gridSection.Tile(x, y)
		if ok && !tile.HasGopher() {
			tile.Gopher = gopher
			gridSection.addGopherTileLocation(x, y, tile)
		} else if !ok {
			newTile := NewTile(gopher, nil)
			gridSection.SetMapPoint(x, y, &newTile)
			gridSection.addGopherTileLocation(x, y, &newTile)
		} else {
			return errors.New("MapPoint already contains Gopher, Gopher can not be inserted")
		}

	} else {
		return errors.New("Gopher inserted outside of grid section bounds")
	}

	return nil

}

func (gridSection *GridSection) addGopherTileLocation(x int, y int, tile *Tile) {
	x, y = gridSection.convertToGridCoordinates(x, y)
	gridSection.gopherTileLocations[calc.Hashcode(x, y)] = tile
}

func (gridSection *GridSection) addFoodTileLocations(x int, y int, tile *Tile) {
	x, y = gridSection.convertToGridCoordinates(x, y)
	gridSection.foodTileLocations[calc.Hashcode(x, y)] = tile
}

func (gridSection *GridSection) deleteGopherTileLocation(x int, y int) {
	x, y = gridSection.convertToGridCoordinates(x, y)
	delete(gridSection.gopherTileLocations, calc.Hashcode(x, y))
}

func (gridSection *GridSection) deleteFoodTileLocation(x int, y int) {
	x, y = gridSection.convertToGridCoordinates(x, y)
	delete(gridSection.foodTileLocations, calc.Hashcode(x, y))
}

func (gridSection *GridSection) InsertFood(x int, y int, food *Food) error {

	food.Position.X = x
	food.Position.Y = y

	if gridSection.Contains(x, y) {

		tile, ok := gridSection.Tile(x, y)
		if ok && !tile.HasFood() {
			tile.Food = food
			gridSection.addFoodTileLocations(x, y, tile)
		} else if !ok {
			newTile := NewTile(nil, food)
			gridSection.SetMapPoint(x, y, &newTile)
			gridSection.addFoodTileLocations(x, y, &newTile)
		} else {
			return errors.New("Tile already contains Food, Food can not be inserted")
		}

	} else {
		return errors.New("Food inserted outside of grid section bounds")
	}

	return nil

}

func (gridSection *GridSection) RemoveGopher(gopher *Gopher) {

	x, y := gopher.Position.GetX(), gopher.Position.GetY()

	tile, ok := gridSection.Tile(x, y)
	if ok {
		tile.ClearGopher()
		gridSection.deleteGopherTileLocation(x, y)
		if tile.isEmpty() {
			tile = nil
		}
	}
}

func (gridSection *GridSection) RemoveFood(food *Food) {

	x, y := food.Position.GetX(), food.Position.GetY()

	tile, ok := gridSection.Tile(x, y)
	if ok {
		tile.ClearFood()
		gridSection.deleteFoodTileLocation(x, y)
		if tile.isEmpty() {
			tile = nil
		}
	}
}

func (gridsection *GridSection) Contains(x int, y int) bool {

	return x >= gridsection.x &&
		x < gridsection.x+gridsection.width &&
		y >= gridsection.y &&
		y < gridsection.y+gridsection.height

}

//GetMapPointFromGrid returns a map point of the given coordinates if the map point exists
//inside the grid
func (gridsection *GridSection) GetMapPointFromGrid(x int, y int) (*Tile, bool) {

	tile, ok := gridsection.Tile(x, y)
	if ok {
		return tile, true
	}

	return nil, false
}
