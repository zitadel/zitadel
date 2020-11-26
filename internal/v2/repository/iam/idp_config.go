package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp"
)

const (
	IDPConfigAddedEventType       eventstore.EventType = "iam.idp.config.added"
	IDPConfigChangedEventType     eventstore.EventType = "iam.idp.config.changed"
	IDPConfigRemovedEventType     eventstore.EventType = "iam.idp.config.removed"
	IDPConfigDeactivatedEventType eventstore.EventType = "iam.idp.config.deactivated"
	IDPConfigReactivatedEventType eventstore.EventType = "iam.idp.config.reactivated"
)

type IDPConfigReadModel struct {
	idp.ConfigReadModel
}

func (rm *IDPConfigReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPConfigAddedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPConfigChangedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigChangedEvent)
		case *IDPConfigDeactivatedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigDeactivatedEvent)
		case *IDPConfigReactivatedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigReactivatedEvent)
		case *IDPConfigRemovedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigRemovedEvent)
		case *IDPOIDCConfigAddedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPOIDCConfigChangedEvent:
			rm.ConfigReadModel.AppendEvents(&e.ConfigChangedEvent)
		}
	}
}

type IDPConfigWriteModel struct {
	idp.ConfigWriteModel
}

func (rm *IDPConfigWriteModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *IDPConfigAddedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigAddedEvent)
		case *IDPConfigChangedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigChangedEvent)
		case *IDPConfigDeactivatedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigDeactivatedEvent)
		case *IDPConfigReactivatedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigReactivatedEvent)
		case *IDPConfigRemovedEvent:
			rm.ConfigWriteModel.AppendEvents(&e.ConfigRemovedEvent)
		case *idp.ConfigAddedEvent,
			*idp.ConfigChangedEvent,
			*idp.ConfigDeactivatedEvent,
			*idp.ConfigReactivatedEvent,
			*idp.ConfigRemovedEvent:

			rm.ConfigWriteModel.AppendEvents(e)
		}
	}
}

type IDPConfigAddedEvent struct {
	idp.ConfigAddedEvent
}

func NewIDPConfigAddedEvent(
	ctx context.Context,
	configID string,
	name string,
	configType idp.ConfigType,
	stylingType idp.StylingType,
) *IDPConfigAddedEvent {

	return &IDPConfigAddedEvent{
		ConfigAddedEvent: *idp.NewConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
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
	e, err := idp.ConfigAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigAddedEvent{ConfigAddedEvent: *e}, nil
}

type IDPConfigChangedEvent struct {
	idp.ConfigChangedEvent
}

func NewIDPConfigChangedEvent(
	ctx context.Context,
	current *IDPConfigWriteModel,
	configID string,
	name string,
	configType idp.ConfigType,
	stylingType idp.StylingType,
) (*IDPConfigChangedEvent, error) {
	event, err := idp.NewConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			IDPConfigChangedEventType,
		),
		&current.ConfigWriteModel,
		name,
		stylingType,
	)

	if err != nil {
		return nil, err
	}

	return &IDPConfigChangedEvent{
		ConfigChangedEvent: *event,
	}, nil
}

func IDPConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idp.ConfigChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigChangedEvent{ConfigChangedEvent: *e}, nil
}

type IDPConfigRemovedEvent struct {
	idp.ConfigRemovedEvent
}

func NewIDPConfigRemovedEvent(
	ctx context.Context,
	configID string,
) *IDPConfigRemovedEvent {

	return &IDPConfigRemovedEvent{
		ConfigRemovedEvent: *idp.NewConfigRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				IDPConfigRemovedEventType,
			),
			configID,
		),
	}
}

func IDPConfigRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idp.ConfigRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigRemovedEvent{ConfigRemovedEvent: *e}, nil
}

type IDPConfigDeactivatedEvent struct {
	idp.ConfigDeactivatedEvent
}

func NewIDPConfigDeactivatedEvent(
	ctx context.Context,
	configID string,
) *IDPConfigDeactivatedEvent {

	return &IDPConfigDeactivatedEvent{
		ConfigDeactivatedEvent: *idp.NewConfigDeactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				IDPConfigDeactivatedEventType,
			),
			configID,
		),
	}
}

func IDPConfigDeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idp.ConfigDeactivatedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigDeactivatedEvent{ConfigDeactivatedEvent: *e}, nil
}

type IDPConfigReactivatedEvent struct {
	idp.ConfigReactivatedEvent
}

func NewIDPConfigReactivatedEvent(
	ctx context.Context,
	configID string,
) *IDPConfigReactivatedEvent {

	return &IDPConfigReactivatedEvent{
		ConfigReactivatedEvent: *idp.NewConfigReactivatedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				IDPConfigReactivatedEventType,
			),
			configID,
		),
	}
}

func IDPConfigReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idp.ConfigReactivatedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPConfigReactivatedEvent{ConfigReactivatedEvent: *e}, nil
}
