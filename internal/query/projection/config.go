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
	MaxActiveInstances    uint32
	TransactionDuration   time.Duration
	ActiveInstancer       interface {
		ActiveInstances() []string
	}
	MaxParallelTriggers uint16
}

type CustomConfig struct {
	RequeueEvery        *time.Duration
	RetryFailedAfter    *time.Duration
	MaxFailureCount     *uint8
	ConcurrentInstances *uint
	BulkLimit           *uint16
	TransactionDuration *time.Duration
	SkipInstanceIDs     []string
}
