package eventsourcing

import (
	"github.com/caos/zitadel/internal/cache/config"
	es_int "github.com/caos/zitadel/internal/eventstore"
)

type IAMConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}
