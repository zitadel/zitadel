package instance

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	fileType = ".file"
)

var (
	DebugNotificationProviderFileAddedEventType   = iamEventTypePrefix + settings.DebugNotificationPrefix + fileType + settings.DebugNotificationProviderAdded
	DebugNotificationProviderFileChangedEventType = iamEventTypePrefix + settings.DebugNotificationPrefix + fileType + settings.DebugNotificationProviderChanged
	DebugNotificationProviderFileRemovedEventType = iamEventTypePrefix + settings.DebugNotificationPrefix + fileType + settings.DebugNotificationProviderRemoved
)

type DebugNotificationProviderFileAddedEvent struct {
	settings.DebugNotificationProviderAddedEvent
}

func NewDebugNotificationProviderFileAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	compact bool,
) *DebugNotificationProviderFileAddedEvent {
	return &DebugNotificationProviderFileAddedEvent{
		DebugNotificationProviderAddedEvent: *settings.NewDebugNotificationProviderAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderFileAddedEventType),
			compact),
	}
}

func DebugNotificationProviderFileAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderFileAddedEvent{DebugNotificationProviderAddedEvent: *e.(*settings.DebugNotificationProviderAddedEvent)}, nil
}

type DebugNotificationProviderFileChangedEvent struct {
	settings.DebugNotificationProviderChangedEvent
}

func NewDebugNotificationProviderFileChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []settings.DebugNotificationProviderChanges,
) (*DebugNotificationProviderFileChangedEvent, error) {
	changedEvent, err := settings.NewDebugNotificationProviderChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DebugNotificationProviderFileChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &DebugNotificationProviderFileChangedEvent{DebugNotificationProviderChangedEvent: *changedEvent}, nil
}

func DebugNotificationProviderFileChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderFileChangedEvent{DebugNotificationProviderChangedEvent: *e.(*settings.DebugNotificationProviderChangedEvent)}, nil
}

type DebugNotificationProviderFileRemovedEvent struct {
	settings.DebugNotificationProviderRemovedEvent
}

func NewDebugNotificationProviderFileRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DebugNotificationProviderFileRemovedEvent {
	return &DebugNotificationProviderFileRemovedEvent{
		DebugNotificationProviderRemovedEvent: *settings.NewDebugNotificationProviderRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DebugNotificationProviderFileRemovedEventType),
		),
	}
}

func DebugNotificationProviderFileRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := settings.DebugNotificationProviderRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DebugNotificationProviderFileRemovedEvent{DebugNotificationProviderRemovedEvent: *e.(*settings.DebugNotificationProviderRemovedEvent)}, nil
}
