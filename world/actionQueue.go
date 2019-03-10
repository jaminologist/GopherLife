package world

import "fmt"

type ActionQueuer interface {
	Add(action func())
	Process()
}

type BasicActionQueue struct {
	actionQueue chan func()
	maxActions  int
}

func NewBasicActionQueue(maxActions int) BasicActionQueue {
	return BasicActionQueue{
		actionQueue: make(chan func(), maxActions),
		maxActions:  maxActions,
	}
}

func (basicActionQueue *BasicActionQueue) Add(action func()) {
	select {
	case basicActionQueue.actionQueue <- action: // Put 2 in the channel unless it is full
	default:
		fmt.Println("Channel full. Discarding value")
	}
}

func (basicActionQueue *BasicActionQueue) Process() {
	actionChannel := basicActionQueue.actionQueue
	basicActionQueue.actionQueue = make(chan func(), basicActionQueue.maxActions)
	close(actionChannel)
	for action := range actionChannel {
		action()
	}
}
