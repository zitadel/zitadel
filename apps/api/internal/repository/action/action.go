package action

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueActionNameType = "action_names"
	eventTypePrefix      = eventstore.EventType("action.")
	AddedEventType       = eventTypePrefix + "added"
	ChangedEventType     = eventTypePrefix + "changed"
	DeactivatedEventType = eventTypePrefix + "deactivated"
	ReactivatedEventType = eventTypePrefix + "reactivated"
	RemovedEventType     = eventTypePrefix + "removed"
)

func NewAddActionNameUniqueConstraint(actionName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueActionNameType,
		actionName+":"+resourceOwner,
		"Errors.Action.AlreadyExists")
}

func NewRemoveActionNameUniqueConstraint(actionName, resourceOwner string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueActionNameType,
		actionName+":"+resourceOwner)
}

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name          string        `json:"name"`
	Script        string        `json:"script,omitempty"`
	Timeout       time.Duration `json:"timeout,omitempty"`
	AllowedToFail bool          `json:"allowedToFail"`
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddActionNameUniqueConstraint(e.Name, e.Aggregate().ResourceOwner)}
}

func NewAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name,
	script string,
	timeout time.Duration,
	allowedToFail bool,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedEventType,
		),
		Name:          name,
		Script:        script,
		Timeout:       timeout,
		AllowedToFail: allowedToFail,
	}
}

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ACTION-4n8vs", "unable to unmarshal action added")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Name          *string        `json:"name,omitempty"`
	Script        *string        `json:"script,omitempty"`
	Timeout       *time.Duration `json:"timeout,omitempty"`
	AllowedToFail *bool          `json:"allowedToFail,omitempty"`
	oldName       string
}

func (e *ChangedEvent) Payload() interface{} {
	return e
}

func (e *ChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	if e.oldName == "" {
		return nil
	}
	return []*eventstore.UniqueConstraint{
		NewRemoveActionNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
		NewAddActionNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
	}
}

func NewChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []ActionChanges,
) (*ChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ACTION-dg4t2", "Errors.NoChangesFound")
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

type ActionChanges func(event *ChangedEvent)

func ChangeName(name, oldName string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Name = &name
		e.oldName = oldName
	}
}

func ChangeScript(script string) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Script = &script
	}
}

func ChangeTimeout(timeout time.Duration) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.Timeout = &timeout
	}
}

func ChangeAllowedToFail(allowedToFail bool) func(event *ChangedEvent) {
	return func(e *ChangedEvent) {
		e.AllowedToFail = &allowedToFail
	}
}

func ChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ACTION-4n8vs", "unable to unmarshal action changed")
	}

	return e, nil
}

type DeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *DeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *DeactivatedEvent {
	return &DeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DeactivatedEventType,
		),
	}
}

func DeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &DeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type ReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *ReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *ReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ReactivatedEvent {
	return &ReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			ReactivatedEventType,
		),
	}
}

func ReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &ReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	name string
}

func (e *RemovedEvent) Payload() interface{} {
	return e
}

func (e *RemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveActionNameUniqueConstraint(e.name, e.Aggregate().ResourceOwner)}
}

func NewRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	name string,
) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			RemovedEventType,
		),
		name: name,
	}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
