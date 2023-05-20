package eventstore

import "github.com/zitadel/zitadel/internal/eventstore/v3"

type EventUniqueConstraint = eventstore.UniqueConstraint

type UniqueConstraintAction = eventstore.UniqueConstraintAction

const (
	UniqueConstraintAdd            = eventstore.UniqueConstraintAdd
	UniqueConstraintRemove         = eventstore.UniqueConstraintRemove
	UniqueConstraintInstanceRemove = eventstore.UniqueConstraintInstanceRemove
)

func NewAddEventUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveEventUniqueConstraint(
	uniqueType,
	uniqueField string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		Action:      UniqueConstraintRemove,
	}
}

func NewRemoveInstanceUniqueConstraints() *EventUniqueConstraint {
	return &EventUniqueConstraint{
		Action: UniqueConstraintInstanceRemove,
	}
}

func NewAddGlobalEventUniqueConstraint(
	uniqueType,
	uniqueField,
	errMessage string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:   uniqueType,
		UniqueField:  uniqueField,
		ErrorMessage: errMessage,
		IsGlobal:     true,
		Action:       UniqueConstraintAdd,
	}
}

func NewRemoveGlobalEventUniqueConstraint(
	uniqueType,
	uniqueField string) *EventUniqueConstraint {
	return &EventUniqueConstraint{
		UniqueType:  uniqueType,
		UniqueField: uniqueField,
		IsGlobal:    true,
		Action:      UniqueConstraintRemove,
	}
}
