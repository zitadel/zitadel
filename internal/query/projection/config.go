package projection

import (
	"time"
)

type Config struct {
	RequeueEvery          time.Duration
	RetryFailedAfter      time.Duration
	MaxFailureCount       uint
	ConcurrentInstances   uint
	BulkLimit             uint64
	Customizations        map[string]CustomConfig
	HandleActiveInstances time.Duration
}

type CustomConfig struct {
	RequeueEvery          *time.Duration
	RetryFailedAfter      *time.Duration
	MaxFailureCount       *uint
	ConcurrentInstances   *uint
	BulkLimit             *uint64
	HandleActiveInstances *time.Duration
}
