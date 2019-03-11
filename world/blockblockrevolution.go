package world

import (
	"fmt"
	"gopherlife/calc"
	"time"
)

//10 x 20

type BlockBlockRevolutionMap struct {
	grid [][]*BlockBlockRevolutionTile
	Dimensions
	Containable
	CurrentTetromino Tetromino

	ActionQueuer

	FrameTimer calc.StopWatch
	FrameSpeed time.Duration
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
		FrameSpeed:   30,
	}

	bbrm.AddNewBlock()

	return bbrm

}

func (bbrm *BlockBlockRevolutionMap) Update() bool {

	bbrm.FrameTimer.Start()
	bbrm.Process()

	if !bbrm.MoveCurrentTetrominoDown() {
		//bbrm.CurrentTetromino = nil
		bbrm.CheckForAndClearLines()
		bbrm.AddNewBlock()
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
			//for i:= 0; i < bbrm.Width

			//linesToClear = append(linesToClear, y)
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
				fmt.Println("pre", tile.Block)
				blocks = append(blocks, tile.Block)
				tile.RemoveBlock()
			}
		}
	}

	for _, block := range blocks {
		fmt.Println(block)
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
	block, ok := NewLTetrominoes(bbrm.Width/2, bbrm.Height-2, bbrm)

	if ok {
		bbrm.CurrentTetromino = block
		return true
	}

	return false

}

type BlockInserterAndRemover interface {
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
	Shift(d Direction)
	Rotate()
	Blocks() []*Block
}

type BlockRotator interface {
	Rotate()
}

type Block struct {
	Position
}

type SquareTetrominoes struct {
	BlockInserterAndRemover
	blocks []*Block
	Position
}

func NewSquareTetrominoes(x int, y int, bir BlockInserterAndRemover) (*SquareTetrominoes, bool) {

	topLeftBlock := &Block{}
	topRightBlock := &Block{}
	bottomLeftBlock := &Block{}
	bottomRightBlock := &Block{}

	blocks := []*Block{topLeftBlock, topRightBlock, bottomLeftBlock, bottomRightBlock}

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

func (s *SquareTetrominoes) Shift(d Direction) {

}

func (s *SquareTetrominoes) Blocks() []*Block {
	return s.blocks
}

type LTetrominoes struct {
	BlockInserterAndRemover

	middleBlock        *Block
	leftMiddleBlock    *Block
	rightMiddleBlock   *Block
	rightMiddleUpBlock *Block

	rotatePosition int

	blocks []*Block
	Position
}

//NewLTetrominoes Creates a new L-Tetromino. Assume the starting position is an 3 blocks in a straight line and then one at the end
func NewLTetrominoes(x int, y int, bir BlockInserterAndRemover) (*LTetrominoes, bool) {

	middleBlock := &Block{}
	leftMiddleBlock := &Block{}
	rightMiddleBlock := &Block{}
	rightMiddleUpBlock := &Block{}

	blocks := []*Block{middleBlock, leftMiddleBlock, rightMiddleBlock, rightMiddleUpBlock}

	if !bir.InsertBlock(x, y, middleBlock) {
		return nil, false
	}

	if !bir.InsertBlock(x-1, y, leftMiddleBlock) {
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
		leftMiddleBlock:         leftMiddleBlock,
		rightMiddleBlock:        rightMiddleBlock,
		rightMiddleUpBlock:      rightMiddleUpBlock,
		BlockInserterAndRemover: bir,
	}

	return &lt, true
}

func (l *LTetrominoes) Rotate() {

	x, y := l.middleBlock.GetX(), l.middleBlock.GetY()
	var newMPos, newLMPos, newRMpos, newRMUpos calc.Coordinates

	newMPos = calc.NewCoordinate(x, y)

	if l.rotatePosition == 3 {
		l.rotatePosition = 0
	} else {
		l.rotatePosition++
	}

	switch l.rotatePosition {
	case 0:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x+1, y+1)
	case 1:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x, y+1), calc.NewCoordinate(x, y-1), calc.NewCoordinate(x+1, y-1)
	case 2:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x-1, y), calc.NewCoordinate(x+1, y), calc.NewCoordinate(x-1, y-1)
	case 3:
		newLMPos, newRMpos, newRMUpos = calc.NewCoordinate(x, y-1), calc.NewCoordinate(x, y+1), calc.NewCoordinate(x-1, y+1)
	}

	newCoordinateSlice := []calc.Coordinates{newMPos, newLMPos, newRMpos, newRMUpos}

	for i := 0; i < len(l.blocks); i++ {
		l.RemoveBlock(l.blocks[i].GetX(), l.blocks[i].GetY())
	}
	for _, block := range l.blocks {
		l.RemoveBlock(block.GetX(), block.GetY())
	}

	canRotate := true

	for _, coords := range newCoordinateSlice {
		if _, ok := l.ContainsBlock(coords.GetX(), coords.GetY()); ok {
			fmt.Println(coords.GetX(), coords.GetY())
			canRotate = false
		}
	}

	if canRotate {
		l.InsertBlock(newMPos.GetX(), newMPos.GetY(), l.middleBlock)
		l.InsertBlock(newLMPos.GetX(), newLMPos.GetY(), l.leftMiddleBlock)
		l.InsertBlock(newRMpos.GetX(), newRMpos.GetY(), l.rightMiddleBlock)
		l.InsertBlock(newRMUpos.GetX(), newRMUpos.GetY(), l.rightMiddleUpBlock)
	} else {
		if l.rotatePosition == 0 {
			l.rotatePosition = 3
		} else {
			l.rotatePosition--
		}
	}

}

func (l *LTetrominoes) Shift(d Direction) {

}

func (l *LTetrominoes) Blocks() []*Block {
	return l.blocks
}
