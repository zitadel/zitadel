package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserUnlockedEvent userUnlockedEvent
type userUnlockedEvent = eventstore.Event[struct{}]

func UserUnlockedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UserUnlockedEvent, error) {
	event, err := eventstore.EventFromStorage[userUnlockedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserUnlockedEvent)(event), nil
}
