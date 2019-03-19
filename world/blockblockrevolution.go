package world

import (
	"gopherlife/colors"
	"gopherlife/geometry"
	"gopherlife/timer"
	"image/color"
	"math/rand"
	"time"
)

type BlockBlockRevolutionSettings struct {
	Dimensions
	BlockSpeedReduction int
}

type BlockBlockRevolutionMap struct {
	grid [][]*BlockBlockRevolutionTile
	BlockBlockRevolutionSettings
	Container
	CurrentTetromino Tetromino

	newBlockFunctions []func(int, int, BlockInserterAndRemover) (Tetromino, bool)

	nextNewBlockFunctions []func(int, int, BlockInserterAndRemover) (Tetromino, bool)

	ActionQueuer

	FrameTimer timer.StopWatch
	FrameSpeed time.Duration

	DownToNextLineCount int

	Score      int
	IsGameOver bool
}

func NewBlockBlockRevolutionMap(settings BlockBlockRevolutionSettings) BlockBlockRevolutionMap {

	r := geometry.NewRectangle(0, 0, settings.Width, settings.Height)

	grid := make([][]*BlockBlockRevolutionTile, settings.Width)

	for i := 0; i < settings.Width; i++ {
		grid[i] = make([]*BlockBlockRevolutionTile, settings.Height)

		for j := 0; j < settings.Height; j++ {
			tile := BlockBlockRevolutionTile{
				Coordinates: geometry.Coordinates{
					X: i,
					Y: j,
				},
			}
			grid[i][j] = &tile
		}
	}

	qa := NewFiniteActionQueue(1)

	bbrm := BlockBlockRevolutionMap{
		Container:                    &r,
		BlockBlockRevolutionSettings: settings,
		ActionQueuer:                 &qa,
		grid:                         grid,
		FrameSpeed:                   5,
		newBlockFunctions: []func(int, int, BlockInserterAndRemover) (Tetromino, bool){
			NewSquareTetrominoes,
			NewLTetrominoes,
			NewJTetrominoes,
			NewSTetrominoes,
			NewZTetrominoes,
			NewTTetrominoes,
			NewITetrominoes,
		},
	}

	bbrm.AddNewBlock()

	return bbrm

}

func (bbrm *BlockBlockRevolutionMap) Update() bool {

	if bbrm.IsGameOver {
		return false
	}

	bbrm.FrameTimer.Start()
	bbrm.Process()

	if bbrm.DownToNextLineCount < 5 {
		bbrm.DownToNextLineCount++
	} else {
		bbrm.DownToNextLineCount = 0
		if !bbrm.MoveCurrentTetrominoDown() {
			//bbrm.CurrentTetromino = nil
			bbrm.CheckForAndClearLines()

			if !bbrm.AddNewBlock() {
				bbrm.IsGameOver = true
			}
		}
	}

	for bbrm.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*time.Duration(bbrm.BlockBlockRevolutionSettings.BlockSpeedReduction) {
	}

	return true
}

func (bbrm *BlockBlockRevolutionMap) Tile(x int, y int) (*BlockBlockRevolutionTile, bool) {
	if bbrm.Contains(x, y) {
		return bbrm.grid[x][y], true
	}
	return nil, false
}

func (bbrm *BlockBlockRevolutionMap) MoveCurrentTetrominoDown() bool {
	return bbrm.MoveCurrentTetromino(0, -1)
}

func (bbrm *BlockBlockRevolutionMap) MoveCurrentTetrominoLeft() bool {
	return bbrm.MoveCurrentTetromino(-1, 0)
}

func (bbrm *BlockBlockRevolutionMap) MoveCurrentTetrominoRight() bool {
	return bbrm.MoveCurrentTetromino(1, 0)
}

func (bbrm *BlockBlockRevolutionMap) InstantDown() {
	for bbrm.MoveCurrentTetrominoDown() {
	}
}

func (bbrm *BlockBlockRevolutionMap) RotateTetromino() {
	bbrm.Add(func() {
		bbrm.CurrentTetromino.Rotate()
	})
}

