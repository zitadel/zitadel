package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type machineSecretSetPayload struct {
	// New events only use EncodedHash. However, the ClientSecret field
	// is preserved to handle events older than the switch to Passwap.
	// ClientSecret *crypto.CryptoValue `json:"clientSecret,omitempty"`
	HashedSecret string `json:"hashedSecret,omitempty"`
}

type MachineSecretSetEvent machineSecretSetEvent
type machineSecretSetEvent = eventstore.StorageEvent[machineSecretSetPayload]

func MachineSecretSetEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*MachineSecretSetEvent, error) {
	event, err := eventstore.EventFromStorage[machineSecretSetEvent](e)
	if err != nil {
		return nil, err
	}
	return (*MachineSecretSetEvent)(event), nil
}
