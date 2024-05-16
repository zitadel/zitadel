package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type machineSecretHashUpdatedPayload struct {
	HashedSecret string `json:"hashedSecret,omitempty"`
}

type MachineSecretHashUpdatedEvent machineSecretHashUpdatedEvent
type machineSecretHashUpdatedEvent = eventstore.StorageEvent[machineSecretHashUpdatedPayload]

func MachineSecretHashUpdatedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*MachineSecretHashUpdatedEvent, error) {
	event, err := eventstore.EventFromStorage[machineSecretHashUpdatedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*MachineSecretHashUpdatedEvent)(event), nil
}
