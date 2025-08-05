package user

import (
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type machineSecretSetPayload struct {
	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string              `json:"hashedSecret,omitempty"`
}

type MachineSecretHashSetEvent eventstore.Event[machineSecretSetPayload]

const MachineSecretHashSetType = machinePrefix + ".secret.set"

var _ eventstore.TypeChecker = (*MachineSecretHashSetEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *MachineSecretHashSetEvent) ActionType() string {
	return MachineSecretHashSetType
}

func MachineSecretHashSetEventFromStorage(event *eventstore.StorageEvent) (e *MachineSecretHashSetEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-DzycT", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[machineSecretSetPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &MachineSecretHashSetEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
