package world

import (
	"gopherlife/calc"
	"gopherlife/colors"
	"image/color"
	"math/rand"
	"time"
)

//10 x 20

type BlockBlockRevolutionMap struct {
	grid [][]*BlockBlockRevolutionTile
	Dimensions
	Containable
	CurrentTetromino Tetromino

	newBlockFunctions []func(int, int, BlockInserterAndRemover) (Tetromino, bool)

	nextNewBlockFunctions []func(int, int, BlockInserterAndRemover) (Tetromino, bool)

	ActionQueuer

	FrameTimer calc.StopWatch
	FrameSpeed time.Duration

	DownToNextLineCount int
}

const FrameSpeedMultiplier = time.Duration(7)

func NewBlockBlockRevolutionMap(d Dimensions, speed int) BlockBlockRevolutionMap {

	r := NewRectangle(0, 0, d.Width, d.Height)

	grid := make([][]*BlockBlockRevolutionTile, d.Width)

	for i := 0; i < d.Width; i++ {
		grid[i] = make([]*BlockBlockRevolutionTile, d.Height)

		for j := 0; j < d.Height; j++ {
			tile := BlockBlockRevolutionTile{
				Position: Position{
					X: i,
					Y: j,
				},
			}
			grid[i][j] = &tile
		}
	}

	qa := NewBasicActionQueue(1)

	bbrm := BlockBlockRevolutionMap{
		Containable:  &r,
		Dimensions:   d,
		ActionQueuer: &qa,
		grid:         grid,
		FrameSpeed:   5,
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

	bbrm.FrameTimer.Start()
	bbrm.Process()

	if bbrm.DownToNextLineCount < 5 {
		bbrm.DownToNextLineCount++
	} else {
		bbrm.DownToNextLineCount = 0
		if !bbrm.MoveCurrentTetrominoDown() {
			//bbrm.CurrentTetromino = nil
			bbrm.CheckForAndClearLines()
			bbrm.AddNewBlock()
		}
	}

	for bbrm.FrameTimer.GetCurrentElaspedTime() < time.Millisecond*FrameSpeedMultiplier*bbrm.FrameSpeed {
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
			y--
		}

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

type BlockInserterAndRemover interface {
	Containable
	InsertBlock(x int, y int, b *Block) bool
	ContainsBlock(x int, y int) (*Block, bool)
	RemoveBlock(x int, y int)
}

type BlockBlockRevolutionTile struct {
	Position
	Block *Block
}

func (bbrt *BlockBlockRevolutionTile) InsertBlock(b *Block) bool {

	if bbrt.Block == nil {
		b.SetPosition(bbrt.GetX(), bbrt.GetY())
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

type Block struct {
	Position
	Color color.RGBA
}

type SquareTetrominoes struct {
	BlockInserterAndRemover
	blocks []*Block
	Position
}

func SetColorOfBlocks(blocks []*Block, c color.RGBA) {
	for _, block := range blocks {
		block.Color = c
	}
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

	middleBlock        *Block
	centerLeftBlock    *Block
	rightMiddleBlock   *Block
	rightMiddleUpBlock *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewLTetrominoes Creates a new L-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewLTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	middleBlock := &Block{}
	centerLeftBlock := &Block{}
	rightMiddleBlock := &Block{}
	rightMiddleUpBlock := &Block{}

	blocks := []*Block{middleBlock, centerLeftBlock, rightMiddleBlock, rightMiddleUpBlock}
	SetColorOfBlocks(blocks, colors.Orange)

	if !bir.InsertBlock(x, y, middleBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x-1, y, centerLeftBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x+1, y, rightMiddleBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x+1, y+1, rightMiddleUpBlock) {
		return nil, false
	}

	lt := LTetrominoes{
		blocks:                  blocks,
		middleBlock:             middleBlock,
		centerLeftBlock:         centerLeftBlock,
		rightMiddleBlock:        rightMiddleBlock,
		rightMiddleUpBlock:      rightMiddleUpBlock,
		BlockInserterAndRemover: bir,
		rotateDirection:         Up,
	}

	return &lt, true
}

func (l *LTetrominoes) Rotate() {

	x, y := l.middleBlock.GetX(), l.middleBlock.GetY()
	var newMPos, newLMPos, newRMpos, newRMUpos calc.Coordinates

	newMPos = calc.NewCoordinate(x, y)

	l.rotateDirection = l.rotateDirection.TurnClockWise90()

	switch l.rotateDirection {
	case Up:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x+1, y+1)
	case Left:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x, y+1), calc.NewCoordinate(x, y-1), calc.NewCoordinate(x+1, y-1)
	case Down:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x-1, y-1)
	case Right:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x, y-1), calc.NewCoordinate(x, y+1), calc.NewCoordinate(x-1, y+1)
	}

	newCoordinateSlice := []calc.Coordinates{newMPos, newLMPos, newRMpos, newRMUpos}

	for _, block := range l.blocks {
		l.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := true

	for _, coords := range newCoordinateSlice {
		if _, ok := l.ContainsBlock(coords.GetX(), coords.GetY()); ok {
			canRotate = false
		}
	}

	if canRotate {
		l.InsertBlock(newMPos.GetX(), newMPos.GetY(), l.middleBlock)
		l.InsertBlock(newLMPos.GetX(), newLMPos.GetY(), l.centerLeftBlock)
		l.InsertBlock(newRMpos.GetX(), newRMpos.GetY(), l.rightMiddleBlock)
		l.InsertBlock(newRMUpos.GetX(), newRMUpos.GetY(), l.rightMiddleUpBlock)
	} else {
		l.rotateDirection = l.rotateDirection.TurnAntiClockWise90()
	}

}

func (l *LTetrominoes) Blocks() []*Block {
	return l.blocks
}

type JTetrominoes struct {
	BlockInserterAndRemover

	middleBlock       *Block
	centerLeftBlock   *Block
	rightMiddleBlock  *Block
	centerLeftUpBlock *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewJTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	middleBlock := &Block{}
	centerLeftBlock := &Block{}
	rightMiddleBlock := &Block{}
	centerLeftUpBlock := &Block{}

	blocks := []*Block{middleBlock, centerLeftBlock, rightMiddleBlock, centerLeftUpBlock}
	SetColorOfBlocks(blocks, colors.MingBlue)

	if !bir.InsertBlock(x, y, middleBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x-1, y, centerLeftBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x+1, y, rightMiddleBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x-1, y+1, centerLeftUpBlock) {
		return nil, false
	}

	jt := JTetrominoes{
		blocks:                  blocks,
		middleBlock:             middleBlock,
		centerLeftBlock:         centerLeftBlock,
		rightMiddleBlock:        rightMiddleBlock,
		centerLeftUpBlock:       centerLeftUpBlock,
		BlockInserterAndRemover: bir,
		rotateDirection:         Up,
	}

	return &jt, true
}

func (l *JTetrominoes) Rotate() {

	x, y := l.middleBlock.GetX(), l.middleBlock.GetY()
	var newMPos, newLMPos, newRMpos, newcenterLeftUpPos calc.Coordinates

	newMPos = calc.NewCoordinate(x, y)

	l.rotateDirection = l.rotateDirection.TurnClockWise90()

	switch l.rotateDirection {
	case Up:
		newLMPos, newRMpos, newcenterLeftUpPos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x-1, y+1)
	case Left:
		newLMPos, newRMpos, newcenterLeftUpPos = calc.NewCoordinate(x, y+1), calc.NewCoordinate(x, y-1), calc.NewCoordinate(x+1, y+1)
	case Down:
		newLMPos, newRMpos, newcenterLeftUpPos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x+1, y-1)
	case Right:
		newLMPos, newRMpos, newcenterLeftUpPos = calc.NewCoordinate(x, y-1), calc.NewCoordinate(x, y+1), calc.NewCoordinate(x-1, y+1)
	}

	newCoordinateSlice := []calc.Coordinates{newMPos, newLMPos, newRMpos, newcenterLeftUpPos}

	for _, block := range l.blocks {
		l.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := true

	for _, coords := range newCoordinateSlice {
		if _, ok := l.ContainsBlock(coords.GetX(), coords.GetY()); ok {
			canRotate = false
		}
	}

	if canRotate {
		l.InsertBlock(newMPos.GetX(), newMPos.GetY(), l.middleBlock)
		l.InsertBlock(newLMPos.GetX(), newLMPos.GetY(), l.centerLeftBlock)
		l.InsertBlock(newRMpos.GetX(), newRMpos.GetY(), l.rightMiddleBlock)
		l.InsertBlock(newcenterLeftUpPos.GetX(), newcenterLeftUpPos.GetY(), l.centerLeftUpBlock)
	} else {
		l.rotateDirection = l.rotateDirection.TurnAntiClockWise90()
	}

}

