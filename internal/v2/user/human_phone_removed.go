package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanPhoneRemovedEvent eventstore.Event[eventstore.EmptyPayload]

const HumanPhoneRemovedType = humanPrefix + ".phone.removed"

var _ eventstore.TypeChecker = (*HumanPhoneRemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanPhoneRemovedEvent) ActionType() string {
	return HumanPhoneRemovedType
}

func HumanPhoneRemovedEventFromStorage(event *eventstore.StorageEvent) (e *HumanPhoneRemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-vaD75", "Errors.Invalid.Event.Type")
	}

	return &HumanPhoneRemovedEvent{
		StorageEvent: event,
	}, nil
}
