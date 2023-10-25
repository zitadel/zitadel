package settings

import (
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
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

func (e *DebugNotificationProviderAddedEvent) Payload() interface{} {
	return e
}

func (e *DebugNotificationProviderAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func DebugNotificationProviderAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DebugNotificationProviderAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "SET-f93ns", "unable to unmarshal debug notification added")
	}

	return e, nil
}

type DebugNotificationProviderChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Compact *bool `json:"compact,omitempty"`
}

func (e *DebugNotificationProviderChangedEvent) Payload() interface{} {
	return e
}

func (e *DebugNotificationProviderChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func DebugNotificationProviderChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DebugNotificationProviderChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type DebugNotificationProviderRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DebugNotificationProviderRemovedEvent) Payload() interface{} {
	return nil
}

func (e *DebugNotificationProviderRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDebugNotificationProviderRemovedEvent(base *eventstore.BaseEvent) *DebugNotificationProviderRemovedEvent {
	return &DebugNotificationProviderRemovedEvent{
		BaseEvent: *base,
	}
}

func DebugNotificationProviderRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &DebugNotificationProviderRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
