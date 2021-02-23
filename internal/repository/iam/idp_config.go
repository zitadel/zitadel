package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

const (
	IDPConfigAddedEventType       eventstore.EventType = "iam.idp.config.added"
	IDPConfigChangedEventType     eventstore.EventType = "iam.idp.config.changed"
	IDPConfigRemovedEventType     eventstore.EventType = "iam.idp.config.removed"
	IDPConfigDeactivatedEventType eventstore.EventType = "iam.idp.config.deactivated"
	IDPConfigReactivatedEventType eventstore.EventType = "iam.idp.config.reactivated"
)

type IDPConfigAddedEvent struct {
	idpconfig.IDPConfigAddedEvent
}

func NewIDPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	configID,
	name string,
	configType domain.IDPConfigType,
	stylingType domain.IDPConfigStylingType,
) *IDPConfigAddedEvent {

	return &IDPConfigAddedEvent{
		IDPConfigAddedEvent: *idpconfig.NewIDPConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPConfigAddedEventType,
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
	aggregate *eventstore.Aggregate,
	configID,
	oldName string,
	changes []idpconfig.IDPConfigChanges,
) (*IDPConfigChangedEvent, error) {
	changeEvent, err := idpconfig.NewIDPConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IDPConfigChangedEventType),
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
	aggregate *eventstore.Aggregate,
	configID,
	name string,
) *IDPConfigRemovedEvent {
	return &IDPConfigRemovedEvent{
		IDPConfigRemovedEvent: *idpconfig.NewIDPConfigRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPConfigRemovedEventType,
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
	aggregate *eventstore.Aggregate,
	configID string,
) *IDPConfigDeactivatedEvent {

	return &IDPConfigDeactivatedEvent{
		IDPConfigDeactivatedEvent: *idpconfig.NewIDPConfigDeactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
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
	aggregate *eventstore.Aggregate,
	configID string,
) *IDPConfigReactivatedEvent {

	return &IDPConfigReactivatedEvent{
		IDPConfigReactivatedEvent: *idpconfig.NewIDPConfigReactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
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
