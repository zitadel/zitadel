package settings

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	DebugNotificationPrefix           = "notification.provider.debug"
	DebugNotificationProviderAdded    = "added"
	DebugNotificationProviderChanged  = "changed"
	DebugNotificationProviderEnabled  = "enabled"
	DebugNotificationProviderDisabled = "disabled"
	DebugNotificationProviderRemoved  = "removed"
)

type DebugNotificationProviderAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Compact bool `json:"compact,omitempty"`
}

func (e *DebugNotificationProviderAddedEvent) Data() interface{} {
	return e
}

func (e *DebugNotificationProviderAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDebugNotificationProviderAddedEvent(
	base *eventstore.BaseEvent,
	compact bool,
) *DebugNotificationProviderAddedEvent {
	return &DebugNotificationProviderAddedEvent{
		BaseEvent: *base,
		Compact:   compact,
	}
}

func DebugNotificationProviderAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DebugNotificationProviderAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SET-f93ns", "unable to unmarshal debug notification added")
	}

	return e, nil
}

type DebugNotificationProviderChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Compact *bool `json:"compact,omitempty"`
}

func (e *DebugNotificationProviderChangedEvent) Data() interface{} {
	return e
}

func (e *DebugNotificationProviderChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDebugNotificationProviderChangedEvent(
	base *eventstore.BaseEvent,
	changes []DebugNotificationProviderChanges,
) (*DebugNotificationProviderChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "SET-hj90s", "Errors.NoChangesFound")
	}
	changeEvent := &DebugNotificationProviderChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type DebugNotificationProviderChanges func(*DebugNotificationProviderChangedEvent)

func ChangeCompact(compact bool) func(*DebugNotificationProviderChangedEvent) {
	return func(e *DebugNotificationProviderChangedEvent) {
		e.Compact = &compact
	}
}

func DebugNotificationProviderChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DebugNotificationProviderChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type DebugNotificationProviderEnabledEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DebugNotificationProviderEnabledEvent) Data() interface{} {
	return nil
}

func (e *DebugNotificationProviderEnabledEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDebugNotificationProviderEnabledEvent(base *eventstore.BaseEvent) *DebugNotificationProviderEnabledEvent {
	return &DebugNotificationProviderEnabledEvent{
		BaseEvent: *base,
	}
}

func DebugNotificationProviderEnabledEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &DebugNotificationProviderEnabledEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type DebugNotificationProviderDisabledEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DebugNotificationProviderDisabledEvent) Data() interface{} {
	return nil
}

func (e *DebugNotificationProviderDisabledEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDebugNotificationProviderDisabledEvent(base *eventstore.BaseEvent) *DebugNotificationProviderDisabledEvent {
	return &DebugNotificationProviderDisabledEvent{
		BaseEvent: *base,
	}
}

func DebugNotificationProviderDisabledEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &DebugNotificationProviderDisabledEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type DebugNotificationProviderRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DebugNotificationProviderRemovedEvent) Data() interface{} {
	return nil
}

func (e *DebugNotificationProviderRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDebugNotificationProviderRemovedEvent(base *eventstore.BaseEvent) *DebugNotificationProviderRemovedEvent {
	return &DebugNotificationProviderRemovedEvent{
		BaseEvent: *base,
	}
}

func DebugNotificationProviderRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &DebugNotificationProviderRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
