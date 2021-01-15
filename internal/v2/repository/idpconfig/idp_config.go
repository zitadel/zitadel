package idpconfig

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

type IDPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID    string                      `json:"idpConfigId"`
	Name        string                      `json:"name,omitempty"`
	Typ         domain.IDPConfigType        `json:"idpType,omitempty"`
	StylingType domain.IDPConfigStylingType `json:"stylingType,omitempty"`
}

func NewIDPConfigAddedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
	configType domain.IDPConfigType,
	stylingType domain.IDPConfigStylingType,
) *IDPConfigAddedEvent {

	return &IDPConfigAddedEvent{
		BaseEvent:   *base,
		ConfigID:    configID,
		Name:        name,
		StylingType: stylingType,
		Typ:         configType,
	}
}

func (e *IDPConfigAddedEvent) Data() interface{} {
	return e
}

func IDPConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IDPConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID    string                       `json:"idpConfigId"`
	Name        *string                      `json:"name,omitempty"`
	StylingType *domain.IDPConfigStylingType `json:"stylingType,omitempty"`
}

func (e *IDPConfigChangedEvent) Data() interface{} {
	return e
}

func NewIDPConfigChangedEvent(
	base *eventstore.BaseEvent,
	configID string,
	changes []IDPConfigChanges,
) *IDPConfigChangedEvent {
	changeEvent := &IDPConfigChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent
}

type IDPConfigChanges func(*IDPConfigChangedEvent)

func ChangeName(name string) func(*IDPConfigChangedEvent) {
	return func(e *IDPConfigChangedEvent) {
		e.Name = &name
	}
}

func ChangeStyleType(styleType domain.IDPConfigStylingType) func(*IDPConfigChangedEvent) {
	return func(e *IDPConfigChangedEvent) {
		e.StylingType = &styleType
	}
}

func IDPConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IDPConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewIDPConfigDeactivatedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *IDPConfigDeactivatedEvent {

	return &IDPConfigDeactivatedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *IDPConfigDeactivatedEvent) Data() interface{} {
	return e
}

func IDPConfigDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IDPConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewIDPConfigReactivatedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *IDPConfigReactivatedEvent {

	return &IDPConfigReactivatedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *IDPConfigReactivatedEvent) Data() interface{} {
	return e
}

func IDPConfigReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IDPConfigReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `idpConfigId`
}

func NewIDPConfigRemovedEvent(
	base *eventstore.BaseEvent,
	configID string,
) *IDPConfigRemovedEvent {

	return &IDPConfigRemovedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
	}
}

func (e *IDPConfigRemovedEvent) Data() interface{} {
	return e
}

func IDPConfigRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &IDPConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
