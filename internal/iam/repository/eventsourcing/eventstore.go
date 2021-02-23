package eventsourcing

import (
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/v1"
)

type IAMConfig struct {
	v1.Eventstore
	Cache *config.CacheConfig
}
