package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	UserGrantAggregate models.AggregateType = "usergrant"

	UserGrantAdded       models.EventType = "user.grant.added"
	UserGrantChanged     models.EventType = "user.grant.changed"
	UserGrantRemoved     models.EventType = "user.grant.removed"
	UserGrantDeactivated models.EventType = "user.grant.deactivated"
	UserGrantReactivated models.EventType = "user.grant.reactivated"

	UserGrantCascadeRemoved models.EventType = "user.grant.cascade.removed"
	UserGrantCascadeChanged models.EventType = "user.grant.cascade.changed"
)
