package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	PolicyAggregate models.AggregateType = "policy"

	PolicyAdded   models.EventType = "policy.added"
	PolicyChanged models.EventType = "policy.changed"
)
