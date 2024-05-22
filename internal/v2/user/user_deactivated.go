package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type DeactivatedEvent eventstore.Event[eventstore.EmptyPayload]

const DeactivatedType = AggregateType + ".deactivated"

var _ eventstore.TypeChecker = (*DeactivatedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *DeactivatedEvent) ActionType() string {
	return DeactivatedType
}

func DeactivatedEventFromStorage(event *eventstore.StorageEvent) (e *DeactivatedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-SBLu2", "Errors.Invalid.Event.Type")
	}

	return &DeactivatedEvent{
		StorageEvent: event,
	}, nil
}
