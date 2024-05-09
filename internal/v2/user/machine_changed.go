package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type machineChangedPayload struct {
	Name            *string               `json:"name,omitempty"`
	Description     *string               `json:"description,omitempty"`
	AccessTokenType *domain.OIDCTokenType `json:"accessTokenType,omitempty"`
}

type MachineChangedEvent machineChangedEvent
type machineChangedEvent = eventstore.Event[machineChangedPayload]

func MachineChangedEventFromStorage(e *eventstore.Event[eventstore.StoragePayload]) (*MachineChangedEvent, error) {
	event, err := eventstore.EventFromStorage[machineChangedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*MachineChangedEvent)(event), nil
}
