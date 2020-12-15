package idpprovider

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	AddedEventType   = "idpprovider.added"
	RemovedEventType = "idpprovider.removed"
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
		*base,
		idpConfigID,
		idpProviderType,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-bfNnp", "Errors.Internal")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent

	IDPConfigID string `json:"idpConfigId"`
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	idpConfigID string,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent:   *base,
		IDPConfigID: idpConfigID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-6H0KQ", "Errors.Internal")
	}

	return e, nil
}
