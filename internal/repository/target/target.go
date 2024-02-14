package target

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	eventTypePrefix  eventstore.EventType = "target."
	AddedEventType                        = eventTypePrefix + "added"
	ChangedEventType                      = eventTypePrefix + "changed"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Name             string            `json:"name"`
	ExecutionType    domain.TargetType `json:"executionType"`
	URL              string            `json:"url"`
	Timeout          time.Duration     `json:"timeout"`
	Async            bool              `json:"async"`
	InterruptOnError bool              `json:"interruptOnError"`
}

func (e *AddedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *AddedEvent) Payload() any {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return NewAddUniqueConstraints(e.Name)
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
	executionType domain.TargetType,
	url string,
	timeout time.Duration,
	async bool,
	interruptOnError bool,
) *AddedEvent {
	return &AddedEvent{
		eventstore.NewBaseEventForPush(
			ctx, aggregate, AddedEventType,
		),
		name, executionType, url, timeout, async, interruptOnError}
}

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-fx8f8yfbn1", "unable to unmarshal execution added")
	}

	return added, nil
}

type ChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Name             *string            `json:"name,omitempty"`
	ExecutionType    *domain.TargetType `json:"executionType,omitempty"`
	URL              *string            `json:"url,omitempty"`
	Timeout          *time.Duration     `json:"timeout,omitempty"`
	Async            *bool              `json:"async,omitempty"`
	InterruptOnError *bool              `json:"interruptOnError,omitempty"`

	oldName string
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldName == "" {
		return nil
	}
	return NewUpdateUniqueConstraints(e.oldName, *e.Name)
}

func NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []Changes,
) *ChangedEvent {
	changeEvent := &ChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
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

func ChangeExecutionType(executionType domain.TargetType) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.ExecutionType = &executionType
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

func ChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	changed := &ChangedEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(changed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-w6402p4ek7", "unable to unmarshal execution changed")
	}

	return changed, nil
}

type RemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	name string
}

func (e *RemovedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = b
}

func (e *RemovedEvent) Payload() any {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return NewRemoveUniqueConstraints(e.name)
}

func NewRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate, name string) *RemovedEvent {
	return &RemovedEvent{eventstore.NewBaseEventForPush(ctx, aggregate, RemovedEventType), name}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	removed := &RemovedEvent{
		BaseEvent: eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(removed)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "EXEC-0kuc12c7bc", "unable to unmarshal execution removed")
	}

	return removed, nil
}
