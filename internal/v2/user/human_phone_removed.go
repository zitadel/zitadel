package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanPhoneRemovedEvent humanPhoneRemovedEvent
type humanPhoneRemovedEvent = eventstore.StorageEvent[struct{}]

func HumanPhoneRemovedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*HumanPhoneRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[humanPhoneRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanPhoneRemovedEvent)(event), nil
}
