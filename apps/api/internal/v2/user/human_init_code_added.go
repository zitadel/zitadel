package user

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type humanInitCodeAddedPayload struct {
	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
	AuthRequestID     string              `json:"authRequestID,omitempty"`
}

type HumanInitCodeAddedEvent eventstore.Event[humanInitCodeAddedPayload]

const HumanInitCodeAddedType = humanPrefix + ".initialization.code.added"

var _ eventstore.TypeChecker = (*HumanInitCodeAddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanInitCodeAddedEvent) ActionType() string {
	return HumanInitCodeAddedType
}

func HumanInitCodeAddedEventFromStorage(event *eventstore.StorageEvent) (e *HumanInitCodeAddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-XaGf6", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanInitCodeAddedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanInitCodeAddedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
