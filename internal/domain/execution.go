package domain

type ExecutionType uint

const (
	ExecutionTypeUndefined ExecutionType = iota // undefined
	ExecutionTypeWebhook
	ExecutionTypeRequestResponse

	executionTypeStateCount // invalid
)

type ExecutionState int32

const (
	ExecutionUnspecified ExecutionState = iota
	ExecutionActive
	ExecutionRemoved
	executionStateCount
)

func (s ExecutionState) Valid() bool {
	return s >= 0 && s < executionStateCount
}

func (s ExecutionState) Exists() bool {
	return s != ExecutionUnspecified && s != ExecutionRemoved
}
