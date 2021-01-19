package idpconfig

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	uniqueIDPConfigNameTable = "unique_idp_config_names"
)

type IDPConfigNameUniqueConstraint struct {
	tableName     string
	idpConfigName string
	action        eventstore.UniqueConstraintAction
}

func NewAddIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *IDPConfigNameUniqueConstraint {
	return &IDPConfigNameUniqueConstraint{
		tableName:     uniqueIDPConfigNameTable,
		idpConfigName: idpConfigName + resourceOwner,
		action:        eventstore.UniqueConstraintAdd,
	}
}

func NewRemoveIDPConfigNameUniqueConstraint(idpConfigName, resourceOwner string) *IDPConfigNameUniqueConstraint {
	return &IDPConfigNameUniqueConstraint{
		tableName:     uniqueIDPConfigNameTable,
		idpConfigName: idpConfigName + resourceOwner,
		action:        eventstore.UniqueConstraintRemoved,
	}
}

func (e *IDPConfigNameUniqueConstraint) TableName() string {
	return e.tableName
}

func (e *IDPConfigNameUniqueConstraint) UniqueField() string {
	return e.idpConfigName
}

func (e *IDPConfigNameUniqueConstraint) Action() eventstore.UniqueConstraintAction {
	return e.action
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

func (e *IDPConfigAddedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return []eventstore.EventUniqueConstraint{NewAddIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner())}
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

	ConfigID    string                      `json:"idpConfigId"`
	Name        string                      `json:"name,omitempty"`
	StylingType domain.IDPConfigStylingType `json:"stylingType,omitempty"`
}

func (e *IDPConfigChangedEvent) Data() interface{} {
	return e
}

func (e *IDPConfigChangedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
}

func NewIDPConfigChangedEvent(
	base *eventstore.BaseEvent,
) *IDPConfigChangedEvent {
	return &IDPConfigChangedEvent{
		BaseEvent: *base,
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

func (e *IDPConfigDeactivatedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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

func (e *IDPConfigReactivatedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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

	ConfigID string `idpConfigId`
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

func (e *IDPConfigRemovedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return []eventstore.EventUniqueConstraint{NewRemoveIDPConfigNameUniqueConstraint(e.Name, e.ResourceOwner())}
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
