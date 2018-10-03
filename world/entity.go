package world

type entity interface {
	PerformAction(world *World)
}
