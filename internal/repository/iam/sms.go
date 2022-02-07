package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	smsConfigPrefix                     = "sms.config"
	smsConfigTwilioPrefix               = "twilio."
	SMSConfigTwilioAddedEventType       = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "added"
	SMSConfigTwilioChangedEventType     = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "changed"
	SMSConfigTwilioActivatedEventType   = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "activated"
	SMSConfigTwilioDeactivatedEventType = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "deactivated"
	SMSConfigTwilioRemovedEventType     = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "removed"
)

type SMSConfigTwilioAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID    string `json:"id,omitempty"`
	SID   string `json:"sid,omitempty"`
	Token string `json:"token,omitempty"`
	From  string `json:"from,omitempty"`
}

func NewSMSConfigTwilioAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	sid,
	token,
	from string,
) *SMSConfigTwilioAddedEvent {
	return &SMSConfigTwilioAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioAddedEventType,
		),
		ID:    id,
		SID:   sid,
		Token: token,
		From:  from,
	}
}

func (e *SMSConfigTwilioAddedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigAdded := &SMSConfigTwilioAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
	}

	return smsConfigAdded, nil
}

type SMSConfigTwilioChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID    string  `json:"id,omitempty"`
	SID   *string `json:"sid,omitempty"`
	Token *string `json:"token,omitempty"`
	From  *string `json:"from,omitempty"`
}

func NewSMSConfigTwilioChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []SMSConfigTwilioChanges,
) (*SMSConfigTwilioChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-smn8e", "Errors.NoChangesFound")
	}
	changeEvent := &SMSConfigTwilioChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SMSConfigTwilioChanges func(event *SMSConfigTwilioChangedEvent)

func ChangeSMSConfigTwilioSID(sid string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.SID = &sid
	}
}

func ChangeSMSConfigTwilioToken(token string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.Token = &token
	}
}

func ChangeSMSConfigTwilioFrom(from string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.From = &from
	}
}

func (e *SMSConfigTwilioChangedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigChanged := &SMSConfigTwilioChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigChanged)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
	}

	return smsConfigChanged, nil
}

type SMSConfigTwilioActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigTwilioActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigTwilioActivatedEvent {
	return &SMSConfigTwilioActivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioActivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigTwilioActivatedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioActivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioActivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigActivated := &SMSConfigTwilioActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigActivated)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio activated changed")
	}

	return smsConfigActivated, nil
}

type SMSConfigTwilioDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigTwilioDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigTwilioDeactivatedEvent {
	return &SMSConfigTwilioDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigTwilioDeactivatedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioDeactivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigDeactivated := &SMSConfigTwilioDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigDeactivated)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio deactivated changed")
	}

	return smsConfigDeactivated, nil
}

type SMSConfigTwilioRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigTwilioRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigTwilioRemovedEvent {
	return &SMSConfigTwilioRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioRemovedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigTwilioRemovedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigRemoved := &SMSConfigTwilioRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal sms config password changed")
	}

	return smsConfigRemoved, nil
}
