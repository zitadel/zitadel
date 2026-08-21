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

type GroupGrantState int32

const (
	GroupGrantStateUnspecified GroupGrantState = iota
	GroupGrantStateActive
	GroupGrantStateRemoved
)

func (g GroupGrantState) Exists() bool {
	return g != GroupGrantStateRemoved && g != GroupGrantStateUnspecified
}
