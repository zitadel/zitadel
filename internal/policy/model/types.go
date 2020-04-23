package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	// complexity
	PasswordComplexityPolicyAggregate models.AggregateType = "policy.password.complexity"

	PasswordComplexityPolicyAdded   models.EventType = "policy.password.complexity.added"
	PasswordComplexityPolicyChanged models.EventType = "policy.password.complexity.changed"

	// age
	PasswordAgePolicyAggregate models.AggregateType = "policy.password.age"

	PasswordAgePolicyAdded   models.EventType = "policy.password.age.added"
	PasswordAgePolicyChanged models.EventType = "policy.password.age.changed"

	// lokout
	PasswordLokoutPolicyAggregate models.AggregateType = "policy.password.lokout"

	PasswordLokoutPolicyAdded   models.EventType = "policy.password.age.lokout"
	PasswordLokoutPolicyChanged models.EventType = "policy.password.age.lokout"
)
