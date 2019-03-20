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

func (finiteActionQueue *FiniteActionQueue) Add(action func()) {
	select {
	case finiteActionQueue.actionQueue <- action: // Put 2 in the channel unless it is full
	default:
	}
}

func (finiteActionQueue *FiniteActionQueue) Process() {
	actionChannel := finiteActionQueue.actionQueue
	finiteActionQueue.actionQueue = make(chan func(), finiteActionQueue.maxActions)
	close(actionChannel)
	for action := range actionChannel {
		action()
	}
}
