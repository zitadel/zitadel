package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type usernameChangedPayload struct {
	Username string `json:"userName"`
}

type UsernameChangedEvent eventstore.Event[usernameChangedPayload]

const UsernameChangedType = AggregateType + ".username.changed"

var _ eventstore.TypeChecker = (*UsernameChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *UsernameChangedEvent) ActionType() string {
	return UsernameChangedType
}

func UsernameChangedEventFromStorage(event *eventstore.StorageEvent) (e *UsernameChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-hCGsh", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[usernameChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &UsernameChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
