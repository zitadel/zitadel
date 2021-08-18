package user

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/metadata"
)

const (
	MetadataSetType        = userEventTypePrefix + metadata.SetEventType
	MetadataRemovedType    = userEventTypePrefix + metadata.RemovedEventType
	MetadataRemovedAllType = userEventTypePrefix + metadata.RemovedAllEventType
)

type MetadataSetEvent struct {
	metadata.SetEvent
}

func NewMetadataSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, key string, value []byte) *MetadataSetEvent {
	return &MetadataSetEvent{
		SetEvent: *metadata.NewSetEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MetadataSetType),
			key,
			value),
	}
}

func MetadataSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := metadata.SetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetadataSetEvent{SetEvent: *e.(*metadata.SetEvent)}, nil
}

type MetadataRemovedEvent struct {
	metadata.RemovedEvent
}

func NewMetadataRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, key string) *MetadataRemovedEvent {
	return &MetadataRemovedEvent{
		RemovedEvent: *metadata.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MetadataRemovedType),
			key),
	}
}

func MetadataRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := metadata.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetadataRemovedEvent{RemovedEvent: *e.(*metadata.RemovedEvent)}, nil
}

type MetadataRemovedAllEvent struct {
	metadata.RemovedAllEvent
}

func NewMetadataRemovedAllEvent(ctx context.Context, aggregate *eventstore.Aggregate) *MetadataRemovedAllEvent {
	return &MetadataRemovedAllEvent{
		RemovedAllEvent: *metadata.NewRemovedAllEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MetadataRemovedAllType),
		),
	}
}

func MetadataRemovedAllEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := metadata.RemovedAllEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetadataRemovedAllEvent{RemovedAllEvent: *e.(*metadata.RemovedAllEvent)}, nil
}
