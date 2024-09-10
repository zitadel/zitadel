package keypair

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const AddedType = eventTypePrefix + "added"

type addedPayload struct {
	Usage      crypto.KeyUsage `json:"usage"`
	Algorithm  string          `json:"algorithm"`
	PrivateKey *Key            `json:"privateKey"`
	PublicKey  *Key            `json:"publicKey"`
}

type AddedEvent eventstore.Event[addedPayload]

var _ eventstore.TypeChecker = (*AddedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *AddedEvent) ActionType() string {
	return AddedType
}

func AddedEventFromStorage(event *eventstore.StorageEvent) (e *AddedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "KEYPA-q2Ozb", "Errors.Invalid.Event.Type")
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
