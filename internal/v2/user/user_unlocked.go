package user

import (
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type UnlockedEvent eventstore.Event[eventstore.EmptyPayload]

const UnlockedType = AggregateType + ".unlocked"

var _ eventstore.TypeChecker = (*UnlockedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *UnlockedEvent) ActionType() string {
	return UnlockedType
}

func UnlockedEventFromStorage(event *eventstore.StorageEvent) (e *UnlockedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-HB0wi", "Errors.Invalid.Event.Type")
	}

	return &UnlockedEvent{
		StorageEvent: event,
	}, nil
}
