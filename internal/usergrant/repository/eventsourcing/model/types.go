package model

import "github.com/caos/zitadel/internal/eventstore/models"

const (
	UserGrantAggregate       models.AggregateType = "usergrant"
	UserGrantUniqueAggregate models.AggregateType = "usergrant.unique"

	UserGrantAdded       models.EventType = "user.grant.added"
	UserGrantChanged     models.EventType = "user.grant.changed"
	UserGrantRemoved     models.EventType = "user.grant.removed"
	UserGrantDeactivated models.EventType = "user.grant.deactivated"
	UserGrantReactivated models.EventType = "user.grant.reactivated"
	UserGrantReserved    models.EventType = "user.grant.reserved"
	UserGrantReleased    models.EventType = "user.grant.released"
)
