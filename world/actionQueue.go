package world

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
	basicActionQueue.actionQueue <- action
}

func (basicActionQueue *BasicActionQueue) Process() {
	actionChannel := basicActionQueue.actionQueue
	basicActionQueue.actionQueue = make(chan func(), basicActionQueue.maxActions)
	close(actionChannel)
	for action := range actionChannel {
		action()
	}
}
