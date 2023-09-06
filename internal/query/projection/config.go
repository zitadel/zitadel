package projection

import (
	"time"
)

type Config struct {
	RequeueEvery          time.Duration
	RetryFailedAfter      time.Duration
	MaxFailureCount       uint8
	ConcurrentInstances   uint
	BulkLimit             uint64
	Customizations        map[string]CustomConfig
	HandleActiveInstances time.Duration
	TransactionDuration   time.Duration
}

type CustomConfig struct {
	RequeueEvery          *time.Duration
	RetryFailedAfter      *time.Duration
	MaxFailureCount       *uint8
	ConcurrentInstances   *uint
	BulkLimit             *uint16
	HandleActiveInstances *time.Duration
	TransactionDuration   *time.Duration
}
