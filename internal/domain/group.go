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
