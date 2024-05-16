package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type humanEmailChangedPayload struct {
	Address domain.EmailAddress `json:"email,omitempty"`
}

type HumanEmailChangedEvent humanEmailChangedEvent
type humanEmailChangedEvent = eventstore.StorageEvent[humanEmailChangedPayload]

func HumanEmailChangedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*HumanEmailChangedEvent, error) {
	event, err := eventstore.EventFromStorage[humanEmailChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanEmailChangedEvent)(event), nil
}
