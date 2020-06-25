package model

type FailedEvent struct {
	Database       string
	ViewName       string
	FailedSequence uint64
	FailureCount   uint64
	ErrMsg         string
}
