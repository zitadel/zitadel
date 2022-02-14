package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	logType = ".log"
)

var (
	DebugNotificationProviderLogAddedEventType    = iamEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderAdded
	DebugNotificationProviderLogChangedEventType  = iamEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderChanged
	DebugNotificationProviderLogEnabledEventType  = iamEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderEnabled
	DebugNotificationProviderLogDisabledEventType = iamEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderDisabled
	DebugNotificationProviderLogRemovedEventType  = iamEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderRemoved
)

type DebugNotificationProviderLogAddedEvent struct {
	settings.DebugNotificationProviderAddedEvent
}

func NewDebugNotificationProviderLogAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	compact bool,
) *DebugNotificationProviderLogAddedEvent {
	return &DebugNotificationProviderLogAddedEvent{
		DebugNotificationProviderAddedEvent: *settings.NewDebugNotificationProviderAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderLogAddedEventType),
			compact),
	}
}

func DebugNotificationProviderLogAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogAddedEvent{DebugNotificationProviderAddedEvent: *e.(*settings.DebugNotificationProviderAddedEvent)}, nil
}

type DebugNotificationProviderLogChangedEvent struct {
	settings.DebugNotificationProviderChangedEvent
}

func NewDebugNotificationProviderLogChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []settings.DebugNotificationProviderChanges,
) (*DebugNotificationProviderLogChangedEvent, error) {
	changedEvent, err := settings.NewDebugNotificationProviderChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DebugNotificationProviderLogChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &DebugNotificationProviderLogChangedEvent{DebugNotificationProviderChangedEvent: *changedEvent}, nil
}

func DebugNotificationProviderLogChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogChangedEvent{DebugNotificationProviderChangedEvent: *e.(*settings.DebugNotificationProviderChangedEvent)}, nil
}

type DebugNotificationProviderLogEnabledEvent struct {
	settings.DebugNotificationProviderEnabledEvent
}

func NewDebugNotificationProviderLogEnabledEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DebugNotificationProviderLogEnabledEvent {
	return &DebugNotificationProviderLogEnabledEvent{
		DebugNotificationProviderEnabledEvent: *settings.NewDebugNotificationProviderEnabledEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderLogEnabledEventType),
		),
	}
}

func DebugNotificationProviderLogEnabledEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderEnabledEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogEnabledEvent{DebugNotificationProviderEnabledEvent: *e.(*settings.DebugNotificationProviderEnabledEvent)}, nil
}

type DebugNotificationProviderLogDisabledEvent struct {
	settings.DebugNotificationProviderDisabledEvent
}

func NewDebugNotificationProviderLogDisabledEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DebugNotificationProviderLogDisabledEvent {
	return &DebugNotificationProviderLogDisabledEvent{
		DebugNotificationProviderDisabledEvent: *settings.NewDebugNotificationProviderDisabledEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderLogDisabledEventType),
		),
	}
}

func DebugNotificationProviderLogDisabledEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderDisabledEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogDisabledEvent{DebugNotificationProviderDisabledEvent: *e.(*settings.DebugNotificationProviderDisabledEvent)}, nil
}

type DebugNotificationProviderLogRemovedEvent struct {
	settings.DebugNotificationProviderRemovedEvent
}

func NewDebugNotificationProviderLogRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DebugNotificationProviderLogRemovedEvent {
	return &DebugNotificationProviderLogRemovedEvent{
		DebugNotificationProviderRemovedEvent: *settings.NewDebugNotificationProviderRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderLogRemovedEventType),
		),
	}
}

func DebugNotificationProviderLogRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogRemovedEvent{DebugNotificationProviderRemovedEvent: *e.(*settings.DebugNotificationProviderRemovedEvent)}, nil
}
