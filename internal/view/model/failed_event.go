package model

type FailedEvent struct {
	Database     string
	ViewName     string
	EventID      string
	FailureCount uint64
	ErrMsg       string
}
