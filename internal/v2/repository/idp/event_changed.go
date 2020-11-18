package idp

import (
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID          string      `json:"idpConfigId"`
	StylingType StylingType `json:"stylingType,omitempty"`

	hasChanged bool
}

func NewConfigChangedEvent(
	base *eventstore.BaseEvent,
	current *ConfigAggregate,
	changed *ConfigAggregate,
) (*ConfigChangedEvent, error) {

	change := &ConfigChangedEvent{
		BaseEvent: *base,
	}

	if current.ConfigID != changed.ConfigID {
		change.ID = changed.ConfigID
		change.hasChanged = true
	}

	if current.StylingType != changed.StylingType {
		change.StylingType = changed.StylingType
		change.hasChanged = true
	}

	if !change.hasChanged {
		return nil, errors.ThrowPreconditionFailed(nil, "IDP-UBJbB", "Errors.NoChanges")
	}

	return change, nil
}

func (e *ConfigChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ConfigChangedEvent) Data() interface{} {
	if e.current.Name != e.changed.Name {
		e.Name = e.changed.Name
	}
	return e
}
