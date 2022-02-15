package projection

import (
	"time"
)

type Config struct {
	RequeueEvery     time.Duration
	RetryFailedAfter time.Duration
	MaxFailureCount  uint
	BulkLimit        uint64
	Customizations   map[string]CustomConfig
	MaxIterators     int
}

type CustomConfig struct {
	RequeueEvery     *time.Duration
	RetryFailedAfter *time.Duration
	MaxFailureCount  *uint
	BulkLimit        *uint64
}
