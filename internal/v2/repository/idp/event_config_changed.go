package idp

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type ConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID    string      `json:"idpConfigId"`
	Name        string      `json:"name,omitempty"`
	StylingType StylingType `json:"stylingType,omitempty"`
}

func NewConfigChangedEvent(
	base *eventstore.BaseEvent,
	current *ConfigWriteModel,
	name string,
	stylingType StylingType,
) (*ConfigChangedEvent, error) {

	change := &ConfigChangedEvent{
		BaseEvent: *base,
		ConfigID:  current.ConfigID,
	}
	hasChanged := false

	if current.Name != name {
		change.Name = name
		hasChanged = true
	}
	if stylingType != current.StylingType {
		change.StylingType = stylingType
		hasChanged = true
	}

	if !hasChanged {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-UBJbB", "Errors.NoChanges")
	}

	return change, nil
}

func (e *ConfigChangedEvent) Data() interface{} {
	return e
}

func ConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
