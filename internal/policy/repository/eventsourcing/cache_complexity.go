package eventsourcing

import (
	"github.com/caos/logging"
)

func (c *PolicyCache) getComplexityPolicy(id string) (policy *PasswordComplexityPolicy) {
	policy = new(PasswordComplexityPolicy)
	if err := c.policyCache.Get(id, policy); err != nil {
		logging.Log("EVENT-Wgrph").WithError(err).Debug("error in getting cache")
	}
	return policy
}

func (c *PolicyCache) cacheComplexityPolicy(policy *PasswordComplexityPolicy) {
	err := c.policyCache.Set(policy.AggregateID, policy)
	if err != nil {
		logging.Log("EVENT-ylGny").WithError(err).Debug("error in setting policy cache")
	}
}
