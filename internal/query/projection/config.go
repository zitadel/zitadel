package projection

import "github.com/caos/zitadel/internal/config/types"

type Config struct {
	RequeueEvery     types.Duration
	RetryFailedAfter types.Duration
	MaxFailureCount  uint
	BulkLimit        uint64
	CRDB             types.SQL
	Customizations   map[string]CustomConfig
}

type CustomConfig struct {
	RequeueEvery     types.Duration
	RetryFailedAfter types.Duration
	MaxFailureCount  uint
	BulkLimit        uint64
}