func (bbrm *BlockBlockRevolutionMap) MoveCurrentTetromino(moveX int, moveY int) bool {
	blocks := bbrm.CurrentTetromino.Blocks()

	for i := 0; i < len(blocks); i++ {
		bbrm.RemoveBlock(blocks[i].GetX(), blocks[i].GetY())
	}

	for i := 0; i < len(blocks); i++ {
		newX, newY := blocks[i].GetX()+moveX, blocks[i].GetY()+moveY
		if _, ok := bbrm.ContainsBlock(newX, newY); ok || !bbrm.Contains(newX, newY) {

			for i := 0; i < len(blocks); i++ {
				bbrm.InsertBlock(blocks[i].GetX(), blocks[i].GetY(), blocks[i])
			}

			return false
		}
	}

	for i := 0; i < len(blocks); i++ {
		newX, newY := blocks[i].GetX()+moveX, blocks[i].GetY()+moveY
		bbrm.InsertBlock(newX, newY, blocks[i])
	}

	return true
}

func (bbrm *BlockBlockRevolutionMap) ContainsBlock(x int, y int) (*Block, bool) {
	if tile, ok := bbrm.Tile(x, y); ok {
		if tile.Block != nil {
			return tile.Block, true
		}
	}
	return nil, false
}

func (bbrm *BlockBlockRevolutionMap) CheckForAndClearLines() {

	//linesToClear := make([]int, bbrm.Height)

	scoreMultiplier := 0

	for y := 0; y < bbrm.Height; y++ {

		canAddLine := true

		for x := 0; x < bbrm.Width; x++ {
			tile := bbrm.grid[x][y]
			if !tile.ContainsBlock() {
				canAddLine = false
				break
			}

		}

		if canAddLine {
			bbrm.RemoveAllBlocksFromLine(y)
			bbrm.ShiftAllBlocksAboveLineDown(y)
			scoreMultiplier++
			y--
		}

	}

	for i := 0; i < scoreMultiplier; i++ {
		bbrm.Score += scoreMultiplier * 100
	}

}

func (bbrm *BlockBlockRevolutionMap) RemoveAllBlocksFromLine(line int) {
	for i := 0; i < bbrm.Width; i++ {
		bbrm.grid[i][line].RemoveBlock()
	}
}

func (bbrm *BlockBlockRevolutionMap) ShiftAllBlocksAboveLineDown(line int) {

	blocks := make([]*Block, 0)

	for x := 0; x < bbrm.Width; x++ {
		for y := line; y < bbrm.Height; y++ {
			tile := bbrm.grid[x][y]
			if tile.ContainsBlock() {
				blocks = append(blocks, tile.Block)
				tile.RemoveBlock()
			}
		}
	}

	for _, block := range blocks {
		bbrm.InsertBlock(block.GetX(), block.GetY()-1, block)
	}

}

func (bbrm *BlockBlockRevolutionMap) InsertBlock(x int, y int, b *Block) bool {
	if tile, ok := bbrm.Tile(x, y); ok {
		return tile.InsertBlock(b)
	}
	return false
}

func (bbrm *BlockBlockRevolutionMap) RemoveBlock(x int, y int) {
	if tile, ok := bbrm.Tile(x, y); ok {
		tile.RemoveBlock()
	}
}

func (bbrm *BlockBlockRevolutionMap) AddNewBlock() bool {

	if len(bbrm.nextNewBlockFunctions) == 0 {
		for i := 0; i < 3; i++ {
			for _, newblockfunc := range rand.Perm(len(bbrm.newBlockFunctions)) {
				bbrm.nextNewBlockFunctions = append(bbrm.nextNewBlockFunctions, bbrm.newBlockFunctions[newblockfunc])
			}
		}
	}

	x, y := bbrm.nextNewBlockFunctions[0], bbrm.nextNewBlockFunctions[1:]
	bbrm.nextNewBlockFunctions = y

	block, ok := x(bbrm.Width/2, bbrm.Height-2, bbrm)

	if ok {
		bbrm.CurrentTetromino = block
		return true
	}

	return false

}

type BlockContainer interface {
	Container
	ContainsBlock(x int, y int) (*Block, bool)
}

type BlockInserterAndRemover interface {
	BlockContainer
	InsertBlock(x int, y int, b *Block) bool
	RemoveBlock(x int, y int)
}

type BlockBlockRevolutionTile struct {
	geometry.Coordinates
	Block *Block
}

func (bbrt *BlockBlockRevolutionTile) InsertBlock(b *Block) bool {

	if bbrt.Block == nil {
		b.SetXY(bbrt.GetX(), bbrt.GetY())
		bbrt.Block = b
		return true
	}
	return false
}

func (bbrt *BlockBlockRevolutionTile) RemoveBlock() {
	bbrt.Block = nil
}

