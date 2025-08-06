package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/metadata"
)

const (
	MetadataSetType        = projectEventTypePrefix + metadata.SetEventType
	MetadataRemovedType    = projectEventTypePrefix + metadata.RemovedEventType
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

func MetadataSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
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

func MetadataRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := metadata.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetadataRemovedEvent{RemovedEvent: *e.(*metadata.RemovedEvent)}, nil
}
