package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/idp"
	"github.com/caos/zitadel/internal/v2/repository/idp/oidc"
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
		case *idp.ConfigAddedEvent,
			*idp.ConfigChangedEvent,
			*idp.ConfigDeactivatedEvent,
			*idp.ConfigReactivatedEvent,
			*idp.ConfigRemovedEvent,
			*oidc.ConfigAddedEvent,
			*oidc.ConfigChangedEvent:

			rm.ConfigReadModel.AppendEvents(e)
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
