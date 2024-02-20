package domain

type ExecutionType uint

func (s ExecutionType) Valid() bool {
	return s >= 0 && s < executionTypeStateCount
}

const (
	ExecutionTypeUnspecified ExecutionType = iota
	ExecutionTypeRequest
	ExecutionTypeResponse
	ExecutionTypeFunction
	ExecutionTypeEvent

	executionTypeStateCount
)
