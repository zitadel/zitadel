package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const AddedType = eventTypePrefix + "added"

type addedPayload struct {
	Name string `json:"name"`
}

type AddedEvent eventstore.Event[addedPayload]

var _ eventstore.TypeChecker = (*AddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *AddedEvent) ActionType() string {
	return AddedType
}

func AddedEventFromStorage(event *eventstore.StorageEvent) (e *AddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-oRMBW", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[addedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
