package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanInitialCodeSucceededEvent humanInitialCodeSucceededEvent
type humanInitialCodeSucceededEvent = eventstore.StorageEvent[struct{}]

func HumanInitialCodeSucceededEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*HumanInitialCodeSucceededEvent, error) {
	event, err := eventstore.EventFromStorage[humanInitialCodeSucceededEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanInitialCodeSucceededEvent)(event), nil
}
