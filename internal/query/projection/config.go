package projection

import "github.com/caos/zitadel/internal/config/types"

type Config struct {
	RequeueEvery     types.Duration
	RetryFailedAfter types.Duration
	MaxFailureCount  uint
	BulkLimit        uint64
	CRDB             types.SQL
	Customizations   map[string]CustomConfig
	MaxIterators     int
}

type ConfigV2 struct {
	Config
	CRDB types.SQL2
}

type CustomConfig struct {
	RequeueEvery     *types.Duration
	RetryFailedAfter *types.Duration
	MaxFailureCount  *uint
	BulkLimit        *uint64
}
