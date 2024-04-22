package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	Deactivated DeactivatedEvent
)

type DeactivatedEvent deactivatedEvent
type deactivatedEvent = eventstore.Event[struct{}]

func DeactivatedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*DeactivatedEvent, error) {
	event, err := eventstore.EventFromStorage[deactivatedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*DeactivatedEvent)(event), nil
}

func (e DeactivatedEvent) IsType(typ string) bool {
	return typ == "org.deactivated"
}
