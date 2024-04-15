package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type removedPayload struct {
	Name string `json:"domain"`
}

type RemovedEvent removedEvent
type removedEvent = eventstore.Event[removedPayload]

func RemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*RemovedEvent, error) {
	event, err := eventstore.EventFromStorage[removedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*RemovedEvent)(event), nil
}

func (e *RemovedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "domain.removed")
}
