package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanPhoneVerifiedEvent eventstore.Event[eventstore.EmptyPayload]

const HumanPhoneVerifiedType = humanPrefix + ".phone.removed"

var _ eventstore.TypeChecker = (*HumanPhoneVerifiedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanPhoneVerifiedEvent) ActionType() string {
	return HumanPhoneVerifiedType
}

func HumanPhoneVerifiedEventFromStorage(event *eventstore.StorageEvent) (e *HumanPhoneVerifiedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-ycRBi", "Errors.Invalid.Event.Type")
	}

	return &HumanPhoneVerifiedEvent{
		StorageEvent: event,
	}, nil
}
