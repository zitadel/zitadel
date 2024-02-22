package domain

type ExecutionType uint

func (s ExecutionType) Valid() bool {
	return s < executionTypeStateCount
}

const (
	ExecutionTypeUnspecified ExecutionType = iota
	ExecutionTypeRequest
	ExecutionTypeResponse
	ExecutionTypeFunction
	ExecutionTypeEvent

	executionTypeStateCount
)
