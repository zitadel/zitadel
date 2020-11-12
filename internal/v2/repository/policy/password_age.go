package policy

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	PasswordAgePolicyAddedEventType   = "policy.password.age.added"
	PasswordAgePolicyChangedEventType = "policy.password.age.changed"
	PasswordAgePolicyRemovedEventType = "policy.password.age.removed"
)

type PasswordAgePolicyAggregate struct {
	eventstore.Aggregate

	ExpireWarnDays uint16
	MaxAgeDays     uint16
}

type PasswordAgePolicyReadModel struct {
	eventstore.ReadModel

	ExpireWarnDays uint16
	MaxAgeDays     uint16
}

func (rm *PasswordAgePolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *PasswordAgePolicyAddedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		case *PasswordAgePolicyChangedEvent:
			rm.ExpireWarnDays = e.ExpireWarnDays
			rm.MaxAgeDays = e.MaxAgeDays
		}
	}
	return rm.ReadModel.Reduce()
}

type PasswordAgePolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays uint16 `json:"expireWarnDays"`
	MaxAgeDays     uint16 `json:"maxAgeDays"`
}

func (e *PasswordAgePolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays uint16,
) *PasswordAgePolicyAddedEvent {

	return &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyAddedEventType,
		),
		ExpireWarnDays: expireWarnDays,
		MaxAgeDays:     maxAgeDays,
	}
}

func PasswordAgePolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-T3mGp", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays uint16 `json:"expireWarnDays,omitempty"`
	MaxAgeDays     uint16 `json:"maxAgeDays,omitempty"`
}

func (e *PasswordAgePolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordAgePolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordAgePolicyAggregate,
) *PasswordAgePolicyChangedEvent {

	e := &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
	}

	if current.ExpireWarnDays != changed.ExpireWarnDays {
		e.ExpireWarnDays = changed.ExpireWarnDays
	}
	if current.MaxAgeDays != changed.MaxAgeDays {
		e.MaxAgeDays = changed.ExpireWarnDays
	}

	return e
}

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-PqaVq", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordAgePolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordAgePolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordAgePolicyRemovedEvent(
	ctx context.Context,
	current,
	changed *PasswordAgePolicyRemovedEvent,
) *PasswordAgePolicyChangedEvent {

	return &PasswordAgePolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyRemovedEventType,
		),
	}
}

func PasswordAgePolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordAgePolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-02878", "unable to unmarshal policy")
	}

	return e, nil
}