func (bbrt *BlockBlockRevolutionTile) ContainsBlock() bool {
	return bbrt.Block != nil
}

type Tetromino interface {
	Rotate()
	Blocks() []*Block
}

type BlockHolder struct {
	blocks []*Block
}

func (b *BlockHolder) Blocks() []*Block {
	return b.blocks
}

type Block struct {
	geometry.Coordinates
	Color color.RGBA
}

type SquareTetrominoes struct {
	BlockInserterAndRemover
	blocks []*Block
	geometry.Coordinates
}

func InsertAllBlocks(blocks []*Block, bir BlockInserterAndRemover) {
	for _, block := range blocks {
		bir.InsertBlock(block.GetX(), block.GetY(), block)
	}
}

func RemoveAllBlocks(blocks []*Block, bir BlockInserterAndRemover) {
	for _, block := range blocks {
		bir.RemoveBlock(block.GetX(), block.GetY())
	}
}

func SetColorOfBlocks(blocks []*Block, c color.RGBA) {
	for _, block := range blocks {
		block.Color = c
	}
}

func NewCenterLeftBlockPositionUsingDirection(d geometry.Direction) geometry.Coordinates {

	switch d {
	case geometry.Up:
		return geometry.NewCoordinate(1, 0)
	case geometry.Left:
		return geometry.NewCoordinate(0, 1)
	case geometry.Down:
		return geometry.NewCoordinate(-1, 0)
	case geometry.Right:
		return geometry.NewCoordinate(0, -1)
	default:
		panic("Invalid geometry.Direction Used")
	}

}

func NewCenterRightBlockPositionUsingDirection(d geometry.Direction) geometry.Coordinates {
	d = d.TurnClockWise90().TurnClockWise90()
	return NewCenterLeftBlockPositionUsingDirection(d)
}

func NewCenterUpBlockPositionUsingDirection(d geometry.Direction) geometry.Coordinates {
	d = d.TurnClockWise90()
	return NewCenterLeftBlockPositionUsingDirection(d)
}

func NewCenterDownBlockPositionUsingDirection(d geometry.Direction) geometry.Coordinates {
	d = d.TurnAntiClockWise90()
	return NewCenterLeftBlockPositionUsingDirection(d)
}

func NewBlock(x int, y int) *Block {
	b := Block{
		Coordinates: geometry.Coordinates{x, y},
	}
	return &b
}

func CanTetrominoFit(cs []geometry.Coordinates, bc BlockContainer) bool {
	for _, coords := range cs {
		if _, ok := bc.ContainsBlock(coords.GetX(), coords.GetY()); ok || !bc.Contains(coords.GetX(), coords.GetY()) {
			return false
		}
	}
	return true
}

func NewSquareTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	topLeftBlock := &Block{}
	topRightBlock := &Block{}
	bottomLeftBlock := &Block{}
	bottomRightBlock := &Block{}

	blocks := []*Block{topLeftBlock, topRightBlock, bottomLeftBlock, bottomRightBlock}
	SetColorOfBlocks(blocks, colors.Yellow)

	if !bir.InsertBlock(x, y, topLeftBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x+1, y, topRightBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x, y-1, bottomLeftBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x+1, y-1, bottomRightBlock) {
		return nil, false
	}

	sqt := SquareTetrominoes{
		blocks: blocks,
	}

	return &sqt, true
}

func (s *SquareTetrominoes) Rotate() {

}

func (s *SquareTetrominoes) Blocks() []*Block {
	return s.blocks
}

type LTetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerBlock   *Block
	centerLeft    *Block
	centerRight   *Block
	centerRightUp *Block

	rotateDirection geometry.Direction
}

//NewLTetrominoes Creates a new L-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewLTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	lt := LTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerLeft:              NewBlock(x-1, y),
		centerRight:             NewBlock(x+1, y),
		centerRightUp:           NewBlock(x+1, y+1),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	lt.blocks = []*Block{lt.centerBlock, lt.centerLeft, lt.centerRight, lt.centerRightUp}
	SetColorOfBlocks(lt.Blocks(), colors.Orange)

	for _, block := range lt.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &lt, true

}

