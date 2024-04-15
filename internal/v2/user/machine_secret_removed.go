package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type MachineSecretRemovedEvent machineSecretRemovedEvent
type machineSecretRemovedEvent = eventstore.Event[struct{}]

func MachineSecretRemovedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*MachineSecretRemovedEvent, error) {
	event, err := eventstore.EventFromStorage[machineSecretRemovedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*MachineSecretRemovedEvent)(event), nil
}
