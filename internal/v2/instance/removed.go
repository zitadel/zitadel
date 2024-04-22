package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type RemovedEvent removedEvent
type removedEvent = eventstore.Event[struct{}]

func RemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*RemovedEvent, error) {
	event, err := eventstore.EventFromStorage[removedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*RemovedEvent)(event), nil
}
