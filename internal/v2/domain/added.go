package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type addedPayload struct {
	Name string `json:"domain"`
}

type AddedEvent addedEvent
type addedEvent = eventstore.Event[addedPayload]

func AddedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*AddedEvent, error) {
	event, err := eventstore.EventFromStorage[addedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*AddedEvent)(event), nil
}

func (e *AddedEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "domain.added")
}
