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

func (e ExecutionType) String() string {
	switch e {
	case ExecutionTypeUnspecified, executionTypeStateCount:
		return ""
	case ExecutionTypeRequest:
		return "request"
	case ExecutionTypeResponse:
		return "response"
	case ExecutionTypeFunction:
		return "function"
	case ExecutionTypeEvent:
		return "event"
	}
	return ""
}

type ExecutionTargetType uint

func (s ExecutionTargetType) Valid() bool {
	return s < executionTargetTypeStateCount
}

const (
	ExecutionTargetTypeUnspecified ExecutionTargetType = iota
	ExecutionTargetTypeInclude
	ExecutionTargetTypeTarget

	executionTargetTypeStateCount
)
