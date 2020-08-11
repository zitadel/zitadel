package model

import (
	"github.com/caos/zitadel/internal/eventstore/models"
)

const (
	ServiceAccountAggregate models.AggregateType = "serviceaccount"

	ServiceAccountAdded   models.EventType = "serviceaccount.added"
	ServiceAccountChanged models.EventType = "serviceaccount.changed"

	ServiceAccountLocked      models.EventType = "serviceaccount.locked"
	ServiceAccountUnlocked    models.EventType = "serviceaccount.unlocked"
	ServiceAccountDeactivated models.EventType = "serviceaccount.deactivated"
	ServiceAccountReactivated models.EventType = "serviceaccount.reactivated"
	ServiceAccountRemoved     models.EventType = "serviceaccount.removed"

	KeyAdded   models.EventType = "serviceaccount.key.added"
	KeyRemoved models.EventType = "serviceaccount.key.removed"
)
