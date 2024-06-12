package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type machineChangedPayload struct {
	Name            *string               `json:"name,omitempty"`
	Description     *string               `json:"description,omitempty"`
	AccessTokenType *domain.OIDCTokenType `json:"accessTokenType,omitempty"`
}

type MachineChangedEvent eventstore.Event[machineChangedPayload]

const MachineChangedType = machinePrefix + ".changed"

var _ eventstore.TypeChecker = (*MachineChangedEvent)(nil)

// ActionType implements eventstore.Typer.
func (c *MachineChangedEvent) ActionType() string {
	return MachineChangedType
}

func MachineChangedEventFromStorage(event *eventstore.StorageEvent) (e *MachineChangedEvent, _ error) {
	if event.Type != e.ActionType() {
		return nil, zerrors.ThrowInvalidArgument(nil, "USER-JHwNs", "Errors.Invalid.Event.Type")
	}

	payload, err := eventstore.UnmarshalPayload[machineChangedPayload](event.Payload)
	if err != nil {
		return nil, err
	}

	return &MachineChangedEvent{
		StorageEvent: event,
		Payload:      payload,
	}, nil
}
