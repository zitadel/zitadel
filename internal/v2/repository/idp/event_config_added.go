package idp

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type ConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID    string      `json:"idpConfigId"`
	Name        string      `json:"name,omitempty"`
	Typ         ConfigType  `json:"idpType,omitempty"`
	StylingType StylingType `json:"stylingType,omitempty"`
}

func NewConfigAddedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
	configType ConfigType,
	stylingType StylingType,
) *ConfigAddedEvent {

	return &ConfigAddedEvent{
		BaseEvent:   *base,
		ConfigID:    configID,
		Name:        name,
		StylingType: stylingType,
		Typ:         configType,
	}
}

func (e *ConfigAddedEvent) CheckPrevious() bool {
	return true
}

func (e *ConfigAddedEvent) Data() interface{} {
	return e
}

func ConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
