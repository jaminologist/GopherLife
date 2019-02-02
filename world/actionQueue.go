package world

type QueueableActions interface {
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
	basicActionQueue.actionQueue <- action
}

func (basicActionQueue *BasicActionQueue) Process() {
	close(basicActionQueue.actionQueue)

	for action := range basicActionQueue.actionQueue {
		action()
	}

	basicActionQueue.actionQueue = make(chan func(), basicActionQueue.maxActions)
}
