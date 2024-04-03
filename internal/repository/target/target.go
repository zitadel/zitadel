package target

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	eventTypePrefix  eventstore.EventType = "target."
	AddedEventType                        = eventTypePrefix + "added"
	ChangedEventType                      = eventTypePrefix + "changed"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name             string            `json:"name"`
	TargetType       domain.TargetType `json:"targetType"`
	URL              string            `json:"url"`
	Timeout          time.Duration     `json:"timeout"`
	Async            bool              `json:"async"`
	InterruptOnError bool              `json:"interruptOnError"`
}

func (e *AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *AddedEvent) Payload() any {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddUniqueConstraint(e.Name)}
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	targetType domain.TargetType,
	url string,
	timeout time.Duration,
	async bool,
	interruptOnError bool,
) *AddedEvent {
	return &AddedEvent{
		*eventstore.NewBaseEventForPush(
			ctx, aggregate, AddedEventType,
		),
		name, targetType, url, timeout, async, interruptOnError}
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name             *string            `json:"name,omitempty"`
	TargetType       *domain.TargetType `json:"targetType,omitempty"`
	URL              *string            `json:"url,omitempty"`
	Timeout          *time.Duration     `json:"timeout,omitempty"`
	Async            *bool              `json:"async,omitempty"`
	InterruptOnError *bool              `json:"interruptOnError,omitempty"`

	oldName string
}

func (e *ChangedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldName == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveUniqueConstraint(e.oldName),
		NewAddUniqueConstraint(*e.Name),
	}
}

func NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []Changes,
) *ChangedEvent {
	changeEvent := &ChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent
}

type Changes func(event *ChangedEvent)

func ChangeName(oldName, name string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Name = &name
		e.oldName = oldName
	}
}

func ChangeTargetType(targetType domain.TargetType) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.TargetType = &targetType
	}
}

func ChangeURL(url string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.URL = &url
	}
}

func ChangeTimeout(timeout time.Duration) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Timeout = &timeout
	}
}

func ChangeAsync(async bool) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Async = &async
	}
}

func ChangeInterruptOnError(interruptOnError bool) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.InterruptOnError = &interruptOnError
	}
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	name string
}

func (e *RemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func (e *RemovedEvent) Payload() any {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUniqueConstraint(e.name)}
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *RemovedEvent {
	return &RemovedEvent{*eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType), name}
}
