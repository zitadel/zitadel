package provider

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	RemovedEventType = "idpprovider.removed"
)

type RemovedEvent struct {
	eventstore.BaseEvent

	IDPConfigID string `json:"idpConfigId"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
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

func RemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "PROVI-6H0KQ", "Errors.Internal")
	}

	return e, nil
}
