package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	// TODO: use same logic as in [strings.Builder] to get rid of the following line
	Changed ChangedEvent
)

type changedPayload struct {
	Name *string `json:"name,omitempty"`
}

type ChangedEvent changedEvent
type changedEvent = eventstore.Event[changedPayload]

func ChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*ChangedEvent, error) {
	event, err := eventstore.EventFromStorage[changedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*ChangedEvent)(event), nil
}

func (e ChangedEvent) IsType(typ string) bool {
	return typ == "org.changed"
}
