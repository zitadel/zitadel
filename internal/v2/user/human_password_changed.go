package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type humanPasswordChangedPayload struct {
	// New events only use EncodedHash. However, the secret field
	// is preserved to handle events older than the switch to Passwap.
	// Secret            *crypto.CryptoValue `json:"secret,omitempty"`
	EncodedHash       string `json:"encodedHash,omitempty"`
	ChangeRequired    bool   `json:"changeRequired"`
	UserAgentID       string `json:"userAgentID,omitempty"`
	TriggeredAtOrigin string `json:"triggerOrigin,omitempty"`
}

type HumanPasswordChangedEvent humanPasswordChangedEvent
type humanPasswordChangedEvent = eventstore.Event[humanPasswordChangedPayload]

func HumanPasswordChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*HumanPasswordChangedEvent, error) {
	event, err := eventstore.EventFromStorage[humanPasswordChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*HumanPasswordChangedEvent)(event), nil
}
