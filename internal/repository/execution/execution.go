package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	eventTypePrefix  eventstore.EventType = "execution."
	SetEventType                          = eventTypePrefix + "set"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type SetEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ExecutionType domain.ExecutionType `json:"executionType"`
	Targets       []string             `json:"target"`
	Includes      []string             `json:"include"`
}

func (e *SetEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEvent) Payload() any {
	return e
}

func (e *SetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	executionType domain.ExecutionType,
	targets []string,
	includes []string,
) *SetEvent {
	return &SetEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, SetEventType,
		),
		executionType,
		targets, includes,
	}
}

func SetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SetEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-r8e2e6hawz", "unable to unmarshal execution set")
	}

	return added, nil
}

type RemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedEvent) Payload() any {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *RemovedEvent {
	return &RemovedEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType)}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-rsg1cnt5am", "unable to unmarshal execution removed")
	}

	return removed, nil
}
