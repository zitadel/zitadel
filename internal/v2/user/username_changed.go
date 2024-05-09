package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type usernameChangedPayload struct {
	Username string `json:"userName"`
}

type UsernameChangedEvent usernameChangedEvent
type usernameChangedEvent = eventstore.Event[usernameChangedPayload]

func UsernameChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UsernameChangedEvent, error) {
	event, err := eventstore.EventFromStorage[usernameChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UsernameChangedEvent)(event), nil
}
