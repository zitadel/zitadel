package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserDeactivatedEvent userDeactivatedEvent
type userDeactivatedEvent = eventstore.Event[struct{}]

func UserDeactivatedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UserDeactivatedEvent, error) {
	event, err := eventstore.EventFromStorage[userDeactivatedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserDeactivatedEvent)(event), nil
}
