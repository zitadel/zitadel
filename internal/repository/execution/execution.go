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
	Targets       []string             `json:"targets"`
	Includes      []string             `json:"includes"`
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
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, SetEventType,
		),
		ExecutionType: executionType,
		Targets:       targets,
		Includes:      includes,
	}
}


type RemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ExecutionType domain.ExecutionType `json:"executionType"`
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

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, executionType domain.ExecutionType) *RemovedEvent {
	return &RemovedEvent{
		eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType),
		executionType,
	}
}

