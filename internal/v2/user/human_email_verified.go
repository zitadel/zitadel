package user

import (
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type HumanEmailVerifiedEvent eventstore.Event[eventstore.EmptyPayload]

const HumanEmailVerifiedType = humanPrefix + ".email.verified"

var _ eventstore.TypeChecker = (*HumanEmailVerifiedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanEmailVerifiedEvent) ActionType() string {
	return HumanEmailVerifiedType
}

func HumanEmailVerifiedEventFromStorage(event *eventstore.StorageEvent) (e *HumanEmailVerifiedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-X3esB", "Errors.Invalid.Event.Type")
	}

	return &HumanEmailVerifiedEvent{
		StorageEvent: event,
	}, nil
}
