package user

import (
	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/v2/eventstore"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

type humanPhoneChangedPayload struct {
	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`
}

type HumanPhoneChangedEvent eventstore.Event[humanPhoneChangedPayload]

const HumanPhoneChangedType = humanPrefix + ".phone.changed"

var _ eventstore.TypeChecker = (*HumanPhoneChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanPhoneChangedEvent) ActionType() string {
	return HumanPhoneChangedType
}

func HumanPhoneChangedEventFromStorage(event *eventstore.StorageEvent) (e *HumanPhoneChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-d6hGS", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanPhoneChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanPhoneChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
