package user

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/v2/eventstore"
)

type machineAddedPayload struct {
	Username        string               `json:"userName"`
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	AccessTokenType domain.OIDCTokenType `json:"accessTokenType,omitempty"`
}

type MachineAddedEvent machineAddedEvent
type machineAddedEvent = eventstore.StorageEvent[machineAddedPayload]

func MachineAddedEventFromStorage(e *eventstore.StorageEvent[eventstore.StoragePayload]) (*MachineAddedEvent, error) {
	event, err := eventstore.EventFromStorage[machineAddedEvent](e)
	if err != nil {
		return nil, err
	}
	return (*MachineAddedEvent)(event), nil
}
