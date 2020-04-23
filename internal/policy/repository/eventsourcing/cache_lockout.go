package eventsourcing

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore/models"
)

func (c *PolicyCache) getLockoutPolicy(ID string) (policy *PasswordLockoutPolicy) {
	policy = &PasswordLockoutPolicy{ObjectRoot: models.ObjectRoot{}}
	if err := c.policyCache.Get(ID, policy); err != nil {
		logging.Log("EVENT-4eTZh").WithError(err).Debug("error in getting cache")
	}
	return policy
}

func (c *PolicyCache) cacheLockoutPolicy(policy *PasswordLockoutPolicy) {
	err := c.policyCache.Set(policy.ID, policy)
	if err != nil {
		logging.Log("EVENT-ThnBb").WithError(err).Debug("error in setting policy cache")
	}
}
