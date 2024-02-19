package org

type State uint8

const (
	UndefinedState State = iota
	ActiveState
	InactiveState
	RemovedState
	maxState
)

func (s State) IsValid() bool {
	return s != UndefinedState ||
		s < maxState
}

func (s State) Is(state State) bool {
	return s == state
}

func (s State) IsValidState(state State) bool {
	return s.IsValid() && s.Is(state)
}
