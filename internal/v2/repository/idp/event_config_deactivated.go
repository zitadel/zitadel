package idp

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type ConfigDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewConfigDeactivatedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *ConfigDeactivatedEvent {

	return &ConfigDeactivatedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *ConfigDeactivatedEvent) Data() interface{} {
	return e
}

func ConfigDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
