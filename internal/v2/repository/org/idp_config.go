package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/idpconfig"
)

const (
	IDPConfigAddedEventType       eventstore.EventType = "org.idp.config.added"
	IDPConfigChangedEventType     eventstore.EventType = "org.idp.config.changed"
	IDPConfigRemovedEventType     eventstore.EventType = "org.idp.config.removed"
	IDPConfigDeactivatedEventType eventstore.EventType = "org.idp.config.deactivated"
	IDPConfigReactivatedEventType eventstore.EventType = "org.idp.config.reactivated"
)

type IDPConfigAddedEvent struct {
	idpconfig.IDPConfigAddedEvent
}

func NewIDPConfigAddedEvent(
	ctx context.Context,
	resourceOwner,
	configID,
	name string,
	configType domain.IDPConfigType,
	stylingType domain.IDPConfigStylingType,
) *IDPConfigAddedEvent {

	return &IDPConfigAddedEvent{
		IDPConfigAddedEvent: *idpconfig.NewIDPConfigAddedEvent(
			eventstore.NewBaseEventForPushWithResourceOwner(
				ctx,
				IDPConfigAddedEventType,
				resourceOwner,
			),
			configID,
			name,
			configType,
			stylingType,
		),
	}
}

func IDPConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.IDPConfigAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigAddedEvent{IDPConfigAddedEvent: *e.(*idpconfig.IDPConfigAddedEvent)}, nil
}

type IDPConfigChangedEvent struct {
	idpconfig.IDPConfigChangedEvent
}

func NewIDPConfigChangedEvent(
	ctx context.Context,
	resourceOwner,
	configID,
	oldName string,
	changes []idpconfig.IDPConfigChanges,
) (*IDPConfigChangedEvent, error) {
	changeEvent, err := idpconfig.NewIDPConfigChangedEvent(
		eventstore.NewBaseEventForPushWithResourceOwner(ctx, IDPConfigChangedEventType, resourceOwner),
		configID,
		oldName,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &IDPConfigChangedEvent{IDPConfigChangedEvent: *changeEvent}, nil
}

func IDPConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.IDPConfigChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigChangedEvent{IDPConfigChangedEvent: *e.(*idpconfig.IDPConfigChangedEvent)}, nil
}

type IDPConfigRemovedEvent struct {
	idpconfig.IDPConfigRemovedEvent
}

func NewIDPConfigRemovedEvent(
	ctx context.Context,
	resourceOwner,
	configID,
	name string,
) *IDPConfigRemovedEvent {

	return &IDPConfigRemovedEvent{
		IDPConfigRemovedEvent: *idpconfig.NewIDPConfigRemovedEvent(
			eventstore.NewBaseEventForPushWithResourceOwner(
				ctx,
				IDPConfigRemovedEventType,
				resourceOwner,
			),
			configID,
			name,
		),
	}
}

func IDPConfigRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.IDPConfigRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigRemovedEvent{IDPConfigRemovedEvent: *e.(*idpconfig.IDPConfigRemovedEvent)}, nil
}

type IDPConfigDeactivatedEvent struct {
	idpconfig.IDPConfigDeactivatedEvent
}

func NewIDPConfigDeactivatedEvent(
	ctx context.Context,
	configID string,
) *IDPConfigDeactivatedEvent {

	return &IDPConfigDeactivatedEvent{
		IDPConfigDeactivatedEvent: *idpconfig.NewIDPConfigDeactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				IDPConfigDeactivatedEventType,
			),
			configID,
		),
	}
}

func IDPConfigDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.IDPConfigDeactivatedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigDeactivatedEvent{IDPConfigDeactivatedEvent: *e.(*idpconfig.IDPConfigDeactivatedEvent)}, nil
}

type IDPConfigReactivatedEvent struct {
	idpconfig.IDPConfigReactivatedEvent
}

func NewIDPConfigReactivatedEvent(
	ctx context.Context,
	configID string,
) *IDPConfigReactivatedEvent {

	return &IDPConfigReactivatedEvent{
		IDPConfigReactivatedEvent: *idpconfig.NewIDPConfigReactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				IDPConfigReactivatedEventType,
			),
			configID,
		),
	}
}

func IDPConfigReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.IDPConfigReactivatedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigReactivatedEvent{IDPConfigReactivatedEvent: *e.(*idpconfig.IDPConfigReactivatedEvent)}, nil
}
