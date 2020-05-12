package eventsourcing

import (
	"github.com/caos/zitadel/internal/cache/config"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/policy"
	"github.com/sony/sonyflake"
)

var idGenerator = sonyflake.NewSonyflake(sonyflake.Settings{})

type PolicyEventstore struct {
	es_int.Eventstore
	policyCache                     *PolicyCache
	passwordAgePolicyDefault        policy.PasswordAgePolicyDefault
	passwordComplexityPolicyDefault policy.PasswordComplexityPolicyDefault
	passwordLockoutPolicyDefault    policy.PasswordLockoutPolicyDefault
}

type PolicyConfig struct {
	es_int.Eventstore
	Cache *config.CacheConfig
}

func StartPolicy(conf PolicyConfig, systemDefaults sd.SystemDefaults) (*PolicyEventstore, error) {
	policyCache, err := StartCache(conf.Cache)
	if err != nil {
		return nil, err
	}
	return &PolicyEventstore{
		Eventstore:                      conf.Eventstore,
		policyCache:                     policyCache,
		passwordAgePolicyDefault:        systemDefaults.DefaultPolicies.Age,
		passwordComplexityPolicyDefault: systemDefaults.DefaultPolicies.Complexity,
		passwordLockoutPolicyDefault:    systemDefaults.DefaultPolicies.Lockout,
	}, nil
}