func (l *JTetrominoes) Blocks() []*Block {
	return l.blocks
}

func NewMiddleLeftBlockPositionUsingDirection(d Direction) calc.Coordinates {

	switch d {
	case Up:
		return calc.NewCoordinate(1, 0)
	case Left:
		return calc.NewCoordinate(0, 1)
	case Down:
		return calc.NewCoordinate(-1, 0)
	case Right:
		return calc.NewCoordinate(0, -1)
	default:
		panic("Invalid Direction Used")
	}

}

func NewMiddleRightBlockPositionUsingDirection(d Direction) calc.Coordinates {
	d = d.TurnClockWise90().TurnClockWise90()
	return NewMiddleLeftBlockPositionUsingDirection(d)
}

func NewMiddleUpBlockPositionUsingDirection(d Direction) calc.Coordinates {
	d = d.TurnClockWise90()
	return NewMiddleLeftBlockPositionUsingDirection(d)
}

func NewMiddleDownBlockPositionUsingDirection(d Direction) calc.Coordinates {
	d = d.TurnAntiClockWise90()
	return NewMiddleLeftBlockPositionUsingDirection(d)
}

func NewBlock(x int, y int) *Block {
	b := Block{
		Position: Position{x, y},
	}
	return &b
}

type STetrominoes struct {
	BlockInserterAndRemover

	middleBlock   *Block
	middleLeft    *Block
	middleUp      *Block
	middleUpRight *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewSTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	st := STetrominoes{
		middleBlock:             NewBlock(x, y),
		middleLeft:              NewBlock(x-1, y),
		middleUp:                NewBlock(x, y+1),
		middleUpRight:           NewBlock(x+1, y+1),
		rotateDirection:         Up,
		BlockInserterAndRemover: bir,
	}

	st.blocks = []*Block{st.middleBlock, st.middleLeft, st.middleUp, st.middleUpRight}
	SetColorOfBlocks(st.blocks, colors.Green)

	for _, block := range st.blocks {
		if !bir.InsertBlock(block.GetX(), block.GetY(), block) {
			return nil, false
		}
	}

	return &st, true
}

