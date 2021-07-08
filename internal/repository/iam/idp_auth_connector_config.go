package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

const (
	IDPAuthConnectorConfigAddedEventType        eventstore.EventType = "iam.idp." + idpconfig.AuthConnectorConfigAddedEventType
	IDPAuthConnectorConfigChangedEventType      eventstore.EventType = "iam.idp." + idpconfig.AuthConnectorConfigChangedEventType
	IDPAuthConnectorMachineUserRemovedEventType eventstore.EventType = "iam.idp." + idpconfig.AuthConnectorMachineUserRemovedEventType
)

type IDPAuthConnectorConfigAddedEvent struct {
	idpconfig.AuthConnectorConfigAddedEvent
}

func NewIDPAuthConnectorConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID,
	baseURL,
	providerID,
	machineID string,
) *IDPAuthConnectorConfigAddedEvent {

	return &IDPAuthConnectorConfigAddedEvent{
		AuthConnectorConfigAddedEvent: *idpconfig.NewAuthConnectorConfigAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPAuthConnectorConfigAddedEventType,
			),
			idpConfigID,
			baseURL,
			providerID,
			machineID,
		),
	}
}

func IDPAuthConnectorConfigAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.AuthConnectorConfigAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPAuthConnectorConfigAddedEvent{AuthConnectorConfigAddedEvent: *e.(*idpconfig.AuthConnectorConfigAddedEvent)}, nil
}

type IDPAuthConnectorConfigChangedEvent struct {
	idpconfig.AuthConnectorConfigChangedEvent
}

func NewIDPAuthConnectorConfigChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
	changes []idpconfig.AuthConnectorConfigChanges,
) (*IDPAuthConnectorConfigChangedEvent, error) {
	changeEvent, err := idpconfig.NewAuthConnectorConfigChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			IDPAuthConnectorConfigChangedEventType),
		idpConfigID,
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &IDPAuthConnectorConfigChangedEvent{AuthConnectorConfigChangedEvent: *changeEvent}, nil
}

func IDPAuthConnectorConfigChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.AuthConnectorConfigChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPAuthConnectorConfigChangedEvent{AuthConnectorConfigChangedEvent: *e.(*idpconfig.AuthConnectorConfigChangedEvent)}, nil
}

type IDPAuthConnectorMachineUserRemovedEvent struct {
	idpconfig.AuthConnectorMachineUserRemovedEvent
}

func NewIDPAuthConnectorMachineUserRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	idpConfigID string,
) *IDPAuthConnectorMachineUserRemovedEvent {

	return &IDPAuthConnectorMachineUserRemovedEvent{
		AuthConnectorMachineUserRemovedEvent: *idpconfig.NewAuthConnectorMachineUserRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				IDPAuthConnectorMachineUserRemovedEventType,
			),
			idpConfigID,
		),
	}
}

func IDPAuthConnectorMachineUserRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := idpconfig.AuthConnectorMachineUserRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &IDPAuthConnectorMachineUserRemovedEvent{AuthConnectorMachineUserRemovedEvent: *e.(*idpconfig.AuthConnectorMachineUserRemovedEvent)}, nil
}
