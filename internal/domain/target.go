package domain

type TargetState int32

const (
	TargetUnspecified TargetState = iota
	TargetActive
	TargetRemoved
	targetStateCount
)

func (s TargetState) Valid() bool {
	return s >= 0 && s < targetStateCount
}

func (s TargetState) Exists() bool {
	return s != TargetUnspecified && s != TargetRemoved
}
