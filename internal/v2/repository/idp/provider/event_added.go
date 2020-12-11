package provider

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	AddedEventType = "idpprovider.added"
)

type AddedEvent struct {
	eventstore.BaseEvent

	IDPConfigID     string `json:"idpConfigId"`
	IDPProviderType Type   `json:"idpProviderType"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
	idpProviderType Type,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent:       *base,
		IDPConfigID:     idpConfigID,
		IDPProviderType: idpProviderType,
	}
}

func AddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-bfNnp", "Errors.Internal")
	}

	return e, nil
}
