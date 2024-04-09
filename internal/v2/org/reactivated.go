package org

import "github.com/zitadel/zitadel/internal/v2/eventstore"

var (
	Reactivated ReactivatedEvent
)

type ReactivatedEvent reactivatedEvent
type reactivatedEvent = eventstore.Event[struct{}]

func ReactivatedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*ReactivatedEvent, error) {
	event, err := eventstore.EventFromStorage[reactivatedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*ReactivatedEvent)(event), nil
}

func (e ReactivatedEvent) IsType(typ string) bool {
	return typ == "org.reactivated"
}
