package debug_events

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	AddedEventType   = eventTypePrefix + "added"
	ChangedEventType = eventTypePrefix + "changed"
	RemovedEventType = eventTypePrefix + "removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ProjectionSleep      time.Duration `json:"projectionSleep,omitempty"`
	Blob                 *string       `json:"blob,omitempty"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, projectionSleep time.Duration, blob *string) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedEventType,
		),
		Blob:            blob,
		ProjectionSleep: projectionSleep,
	}
}

func DebugAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	debugAdded := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(debugAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal debug added")
	}

	return debugAdded, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ProjectionSleep      time.Duration `json:"projectionSleep,omitempty"`
	Blob                 *string       `json:"blob,omitempty"`
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, projectionSleep time.Duration, blob *string) *ChangedEvent {
	return &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ChangedEventType,
		),
		ProjectionSleep: projectionSleep,
		Blob:            blob,
	}
}

func DebugChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	debugChanged := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(debugChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-Bren2", "unable to unmarshal debug added")
	}

	return debugChanged, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ProjectionSleep      time.Duration `json:"projectionSleep,omitempty"`
}

func (e *RemovedEvent) Payload() interface{} {
	return nil
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, projectionSleep time.Duration) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RemovedEventType,
		),
		ProjectionSleep: projectionSleep,
	}
}

func DebugRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

func AggregateFromWriteModel(ctx context.Context, wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModelCtx(ctx, wm, AggregateType, AggregateVersion)
}
