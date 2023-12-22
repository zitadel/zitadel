package policy

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	NotificationPolicyAddedEventType   = "policy.notification.added"
	NotificationPolicyChangedEventType = "policy.notification.changed"
	NotificationPolicyRemovedEventType = "policy.notification.removed"
)

type NotificationPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PasswordChange bool `json:"passwordChange,omitempty"`
}

func (e *NotificationPolicyAddedEvent) Payload() interface{} {
	return e
}

func (e *NotificationPolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotificationPolicyAddedEvent(
	base *eventstore.BaseEvent,
	passwordChange bool,
) *NotificationPolicyAddedEvent {
	return &NotificationPolicyAddedEvent{
		BaseEvent:      *base,
		PasswordChange: passwordChange,
	}
}

func NotificationPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &NotificationPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-0sp2nios", "unable to unmarshal policy")
	}

	return e, nil
}

type NotificationPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PasswordChange *bool `json:"passwordChange,omitempty"`
}

func (e *NotificationPolicyChangedEvent) Payload() interface{} {
	return e
}

func (e *NotificationPolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotificationPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []NotificationPolicyChanges,
) (*NotificationPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-09sp2m", "Errors.NoChangesFound")
	}
	changeEvent := &NotificationPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type NotificationPolicyChanges func(*NotificationPolicyChangedEvent)

func ChangePasswordChange(passwordChange bool) func(*NotificationPolicyChangedEvent) {
	return func(e *NotificationPolicyChangedEvent) {
		e.PasswordChange = &passwordChange
	}
}

func NotificationPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &NotificationPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-09s2oss", "unable to unmarshal policy")
	}

	return e, nil
}

type NotificationPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *NotificationPolicyRemovedEvent) Payload() interface{} {
	return nil
}

func (e *NotificationPolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewNotificationPolicyRemovedEvent(base *eventstore.BaseEvent) *NotificationPolicyRemovedEvent {
	return &NotificationPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func NotificationPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &NotificationPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
