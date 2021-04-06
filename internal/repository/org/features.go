package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/features"
)

var (
	FeaturesSetEventType     = orgEventTypePrefix + features.FeaturesSetEventType
	FeaturesRemovedEventType = orgEventTypePrefix + features.FeaturesRemovedEventType
)

type FeaturesSetEvent struct {
	features.FeaturesSetEvent
}

func NewFeaturesSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []features.FeaturesChanges,
) (*FeaturesSetEvent, error) {
	changedEvent, err := features.NewFeaturesSetEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FeaturesSetEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &FeaturesSetEvent{FeaturesSetEvent: *changedEvent}, nil
}

func FeaturesSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := features.FeaturesSetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &FeaturesSetEvent{FeaturesSetEvent: *e.(*features.FeaturesSetEvent)}, nil
}

type FeaturesRemovedEvent struct {
	features.FeaturesRemovedEvent
}

func NewFeaturesRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *FeaturesRemovedEvent {
	return &FeaturesRemovedEvent{
		FeaturesRemovedEvent: *features.NewFeaturesRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				FeaturesRemovedEventType),
		),
	}
}

func FeaturesRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := features.FeaturesRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &FeaturesRemovedEvent{FeaturesRemovedEvent: *e.(*features.FeaturesRemovedEvent)}, nil
}