func (l *LTetrominoes) Rotate() {

	x, y := l.centerBlock.GetX(), l.centerBlock.GetY()

	l.rotateDirection = l.rotateDirection.TurnAntiClockWise90()

	newCenter := geometry.NewCoordinate(x, y)
	newCenterLeft := geometry.Add(newCenter, NewCenterLeftBlockPositionUsingDirection(l.rotateDirection))
	newCenterRight := geometry.Add(newCenter, NewCenterRightBlockPositionUsingDirection(l.rotateDirection))
	newCenterRightUp := geometry.Add(newCenter, geometry.Add(NewCenterRightBlockPositionUsingDirection(l.rotateDirection), NewCenterUpBlockPositionUsingDirection(l.rotateDirection)))

	newCoordinateSlice := []geometry.Coordinates{newCenter, newCenterLeft, newCenterRight, newCenterRightUp}

	for _, block := range l.blocks {
		l.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := CanTetrominoFit(newCoordinateSlice, l.BlockInserterAndRemover)

	if canRotate {
		l.InsertBlock(newCenter.GetX(), newCenter.GetY(), l.centerBlock)
		l.InsertBlock(newCenterLeft.GetX(), newCenterLeft.GetY(), l.centerLeft)
		l.InsertBlock(newCenterRight.GetX(), newCenterRight.GetY(), l.centerRight)
		l.InsertBlock(newCenterRightUp.GetX(), newCenterRightUp.GetY(), l.centerRightUp)
	} else {
		l.rotateDirection = l.rotateDirection.TurnClockWise90()
	}

}

type JTetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerBlock       *Block
	centerLeftBlock   *Block
	centerRightBlock  *Block
	centerLeftUpBlock *Block

	rotateDirection geometry.Direction
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewJTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	jt := JTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerLeftBlock:         NewBlock(x-1, y),
		centerRightBlock:        NewBlock(x+1, y),
		centerLeftUpBlock:       NewBlock(x-1, y+1),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	jt.blocks = []*Block{jt.centerBlock, jt.centerLeftBlock, jt.centerRightBlock, jt.centerLeftUpBlock}
	SetColorOfBlocks(jt.Blocks(), colors.MingBlue)

	for _, block := range jt.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &jt, true
}

func (jt *JTetrominoes) Rotate() {

	x, y := jt.centerBlock.GetX(), jt.centerBlock.GetY()

	jt.rotateDirection = jt.rotateDirection.TurnAntiClockWise90()
	newCenter := geometry.NewCoordinate(x, y)
	newCenterLeft := geometry.Add(newCenter, NewCenterLeftBlockPositionUsingDirection(jt.rotateDirection))
	newCenterRight := geometry.Add(newCenter, NewCenterRightBlockPositionUsingDirection(jt.rotateDirection))
	newCenterLeftUp := geometry.Add(newCenter, geometry.Add(NewCenterLeftBlockPositionUsingDirection(jt.rotateDirection), NewCenterUpBlockPositionUsingDirection(jt.rotateDirection)))

	newCoordinateSlice := []geometry.Coordinates{newCenter, newCenterLeft, newCenterRight, newCenterLeftUp}

	RemoveAllBlocks(jt.blocks, jt)

	if CanTetrominoFit(newCoordinateSlice, jt) {
		jt.InsertBlock(newCenter.GetX(), newCenter.GetY(), jt.centerBlock)
		jt.InsertBlock(newCenterLeft.GetX(), newCenterLeft.GetY(), jt.centerLeftBlock)
		jt.InsertBlock(newCenterRight.GetX(), newCenterRight.GetY(), jt.centerRightBlock)
		jt.InsertBlock(newCenterLeftUp.GetX(), newCenterLeftUp.GetY(), jt.centerLeftUpBlock)
	} else {
		InsertAllBlocks(jt.blocks, jt)
		jt.rotateDirection = jt.rotateDirection.TurnClockWise90()
	}

}

type STetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerBlock   *Block
	centerLeft    *Block
	centerUp      *Block
	centerUpRight *Block

	rotateDirection geometry.Direction
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewSTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	st := STetrominoes{
		centerBlock:             NewBlock(x, y),
		centerLeft:              NewBlock(x-1, y),
		centerUp:                NewBlock(x, y+1),
		centerUpRight:           NewBlock(x+1, y+1),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	st.blocks = []*Block{st.centerBlock, st.centerLeft, st.centerUp, st.centerUpRight}
	SetColorOfBlocks(st.blocks, colors.Green)

	for _, block := range st.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &st, true
}

