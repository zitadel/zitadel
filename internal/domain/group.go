package domain

type GroupState int32

const (
	GroupStateUnspecified GroupState = iota
	GroupStateActive
	GroupStateRemoved
)

func (g GroupState) Exists() bool {
	return g != GroupStateRemoved && g != GroupStateUnspecified
}

type GroupUserState int32

const (
	GroupUserStateUnspecified GroupUserState = iota
	GroupUserStateActive
	GroupUserStateRemoved
)

func (g GroupUserState) Exists() bool {
	return g != GroupUserStateRemoved && g != GroupUserStateUnspecified
}
