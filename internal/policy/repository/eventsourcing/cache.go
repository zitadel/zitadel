package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/cache"
	"github.com/caos/zitadel/internal/cache/config"
	"github.com/caos/zitadel/internal/eventstore/models"
)

type PolicyCache struct {
	policyCache cache.Cache
}

func StartCache(conf *config.CacheConfig) (*PolicyCache, error) {
	policyCache, err := conf.Config.NewCache()
	logging.Log("EVENT-vDneN").OnError(err).Panic("unable to create policy cache")

	return &PolicyCache{policyCache: policyCache}, nil
}

func (c *PolicyCache) getPolicy(id string) (policy *PasswordComplexityPolicy) {
	policy = &PasswordComplexityPolicy{ObjectRoot: models.ObjectRoot{AggregateID: id}}
	if err := c.policyCache.Get(id, policy); err != nil {
		logging.Log("EVENT-4eTZh").WithError(err).Debug("error in getting cache")
	}
	return policy
}

func (c *PolicyCache) cachePolicy(policy *PasswordComplexityPolicy) {
	err := c.policyCache.Set(policy.AggregateID, policy)
	if err != nil {
		logging.Log("EVENT-ThnBb").WithError(err).Debug("error in setting policy cache")
	}
}
