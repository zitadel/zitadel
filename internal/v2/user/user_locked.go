package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserLockedEvent userLockedEvent
type userLockedEvent = eventstore.Event[struct{}]

func UserLockedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UserLockedEvent, error) {
	event, err := eventstore.EventFromStorage[userLockedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserLockedEvent)(event), nil
}
