package model

type orgState int32

func (s orgState) String() string {
	return states[s]
}

const (
	Active orgState = iota
	Inactive
)

func ProjectStateToInt(s OrgState) int32 {
	if s == nil {
		return 0
	}
	return int32(s.(orgState))
}

func ProjectStateFromInt(index int32) OrgState {
	return orgState(index)
}
