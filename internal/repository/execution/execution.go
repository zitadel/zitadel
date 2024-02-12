package execution

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	eventTypePrefix  eventstore.EventType = "execution."
	AddedEventType                        = eventTypePrefix + "added"
	ChangedEventType                      = eventTypePrefix + "changed"
	RemovedEventType                      = eventTypePrefix + "removed"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Name             string               `json:"name"`
	ExecutionType    domain.ExecutionType `json:"executionType"`
	URL              *url.URL             `json:"url"`
	Timeout          time.Duration        `json:"timeout"`
	Async            bool                 `json:"async"`
	InterruptOnError bool                 `json:"interruptOnError"`
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
	executionType domain.ExecutionType,
	url *url.URL,
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

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExecutionType    *domain.ExecutionType `json:"executionType,omitempty"`
	URL              *url.URL              `json:"url,omitempty"`
	Timeout          *time.Duration        `json:"timeout,omitempty"`
	Async            *bool                 `json:"async,omitempty"`
	InterruptOnError *bool                 `json:"interruptOnError,omitempty"`
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []Changes,
) (*ChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "EXEC-n1yzrtkyb1", "Errors.NoChangesFound")
	}
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
	return changeEvent, nil
}

type Changes func(event *ChangedEvent)

func ChangeType(executionType domain.ExecutionType) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.ExecutionType = &executionType
	}
}

func ChangeURL(url *url.URL) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.URL = url
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
		BaseEvent: *eventstore.BaseEventFromRepo(event),
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
