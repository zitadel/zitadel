package model

import "time"

type FailedEvent struct {
	Database       string
	ViewName       string
	FailedSequence uint64
	FailureCount   uint64
	ErrMsg         string
	InstanceID     string
	LastFailed     time.Time
}
