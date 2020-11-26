package idp

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type ConfigReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewConfigReactivatedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *ConfigReactivatedEvent {

	return &ConfigReactivatedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *ConfigReactivatedEvent) CheckPrevious() bool {
	return true
}

func (e *ConfigReactivatedEvent) Data() interface{} {
	return e
}

func ConfigReactivatedEventMapper(event *repository.Event) (*ConfigReactivatedEvent, error) {
	e := &ConfigReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
