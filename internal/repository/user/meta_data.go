package user

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/metadata"
)

const (
	MetaDataSetType     = userEventTypePrefix + metadata.SetEventType
	MetaDataRemovedType = userEventTypePrefix + metadata.RemovedEventType
)

type MetaDataSetEvent struct {
	metadata.SetEvent
}

func NewMetaDataSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, key, value string) *MetaDataSetEvent {
	return &MetaDataSetEvent{
		SetEvent: *metadata.NewSetEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MetaDataSetType),
			key,
			value),
	}
}

func MetaDataSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := metadata.SetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetaDataSetEvent{SetEvent: *e.(*metadata.SetEvent)}, nil
}

type MetaDataRemovedEvent struct {
	metadata.RemovedEvent
}

func NewMetaDataRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, key string) *MetaDataRemovedEvent {
	return &MetaDataRemovedEvent{
		RemovedEvent: *metadata.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MetaDataRemovedType),
			key),
	}
}

func MetaDataRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := metadata.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MetaDataRemovedEvent{RemovedEvent: *e.(*metadata.RemovedEvent)}, nil
}
