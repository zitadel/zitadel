package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserRemovedEvent userRemovedEvent
type userRemovedEvent = eventstore.Event[struct{}]

func UserRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UserRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[userRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserRemovedEvent)(event), nil
}
