package idp

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type ConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewConfigRemovedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *ConfigRemovedEvent {

	return &ConfigRemovedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *ConfigRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *ConfigRemovedEvent) Data() interface{} {
	return e
}

func ConfigRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
