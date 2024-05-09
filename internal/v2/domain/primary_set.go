package domain

import (
	"strings"

	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type primarySetPayload struct {
	Name string `json:"domain"`
}

type PrimarySetEvent primarySetEvent
type primarySetEvent = eventstore.Event[primarySetPayload]

func PrimarySetEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*PrimarySetEvent, error) {
	event, err := eventstore.EventFromStorage[primarySetEvent](e)
	if err != nil {
		return nil, err
	}
	return (*PrimarySetEvent)(event), nil
}

func (e *PrimarySetEvent) HasTypeSuffix(typ string) bool {
	return strings.HasSuffix(typ, "domain.primary.set")
}
