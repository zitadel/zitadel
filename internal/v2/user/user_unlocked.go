package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserUnlockedEvent userUnlockedEvent
type userUnlockedEvent = eventstore.StorageEvent[struct{}]

func UserUnlockedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*UserUnlockedEvent, error) {
	event, err := eventstore.EventFromStorage[userUnlockedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserUnlockedEvent)(event), nil
}
