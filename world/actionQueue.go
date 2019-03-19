package world

type ActionQueuer interface {
	Add(action func())
	Process()
}

type FiniteActionQueue struct {
	actionQueue chan func()
	maxActions  int
}

func NewFiniteActionQueue(maxActions int) FiniteActionQueue {
	return FiniteActionQueue{
		actionQueue: make(chan func(), maxActions),
		maxActions:  maxActions,
	}
}

func (basicActionQueue *FiniteActionQueue) Add(action func()) {
	select {
	case basicActionQueue.actionQueue <- action: // Put 2 in the channel unless it is full
	default:
	}
}

func (basicActionQueue *FiniteActionQueue) Process() {
	actionChannel := basicActionQueue.actionQueue
	basicActionQueue.actionQueue = make(chan func(), basicActionQueue.maxActions)
	close(actionChannel)
	for action := range actionChannel {
		action()
	}
}
