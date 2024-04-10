package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanPhoneVerifiedEvent humanPhoneVerifiedEvent
type humanPhoneVerifiedEvent = eventstore.Event[struct{}]

func HumanPhoneVerifiedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanPhoneVerifiedEvent, error) {
	event, err := eventstore.EventFromStorage[humanPhoneVerifiedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanPhoneVerifiedEvent)(event), nil
}