func (st *STetrominoes) Rotate() {

	x, y := st.centerBlock.GetX(), st.centerBlock.GetY()

	st.rotateDirection = st.rotateDirection.TurnClockWise90()

	newCenterBlock := geometry.NewCoordinate(x, y)
	newCenterLeft := geometry.Add(newCenterBlock, NewCenterLeftBlockPositionUsingDirection(st.rotateDirection))
	newCenterUp := geometry.Add(newCenterBlock, NewCenterUpBlockPositionUsingDirection(st.rotateDirection))
	newCenterUpRight := geometry.Add(newCenterBlock, geometry.Add(NewCenterUpBlockPositionUsingDirection(st.rotateDirection), NewCenterRightBlockPositionUsingDirection(st.rotateDirection)))

	newCoordinateSlice := []geometry.Coordinates{newCenterBlock, newCenterLeft, newCenterUp, newCenterUpRight}

	RemoveAllBlocks(st.blocks, st)

	canRotate := CanTetrominoFit(newCoordinateSlice, st.BlockInserterAndRemover)

	if canRotate {
		st.InsertBlock(newCenterBlock.GetX(), newCenterBlock.GetY(), st.centerBlock)
		st.InsertBlock(newCenterLeft.GetX(), newCenterLeft.GetY(), st.centerLeft)
		st.InsertBlock(newCenterUp.GetX(), newCenterUp.GetY(), st.centerUp)
		st.InsertBlock(newCenterUpRight.GetX(), newCenterUpRight.GetY(), st.centerUpRight)
	} else {
		InsertAllBlocks(st.blocks, st)
		st.rotateDirection = st.rotateDirection.TurnAntiClockWise90()
	}

}

type ZTetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerBlock  *Block
	centerRight  *Block
	centerUp     *Block
	centerUpLeft *Block

	rotateDirection geometry.Direction
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewZTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	zt := ZTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerUp:                NewBlock(x, y+1),
		centerUpLeft:            NewBlock(x-1, y+1),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	zt.blocks = []*Block{zt.centerBlock, zt.centerRight, zt.centerUp, zt.centerUpLeft}
	SetColorOfBlocks(zt.blocks, colors.Red)

	for _, block := range zt.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &zt, true
}

func (zt *ZTetrominoes) Rotate() {

	x, y := zt.centerBlock.GetX(), zt.centerBlock.GetY()

	zt.rotateDirection = zt.rotateDirection.TurnClockWise90()

	newCenter := geometry.NewCoordinate(x, y)
	newCenterRight := geometry.Add(newCenter, NewCenterRightBlockPositionUsingDirection(zt.rotateDirection))
	newCenterUp := geometry.Add(newCenter, NewCenterUpBlockPositionUsingDirection(zt.rotateDirection))
	newCenterUpLeft := geometry.Add(newCenter, geometry.Add(NewCenterUpBlockPositionUsingDirection(zt.rotateDirection), NewCenterLeftBlockPositionUsingDirection(zt.rotateDirection)))

	newCoordinateSlice := []geometry.Coordinates{newCenter, newCenterRight, newCenterUp, newCenterUpLeft}

	for _, block := range zt.blocks {
		zt.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := CanTetrominoFit(newCoordinateSlice, zt.BlockInserterAndRemover)

	if canRotate {
		zt.InsertBlock(newCenter.GetX(), newCenter.GetY(), zt.centerBlock)
		zt.InsertBlock(newCenterRight.GetX(), newCenterRight.GetY(), zt.centerRight)
		zt.InsertBlock(newCenterUp.GetX(), newCenterUp.GetY(), zt.centerUp)
		zt.InsertBlock(newCenterUpLeft.GetX(), newCenterUpLeft.GetY(), zt.centerUpLeft)
	} else {
		zt.rotateDirection = zt.rotateDirection.TurnAntiClockWise90()
	}

}

type TTetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerBlock *Block
	centerRight *Block
	centerUp    *Block
	centerLeft  *Block

	rotateDirection geometry.Direction
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewTTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	tt := TTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerUp:                NewBlock(x, y+1),
		centerLeft:              NewBlock(x-1, y),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	tt.blocks = []*Block{tt.centerBlock, tt.centerRight, tt.centerUp, tt.centerLeft}
	SetColorOfBlocks(tt.blocks, colors.Purple)

	for _, block := range tt.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &tt, true
}

