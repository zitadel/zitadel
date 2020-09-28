package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
)

type PolicyCache struct {
	policyCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*PolicyCache, error) {
	policyCache, err := conf.Config.NewCache()
	logging.Log("EVENT-L7ZcH").OnError(err).Panic("unable to create policy cache")

	return &PolicyCache{policyCache: policyCache}, nil
}
