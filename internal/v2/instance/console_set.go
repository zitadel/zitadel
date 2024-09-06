package instance

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const ConsoleSetType = eventTypePrefix + "iam.console.set"

type consoleSetPayload struct {
	ClientID string `json:"clientId"`
	AppID    string `json:"appId"`
}

type ConsoleSetEvent eventstore.Event[consoleSetPayload]

var _ eventstore.TypeChecker = (*ConsoleSetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *ConsoleSetEvent) ActionType() string {
	return ConsoleSetType
}

func ConsoleSetEventFromStorage(event *eventstore.StorageEvent) (e *ConsoleSetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "INSTA-wP2Ie", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[consoleSetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &ConsoleSetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
