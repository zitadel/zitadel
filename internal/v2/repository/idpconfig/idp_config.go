package idpconfig

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	uniqueIDPConfigNameType = "idp_config_names"
)

func NewAddIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueIDPConfigNameType,
		idpConfigName+resourceOwner,
		"Errors.IDPConfig.AlreadyExists")
}

func NewRemoveIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *eventstore.EventUniqueConstraint {
	return eventstore.NewRemoveEventUniqueConstraint(
		uniqueIDPConfigNameType,
		idpConfigName+resourceOwner)
}

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

func (e *IDPConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewAddIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner())}
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

func (e *IDPConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewIDPConfigChangedEvent(
	base *eventstore.BaseEvent,
	configID string,
	changes []IDPConfigChanges,
) (*IDPConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IDPCONFIG-Dsg21", "Errors.NoChangesFound")
	}
	changeEvent := &IDPConfigChangedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
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

func (e *IDPConfigDeactivatedEvent) Data() interface{} {
	return e
}

func (e *IDPConfigDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
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

func (e *IDPConfigReactivatedEvent) Data() interface{} {
	return e
}

func (e *IDPConfigReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
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

	ConfigID string `json:"idpConfigId"`
	Name     string
}

func NewIDPConfigRemovedEvent(
	base *eventstore.BaseEvent,
	configID string,
	name string,
) *IDPConfigRemovedEvent {

	return &IDPConfigRemovedEvent{
		BaseEvent: *base,
		ConfigID:  configID,
		Name:      name,
	}
}

func (e *IDPConfigRemovedEvent) Data() interface{} {
	return e
}

func (e *IDPConfigRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return []*eventstore.EventUniqueConstraint{NewRemoveIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner())}
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
