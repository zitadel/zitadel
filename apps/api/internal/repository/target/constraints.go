package target

import (
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueTarget    = "target"
	DuplicateTarget = "Errors.Target.AlreadyExists"
)

func NewAddUniqueConstraint(name string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueTarget,
		name,
		DuplicateTarget,
	)
}

func NewRemoveUniqueConstraint(name string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueTarget,
		name,
	)
}
