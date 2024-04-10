package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type humanPhoneChangedPayload struct {
	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`
}

type HumanPhoneChangedEvent humanPhoneChangedEvent
type humanPhoneChangedEvent = eventstore.Event[humanPhoneChangedPayload]

func HumanPhoneChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanPhoneChangedEvent, error) {
	event, err := eventstore.EventFromStorage[humanPhoneChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanPhoneChangedEvent)(event), nil
}
