package eventsourcing

import (
	"github.com/caos/logging"
)

func (c *PolicyCache) getAgePolicy(id string) (policy *PasswordAgePolicy) {
	policy = new(PasswordAgePolicy)
	if err := c.policyCache.Get(id, policy); err != nil {
		logging.Log("EVENT-NqMFM").WithError(err).Debug("error in getting cache")
	}
	return policy
}

func (c *PolicyCache) cacheAgePolicy(policy *PasswordAgePolicy) {
	err := c.policyCache.Set(policy.AggregateID, policy)
	if err != nil {
		logging.Log("EVENT-6vRvM").WithError(err).Debug("error in setting policy cache")
	}
}
