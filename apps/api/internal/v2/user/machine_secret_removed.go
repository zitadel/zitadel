package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type MachineSecretRemovedEvent eventstore.Event[eventstore.EmptyPayload]

const MachineSecretRemovedType = machinePrefix + ".secret.removed"

var _ eventstore.TypeChecker = (*MachineSecretRemovedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *MachineSecretRemovedEvent) ActionType() string {
	return MachineSecretRemovedType
}

func MachineSecretRemovedEventFromStorage(event *eventstore.StorageEvent) (e *MachineSecretRemovedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-SMtct", "Errors.Invalid.Event.Type")
	}

	return &MachineSecretRemovedEvent{
		StorageEvent: event,
	}, nil
}
