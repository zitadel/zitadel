package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix  eventstore.EventType = "execution."
	SetEventType                          = eventTypePrefix + "set"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type SetEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Targets  []string `json:"targets"`
	Includes []string `json:"includes"`
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
	targets []string,
	includes []string,
) *SetEvent {
	return &SetEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, SetEventType,
		),
		Targets:  targets,
		Includes: includes,
	}
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
	return &RemovedEvent{
		eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType),
	}
}
