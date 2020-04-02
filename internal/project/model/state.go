package model

// code below could be generated

type state int32

func (s state) String() string {
	return states[s]
}

const (
	Active state = iota
	Inactive
)

func ProjectStateToInt(s ProjectState) int32 {
	if s == nil {
		return 0
	}
	return int32(s.(state))
}

func ProjectStateFromInt(index int32) ProjectState {
	return state(index)
}
