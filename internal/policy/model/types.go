package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	//default
	DefaultPolicy = "0"

	// complexity
	PasswordComplexityPolicyAggregate models.AggregateType = "policy.password.complexity"

	PasswordComplexityPolicyAdded   models.EventType = "policy.password.complexity.added"
	PasswordComplexityPolicyChanged models.EventType = "policy.password.complexity.changed"

	// age
	PasswordAgePolicyAggregate models.AggregateType = "policy.password.age"

	PasswordAgePolicyAdded   models.EventType = "policy.password.age.added"
	PasswordAgePolicyChanged models.EventType = "policy.password.age.changed"

	// lockout
	PasswordLockoutPolicyAggregate models.AggregateType = "policy.password.lockout"

	PasswordLockoutPolicyAdded   models.EventType = "policy.password.lockout.added"
	PasswordLockoutPolicyChanged models.EventType = "policy.password.lockout.changed"
)

type PolicyState int32

const (
	PolicyStateActive PolicyState = iota
	PolicyStateInactive
)
