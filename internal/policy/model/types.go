package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	PasswordComplexityPolicyAggregate models.AggregateType = "policy.password.complexity"

	PasswordComplexityPolicyAdded   models.EventType = "policy.password.complexity.added"
	PasswordComplexityPolicyChanged models.EventType = "policy.password.complexity.changed"
)
