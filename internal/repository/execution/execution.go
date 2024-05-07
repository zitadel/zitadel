package execution

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix  eventstore.EventType = "execution."
	SetEventType                          = eventTypePrefix + "set"
	SetEventV2Type                        = eventTypePrefix + "v2.set"
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

type SetEventV2 struct {
	*eventstore.BaseEvent `json:"-"`

	Targets []*Target `json:"targets"`
}

func (e *SetEventV2) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *SetEventV2) Payload() any {
	return e
}

func (e *SetEventV2) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type Target struct {
	Type   domain.ExecutionTargetType `json:"type"`
	Target string                     `json:"target"`
}

func NewSetEventV2(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	targets []*Target,
) *SetEventV2 {
	return &SetEventV2{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx, aggregate, SetEventV2Type,
		),
		Targets: targets,
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
