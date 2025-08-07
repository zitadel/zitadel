package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/settings"
)

const (
	logType = ".log"
)

var (
	DebugNotificationProviderLogAddedEventType    = instanceEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderAdded
	DebugNotificationProviderLogChangedEventType  = instanceEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderChanged
	DebugNotificationProviderLogEnabledEventType  = instanceEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderEnabled
	DebugNotificationProviderLogDisabledEventType = instanceEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderDisabled
	DebugNotificationProviderLogRemovedEventType  = instanceEventTypePrefix + settings.DebugNotificationPrefix + logType + settings.DebugNotificationProviderRemoved
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

func DebugNotificationProviderLogAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
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

func DebugNotificationProviderLogChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogChangedEvent{DebugNotificationProviderChangedEvent: *e.(*settings.DebugNotificationProviderChangedEvent)}, nil
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

func DebugNotificationProviderLogRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderLogRemovedEvent{DebugNotificationProviderRemovedEvent: *e.(*settings.DebugNotificationProviderRemovedEvent)}, nil
}
