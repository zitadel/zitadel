package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type LockedEvent eventstore.Event[eventstore.EmptyPayload]

const LockedType = AggregateType + ".locked"

var _ eventstore.TypeChecker = (*LockedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *LockedEvent) ActionType() string {
	return LockedType
}

func LockedEventFromStorage(event *eventstore.StorageEvent) (e *LockedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-48jjE", "Errors.Invalid.Event.Type")
	}

	return &LockedEvent{
		StorageEvent: event,
	}, nil
}
