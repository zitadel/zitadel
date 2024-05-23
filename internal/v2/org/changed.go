package org

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const ChangedType = eventTypePrefix + "changed"

type changedPayload struct {
	Name string `json:"name"`
}

type ChangedEvent eventstore.Event[changedPayload]

var _ eventstore.TypeChecker = (*ChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *ChangedEvent) ActionType() string {
	return ChangedType
}

func ChangedEventFromStorage(event *eventstore.StorageEvent) (e *ChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "ORG-pzOfP", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[changedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
