package idp

import "github.com/caos/zitadel/internal/eventstore/v2"

type ConfigAddedEvent struct {
	eventstore.BaseEvent

	ID          string      `json:"idpConfigId"`
	Name        string      `json:"name"`
	Type        ConfigType  `json:"idpType,omitempty"`
	StylingType StylingType `json:"stylingType,omitempty"`
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
	configType ConfigType,
	stylingType StylingType,
) *ConfigAddedEvent {

	return &ConfigAddedEvent{
		BaseEvent:   *base,
		ID:          configID,
		Name:        name,
		StylingType: stylingType,
		Type:        configType,
	}
}

func (e *ConfigAddedEvent) CheckPrevious() bool {
	return true
}

func (e *ConfigAddedEvent) Data() interface{} {
	return e
}

type ConfigType uint32

type StylingType uint32
