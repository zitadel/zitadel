package idpconfig

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueIDPConfigNameType = "idp_config_names"
)

func NewAddIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueIDPConfigNameType,
		idpConfigName+resourceOwner,
		"Errors.IDPConfig.AlreadyExists")
}

func NewRemoveIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueIDPConfigNameType,
		idpConfigName+resourceOwner)
}

type IDPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID     string                      `json:"idpConfigId"`
	Name         string                      `json:"name,omitempty"`
	Typ          domain.IDPConfigType        `json:"idpType,omitempty"`
	StylingType  domain.IDPConfigStylingType `json:"stylingType,omitempty"`
	AutoRegister bool                        `json:"autoRegister,omitempty"`
}

func NewIDPConfigAddedEvent(
	base *eventstore.BaseEvent,
	configID,
	name string,
	configType domain.IDPConfigType,
	stylingType domain.IDPConfigStylingType,
	autoRegister bool,
) *IDPConfigAddedEvent {
	return &IDPConfigAddedEvent{
		BaseEvent:    *base,
		ConfigID:     configID,
		Name:         name,
		StylingType:  stylingType,
		Typ:          configType,
		AutoRegister: autoRegister,
	}
}

func (e *IDPConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *IDPConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddIDPConfigNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func IDPConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IDPConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID     string                       `json:"idpConfigId"`
	Name         *string                      `json:"name,omitempty"`
	StylingType  *domain.IDPConfigStylingType `json:"stylingType,omitempty"`
	AutoRegister *bool                        `json:"autoRegister,omitempty"`
	oldName      string                       `json:"-"`
}

func (e *IDPConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *IDPConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldName == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveIDPConfigNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
		NewAddIDPConfigNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	}
}

func NewIDPConfigChangedEvent(
	base *eventstore.BaseEvent,
	configID,
	oldName string,
	changes []IDPConfigChanges,
) (*IDPConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IDPCONFIG-Dsg21", "Errors.NoChangesFound")
	}
	changeEvent := &IDPConfigChangedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
		oldName:   oldName,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
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

func ChangeAutoRegister(autoRegister bool) func(*IDPConfigChangedEvent) {
	return func(e *IDPConfigChangedEvent) {
		e.AutoRegister = &autoRegister
	}
}

func IDPConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IDPConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `json:"idpConfigId"`
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

func (e *IDPConfigDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *IDPConfigDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func IDPConfigDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IDPConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `json:"idpConfigId"`
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

func (e *IDPConfigReactivatedEvent) Payload() interface{} {
	return e
}

func (e *IDPConfigReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func IDPConfigReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IDPConfigReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}

type IDPConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ConfigID string `json:"idpConfigId"`
	name     string
}

func NewIDPConfigRemovedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
) *IDPConfigRemovedEvent {

	return &IDPConfigRemovedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
		name:      name,
	}
}

func (e *IDPConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *IDPConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveIDPConfigNameUniqueConstraint(e.name, e.Aggregate().ResourceOwner)}
}

func IDPConfigRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &IDPConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "OIDC-plaBZ", "unable to unmarshal event")
	}

	return e, nil
}
