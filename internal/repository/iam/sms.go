package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	smsConfigPrefix                      = "sms.config"
	smsConfigTwilioPrefix                = "twilio."
	SMSConfigTwilioAddedEventType        = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "added"
	SMSConfigTwilioChangedEventType      = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "changed"
	SMSConfigTwilioTokenChangedEventType = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "token.changed"
	SMSConfigActivatedEventType          = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "activated"
	SMSConfigDeactivatedEventType        = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "deactivated"
	SMSConfigRemovedEventType            = iamEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "removed"
)

type SMSConfigTwilioAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID           string              `json:"id,omitempty"`
	SID          string              `json:"sid,omitempty"`
	Token        *crypto.CryptoValue `json:"token,omitempty"`
	SenderNumber string              `json:"senderNumber,omitempty"`
}

func NewSMSConfigTwilioAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	sid,
	senderNumber string,
	token *crypto.CryptoValue,
) *SMSConfigTwilioAddedEvent {
	return &SMSConfigTwilioAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioAddedEventType,
		),
		ID:           id,
		SID:          sid,
		Token:        token,
		SenderNumber: senderNumber,
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

	ID           string  `json:"id,omitempty"`
	SID          *string `json:"sid,omitempty"`
	SenderNumber *string `json:"senderNumber,omitempty"`
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

func ChangeSMSConfigTwilioSenderNumber(senderNumber string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.SenderNumber = &senderNumber
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

type SMSConfigTwilioTokenChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID    string              `json:"id,omitempty"`
	Token *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMSConfigTokenChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	token *crypto.CryptoValue,
) *SMSConfigTwilioTokenChangedEvent {
	return &SMSConfigTwilioTokenChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioTokenChangedEventType,
		),
		ID:    id,
		Token: token,
	}
}

func (e *SMSConfigTwilioTokenChangedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigTwilioTokenChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigTwilioTokenChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smtpConfigTokenChagned := &SMSConfigTwilioTokenChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smtpConfigTokenChagned)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-fi9Wf", "unable to unmarshal sms config token changed")
	}

	return smtpConfigTokenChagned, nil
}

type SMSConfigActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigTwilioActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigActivatedEvent {
	return &SMSConfigActivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigActivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigActivatedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigActivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigActivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigActivated := &SMSConfigActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigActivated)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio activated changed")
	}

	return smsConfigActivated, nil
}

type SMSConfigDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigDeactivatedEvent {
	return &SMSConfigDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigDeactivatedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigDeactivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigDeactivated := &SMSConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigDeactivated)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio deactivated changed")
	}

	return smsConfigDeactivated, nil
}

type SMSConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigRemovedEvent {
	return &SMSConfigRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigRemovedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigRemovedEvent) Data() interface{} {
	return e
}

func (e *SMSConfigRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMSConfigRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smsConfigRemoved := &SMSConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smsConfigRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal sms config removed")
	}

	return smsConfigRemoved, nil
}
