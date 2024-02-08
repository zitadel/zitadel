package org

type State uint8

const (
	UndefinedState State = iota
	ActiveState
	InactiveState
	RemovedState
)
