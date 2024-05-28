package user

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type humanPasswordChangedPayload struct {
	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	Secret            *crypto.CryptoValue `json:"secret,omitempty"`
	EncodedHash       string              `json:"encodedHash,omitempty"`
	ChangeRequired    bool                `json:"changeRequired"`
	UserAgentID       string              `json:"userAgentID,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

type HumanPasswordChangedEvent eventstore.Event[humanPasswordChangedPayload]

const HumanPasswordChangedType = humanPrefix + ".password.changed"

var _ eventstore.TypeChecker = (*HumanPasswordChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *HumanPasswordChangedEvent) ActionType() string {
	return HumanPasswordChangedType
}

func HumanPasswordChangedEventFromStorage(event *eventstore.StorageEvent) (e *HumanPasswordChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-Fx5tr", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[humanPasswordChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
