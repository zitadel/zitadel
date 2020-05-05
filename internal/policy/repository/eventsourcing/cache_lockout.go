package eventsourcing

import (
	"github.com/caos/logging"
)

func (c *PolicyCache) getLockoutPolicy(id string) (policy *PasswordLockoutPolicy) {
	policy = new(PasswordLockoutPolicy)
	if err := c.policyCache.Get(id, policy); err != nil {
		logging.Log("EVENT-Zoljf").WithError(err).Debug("error in getting cache")
	}
	return policy
}

func (c *PolicyCache) cacheLockoutPolicy(policy *PasswordLockoutPolicy) {
	err := c.policyCache.Set(policy.AggregateID, policy)
	if err != nil {
		logging.Log("EVENT-6klAf").WithError(err).Debug("error in setting policy cache")
	}
}
