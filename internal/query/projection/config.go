package projection

import (
	"time"

	"github.com/caos/zitadel/internal/config/types"
)

type Config struct {
	RequeueEvery     time.Duration
	RetryFailedAfter time.Duration
	MaxFailureCount  uint
	BulkLimit        uint64
	CRDB             types.SQL
	Customizations   map[string]CustomConfig
	MaxIterators     int
}

type CustomConfig struct {
	RequeueEvery     *types.Duration
	RetryFailedAfter *types.Duration
	MaxFailureCount  *uint
	BulkLimit        *uint64
}