func (tt *TTetrominoes) Rotate() {

	x, y := tt.centerBlock.GetX(), tt.centerBlock.GetY()

	tt.rotateDirection = tt.rotateDirection.TurnClockWise90()

	newCenter := geometry.NewCoordinate(x, y)
	newCenterRight := geometry.Add(newCenter, NewCenterRightBlockPositionUsingDirection(tt.rotateDirection))
	newCenterUp := geometry.Add(newCenter, NewCenterUpBlockPositionUsingDirection(tt.rotateDirection))
	newCenterLeft := geometry.Add(newCenter, NewCenterLeftBlockPositionUsingDirection(tt.rotateDirection))

	newCoordinateSlice := []geometry.Coordinates{newCenter, newCenterRight, newCenterUp, newCenterLeft}

	RemoveAllBlocks(tt.blocks, tt.BlockInserterAndRemover)

	canRotate := CanTetrominoFit(newCoordinateSlice, tt.BlockInserterAndRemover)

	if canRotate {
		tt.InsertBlock(newCenter.GetX(), newCenter.GetY(), tt.centerBlock)
		tt.InsertBlock(newCenterRight.GetX(), newCenterRight.GetY(), tt.centerRight)
		tt.InsertBlock(newCenterUp.GetX(), newCenterUp.GetY(), tt.centerUp)
		tt.InsertBlock(newCenterLeft.GetX(), newCenterLeft.GetY(), tt.centerLeft)
	} else {

		for _, block := range tt.blocks {
			tt.InsertBlock(block.GetX(), block.GetY(), block)
		}

		tt.rotateDirection = tt.rotateDirection.TurnAntiClockWise90()
	}

}

type ITetrominoes struct {
	BlockInserterAndRemover
	BlockHolder
	geometry.Coordinates

	centerLeftLeft   *Block
	centerLeft       *Block
	centerRight      *Block
	centerRightRight *Block

	rotateDirection geometry.Direction
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewITetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	it := ITetrominoes{
		centerLeftLeft:          NewBlock(x-1, y),
		centerLeft:              NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerRightRight:        NewBlock(x+2, y),
		rotateDirection:         geometry.Up,
		BlockInserterAndRemover: bir,
	}

	it.blocks = []*Block{it.centerLeftLeft, it.centerLeft, it.centerRight, it.centerRightRight}
	SetColorOfBlocks(it.blocks, colors.Cyan)

	for _, block := range it.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &it, true
}

func (it *ITetrominoes) Rotate() {

	x, y := it.centerLeft.GetX(), it.centerLeft.GetY()

	it.rotateDirection = it.rotateDirection.TurnClockWise90()

	var centerLeftLeft, centerLeft, centerRight, centerRightRight geometry.Coordinates

	switch it.rotateDirection {
	case geometry.Up:
		centerLeftLeft, centerLeft, centerRight, centerRightRight = geometry.NewCoordinate(x-1, y), geometry.NewCoordinate(x, y), geometry.NewCoordinate(x+1, y), geometry.NewCoordinate(x+2, y)
	case geometry.Left:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = geometry.NewCoordinate(x, y+2), geometry.NewCoordinate(x, y+1), geometry.NewCoordinate(x, y), geometry.NewCoordinate(x, y-1)
	case geometry.Down:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = geometry.NewCoordinate(x+2, y-1), geometry.NewCoordinate(x+1, y-1), geometry.NewCoordinate(x, y-1), geometry.NewCoordinate(x-1, y-1)
	case geometry.Right:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = geometry.NewCoordinate(x, y-1), geometry.NewCoordinate(x, y), geometry.NewCoordinate(x, y+1), geometry.NewCoordinate(x, y+2)
	}

	newCoordinateSlice := []geometry.Coordinates{centerLeftLeft, centerLeft, centerRight, centerRightRight}

	for _, block := range it.blocks {
		it.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := CanTetrominoFit(newCoordinateSlice, it.BlockInserterAndRemover)

	if canRotate {
		it.InsertBlock(centerLeftLeft.GetX(), centerLeftLeft.GetY(), it.centerLeftLeft)
		it.InsertBlock(centerLeft.GetX(), centerLeft.GetY(), it.centerLeft)
		it.InsertBlock(centerRight.GetX(), centerRight.GetY(), it.centerRight)
		it.InsertBlock(centerRightRight.GetX(), centerRightRight.GetY(), it.centerRightRight)
	} else {

		for _, block := range it.blocks {
			it.InsertBlock(block.GetX(), block.GetY(), block)
		}

		it.rotateDirection = it.rotateDirection.TurnAntiClockWise90()
	}

}
