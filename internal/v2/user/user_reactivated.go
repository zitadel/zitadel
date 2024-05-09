package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type UserReactivatedEvent userReactivatedEvent
type userReactivatedEvent = eventstore.Event[struct{}]

func UserReactivatedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*UserReactivatedEvent, error) {
	event, err := eventstore.EventFromStorage[userReactivatedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*UserReactivatedEvent)(event), nil
}
