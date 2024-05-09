package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type HumanEmailVerifiedEvent humanEmailVerifiedEvent
type humanEmailVerifiedEvent = eventstore.Event[struct{}]

func HumanEmailVerifiedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanEmailVerifiedEvent, error) {
	event, err := eventstore.EventFromStorage[humanEmailVerifiedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanEmailVerifiedEvent)(event), nil
}