func CanTetrominoFit(cs []calc.Coordinates, bir BlockInserterAndRemover) bool {
	for _, coords := range cs {
		if _, ok := bir.ContainsBlock(coords.GetX(), coords.GetY()); ok || !bir.Contains(coords.GetX(), coords.GetY()) {
			return false
		}
	}
	return true
}

func (st *STetrominoes) Rotate() {

	x, y := st.middleBlock.GetX(), st.middleBlock.GetY()

	st.rotateDirection = st.rotateDirection.TurnClockWise90()

	newMPos := calc.NewCoordinate(x, y)
	newMiddleLeft := calc.Add(newMPos, NewMiddleLeftBlockPositionUsingDirection(st.rotateDirection))
	newMiddleUp := calc.Add(newMPos, NewMiddleUpBlockPositionUsingDirection(st.rotateDirection))
	newMiddleUpRight := calc.Add(newMPos, calc.Add(NewMiddleUpBlockPositionUsingDirection(st.rotateDirection), NewMiddleRightBlockPositionUsingDirection(st.rotateDirection)))

	newCoordinateSlice := []calc.Coordinates{newMPos, newMiddleLeft, newMiddleUp, newMiddleUpRight}

	for _, block := range st.blocks {
		st.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := CanTetrominoFit(newCoordinateSlice, st.BlockInserterAndRemover)

	if canRotate {
		st.InsertBlock(newMPos.GetX(), newMPos.GetY(), st.middleBlock)
		st.InsertBlock(newMiddleLeft.GetX(), newMiddleLeft.GetY(), st.middleLeft)
		st.InsertBlock(newMiddleUp.GetX(), newMiddleUp.GetY(), st.middleUp)
		st.InsertBlock(newMiddleUpRight.GetX(), newMiddleUpRight.GetY(), st.middleUpRight)
	} else {
		st.rotateDirection = st.rotateDirection.TurnAntiClockWise90()
	}

}

func (st *STetrominoes) Blocks() []*Block {
	return st.blocks
}

type ZTetrominoes struct {
	BlockInserterAndRemover

	centerBlock  *Block
	centerRight  *Block
	centerUp     *Block
	centerUpLeft *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewZTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	zt := ZTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerUp:                NewBlock(x, y+1),
		centerUpLeft:            NewBlock(x-1, y+1),
		rotateDirection:         Up,
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

	newCenter := calc.NewCoordinate(x, y)
	newCenterRight := calc.Add(newCenter, NewMiddleRightBlockPositionUsingDirection(zt.rotateDirection))
	newCenterUp := calc.Add(newCenter, NewMiddleUpBlockPositionUsingDirection(zt.rotateDirection))
	newCenterUpLeft := calc.Add(newCenter, calc.Add(NewMiddleUpBlockPositionUsingDirection(zt.rotateDirection), NewMiddleLeftBlockPositionUsingDirection(zt.rotateDirection)))

	newCoordinateSlice := []calc.Coordinates{newCenter, newCenterRight, newCenterUp, newCenterUpLeft}

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

func (zt *ZTetrominoes) Blocks() []*Block {
	return zt.blocks
}

type TTetrominoes struct {
	BlockInserterAndRemover

	centerBlock *Block
	centerRight *Block
	centerUp    *Block
	centerLeft  *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewTTetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	tt := TTetrominoes{
		centerBlock:             NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerUp:                NewBlock(x, y+1),
		centerLeft:              NewBlock(x-1, y),
		rotateDirection:         Up,
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

	newCenter := calc.NewCoordinate(x, y)
	newCenterRight := calc.Add(newCenter, NewMiddleRightBlockPositionUsingDirection(tt.rotateDirection))
	newCenterUp := calc.Add(newCenter, NewMiddleUpBlockPositionUsingDirection(tt.rotateDirection))
	newCenterLeft := calc.Add(newCenter, NewMiddleLeftBlockPositionUsingDirection(tt.rotateDirection))

	newCoordinateSlice := []calc.Coordinates{newCenter, newCenterRight, newCenterUp, newCenterLeft}

	for _, block := range tt.blocks {
		tt.RemoveBlock(block.GetX(), block.GetY())
	}

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

func (tt *TTetrominoes) Shift(d Direction) {

}

func (tt *TTetrominoes) Blocks() []*Block {
	return tt.blocks
}

type ITetrominoes struct {
	BlockInserterAndRemover

	centerLeftLeft   *Block
	centerLeft       *Block
	centerRight      *Block
	centerRightRight *Block

	rotateDirection Direction

	blocks []*Block
	Position
}

//NewJTetrominoes Creates a new J-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewITetrominoes(x int, y int, bir BlockInserterAndRemover) (Tetromino, bool) {

	it := ITetrominoes{
		centerLeftLeft:          NewBlock(x-1, y),
		centerLeft:              NewBlock(x, y),
		centerRight:             NewBlock(x+1, y),
		centerRightRight:        NewBlock(x+2, y),
		rotateDirection:         Up,
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

	var centerLeftLeft, centerLeft, centerRight, centerRightRight calc.Coordinates

	switch it.rotateDirection {
	case Up:
		centerLeftLeft, centerLeft, centerRight, centerRightRight = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x+2, y)
	case Left:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = calc.NewCoordinate(x, y+2), calc.NewCoordinate(x, y+1), calc.NewCoordinate(x, y), calc.NewCoordinate(x, y-1)
	case Down:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = calc.NewCoordinate(x+2, y-1), calc.NewCoordinate(x+1, y-1), calc.NewCoordinate(x, y-1), calc.NewCoordinate(x-1, y-1)
	case Right:
		x, y = it.centerRight.GetX(), it.centerRight.GetY()
		centerLeftLeft, centerLeft, centerRight, centerRightRight = calc.NewCoordinate(x, y-1), calc.NewCoordinate(x, y), calc.NewCoordinate(x, y+1), calc.NewCoordinate(x, y+2)
	}

	newCoordinateSlice := []calc.Coordinates{centerLeftLeft, centerLeft, centerRight, centerRightRight}

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

func (it *ITetrominoes) Shift(d Direction) {

}

func (it *ITetrominoes) Blocks() []*Block {
	return it.blocks
}
