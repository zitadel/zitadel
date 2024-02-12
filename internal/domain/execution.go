package domain

type ExecutionType uint

const (
	ExecutionTypeUndefined ExecutionType = iota // undefined
	ExecutionTypeWebhook
	ExecutionTypeRequestResponse

	executionTypeStateCount // invalid
)
