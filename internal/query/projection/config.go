package projection

import (
	"time"
)

type Config struct {
	RequeueEvery            time.Duration
	RetryFailedAfter        time.Duration
	MaxFailureCount         uint
	ConcurrentInstances     uint
	BulkLimit               uint64
	Customizations          map[string]CustomConfig
	HandleInactiveInstances bool
}

type CustomConfig struct {
	RequeueEvery            *time.Duration
	RetryFailedAfter        *time.Duration
	MaxFailureCount         *uint
	ConcurrentInstances     *uint
	BulkLimit               *uint64
	HandleInactiveInstances *bool
}
