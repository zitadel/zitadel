package projection

import "github.com/caos/zitadel/internal/config/types"

type Config struct {
	RequeueEvery    types.Duration
	MaxFailureCount uint
	BulkLimit       uint64
	CRDB            types.SQL
}
