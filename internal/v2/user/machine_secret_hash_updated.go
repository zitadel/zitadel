package user

import (
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type machineSecretHashUpdatedPayload struct {
	HashedSecret string `json:"hashedSecret,omitempty"`
}

type MachineSecretHashUpdatedEvent eventstore.Event[machineSecretHashUpdatedPayload]

const MachineSecretHashUpdatedType = machinePrefix + ".secret.updated"

var _ eventstore.TypeChecker = (*MachineSecretHashUpdatedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *MachineSecretHashUpdatedEvent) ActionType() string {
	return MachineSecretHashUpdatedType
}

func MachineSecretHashUpdatedEventFromStorage(event *eventstore.StorageEvent) (e *MachineSecretHashUpdatedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-y41RK", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[machineSecretHashUpdatedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &MachineSecretHashUpdatedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
