package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	smsConfigPrefix                      = "sms.config"
	smsConfigTwilioPrefix                = "twilio."
	smsConfigHTTPPrefix                  = "http."
	SMSConfigTwilioAddedEventType        = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "added"
	SMSConfigTwilioChangedEventType      = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "changed"
	SMSConfigHTTPAddedEventType          = instanceEventTypePrefix + smsConfigPrefix + smsConfigHTTPPrefix + "added"
	SMSConfigHTTPChangedEventType        = instanceEventTypePrefix + smsConfigPrefix + smsConfigHTTPPrefix + "changed"
	SMSConfigTwilioTokenChangedEventType = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "token.changed"
	SMSConfigTwilioActivatedEventType    = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "activated"
	SMSConfigTwilioDeactivatedEventType  = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "deactivated"
	SMSConfigTwilioRemovedEventType      = instanceEventTypePrefix + smsConfigPrefix + smsConfigTwilioPrefix + "removed"
	SMSConfigActivatedEventType          = instanceEventTypePrefix + smsConfigPrefix + "activated"
	SMSConfigDeactivatedEventType        = instanceEventTypePrefix + smsConfigPrefix + "deactivated"
	SMSConfigRemovedEventType            = instanceEventTypePrefix + smsConfigPrefix + "removed"
)

type SMSConfigTwilioAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID               string              `json:"id,omitempty"`
	Description      string              `json:"description,omitempty"`
	SID              string              `json:"sid,omitempty"`
	Token            *crypto.CryptoValue `json:"token,omitempty"`
	SenderNumber     string              `json:"senderNumber,omitempty"`
	VerifyServiceSID string              `json:"verifyServiceSid,omitempty"`
}

func NewSMSConfigTwilioAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	description string,
	sid,
	senderNumber string,
	token *crypto.CryptoValue,
	verifyServiceSid string,
) *SMSConfigTwilioAddedEvent {
	return &SMSConfigTwilioAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioAddedEventType,
		),
		ID:               id,
		Description:      description,
		SID:              sid,
		Token:            token,
		SenderNumber:     senderNumber,
		VerifyServiceSID: verifyServiceSid,
	}
}

func (e *SMSConfigTwilioAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioAddedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigTwilioChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID               string  `json:"id,omitempty"`
	Description      *string `json:"description,omitempty"`
	SID              *string `json:"sid,omitempty"`
	SenderNumber     *string `json:"senderNumber,omitempty"`
	VerifyServiceSID *string `json:"verifyServiceSid,omitempty"`
}

func NewSMSConfigTwilioChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []SMSConfigTwilioChanges,
) (*SMSConfigTwilioChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-smn8e", "Errors.NoChangesFound")
	}
	changeEvent := &SMSConfigTwilioChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
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

func ChangeSMSConfigTwilioDescription(description string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.Description = &description
	}
}

func ChangeSMSConfigTwilioSenderNumber(senderNumber string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.SenderNumber = &senderNumber
	}
}

func ChangeSMSConfigTwilioVerifyServiceSID(verifyServiceSID string) func(event *SMSConfigTwilioChangedEvent) {
	return func(e *SMSConfigTwilioChangedEvent) {
		e.VerifyServiceSID = &verifyServiceSID
	}
}

func (e *SMSConfigTwilioChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigTwilioTokenChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID    string              `json:"id,omitempty"`
	Token *crypto.CryptoValue `json:"token,omitempty"`
}

func NewSMSConfigTokenChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	token *crypto.CryptoValue,
) *SMSConfigTwilioTokenChangedEvent {
	return &SMSConfigTwilioTokenChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioTokenChangedEventType,
		),
		ID:    id,
		Token: token,
	}
}

func (e *SMSConfigTwilioTokenChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioTokenChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioTokenChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigHTTPAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID          string              `json:"id,omitempty"`
	Description string              `json:"description,omitempty"`
	Endpoint    string              `json:"endpoint,omitempty"`
	SigningKey  *crypto.CryptoValue `json:"signingKey"`
}

func NewSMSConfigHTTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	description,
	endpoint string,
	signingKey *crypto.CryptoValue,
) *SMSConfigHTTPAddedEvent {
	return &SMSConfigHTTPAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigHTTPAddedEventType,
		),
		ID:          id,
		Description: description,
		Endpoint:    endpoint,
		SigningKey:  signingKey,
	}
}

func (e *SMSConfigHTTPAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigHTTPAddedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigHTTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigHTTPChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID          string              `json:"id,omitempty"`
	Description *string             `json:"description,omitempty"`
	Endpoint    *string             `json:"endpoint,omitempty"`
	SigningKey  *crypto.CryptoValue `json:"signingKey,omitempty"`
}

func NewSMSConfigHTTPChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []SMSConfigHTTPChanges,
) (*SMSConfigHTTPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-smn8e", "Errors.NoChangesFound")
	}
	changeEvent := &SMSConfigHTTPChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigHTTPChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SMSConfigHTTPChanges func(event *SMSConfigHTTPChangedEvent)

func ChangeSMSConfigHTTPDescription(description string) func(event *SMSConfigHTTPChangedEvent) {
	return func(e *SMSConfigHTTPChangedEvent) {
		e.Description = &description
	}
}
func ChangeSMSConfigHTTPEndpoint(endpoint string) func(event *SMSConfigHTTPChangedEvent) {
	return func(e *SMSConfigHTTPChangedEvent) {
		e.Endpoint = &endpoint
	}
}
func ChangeSMSConfigHTTPSigningKey(signingKey *crypto.CryptoValue) func(event *SMSConfigHTTPChangedEvent) {
	return func(e *SMSConfigHTTPChangedEvent) {
		e.SigningKey = signingKey
	}
}

func (e *SMSConfigHTTPChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigHTTPChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigHTTPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigTwilioActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioActivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioActivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMSConfigActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigActivatedEvent {
	return &SMSConfigActivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigActivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigActivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigActivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigTwilioDeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioDeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigDeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMSConfigDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigDeactivatedEvent {
	return &SMSConfigDeactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigDeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigTwilioRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigTwilioRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMSConfigRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMSConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMSConfigRemovedEvent {
	return &SMSConfigRemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigRemovedEventType,
		),
		ID: id,
	}
}

func (e *SMSConfigRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMSConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
