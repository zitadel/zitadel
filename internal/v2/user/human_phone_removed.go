package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanPhoneRemovedEvent humanPhoneRemovedEvent
type humanPhoneRemovedEvent = eventstore.Event[struct{}]

func HumanPhoneRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanPhoneRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[humanPhoneRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanPhoneRemovedEvent)(event), nil
}
