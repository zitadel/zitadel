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

func (e *SMSConfigTwilioAddedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigAdded := &SMSConfigTwilioAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
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
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-smn8e", "Errors.NoChangesFound")
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

func (e *SMSConfigTwilioChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigChanged := &SMSConfigTwilioChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
	}

	return smsConfigChanged, nil
}

type SMSConfigTwilioTokenChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

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
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigTwilioTokenChangedEventType,
		),
		ID:    id,
		Token: token,
	}
}

func (e *SMSConfigTwilioTokenChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioTokenChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioTokenChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smtpConfigTokenChagned := &SMSConfigTwilioTokenChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smtpConfigTokenChagned)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-fi9Wf", "unable to unmarshal sms config token changed")
	}

	return smtpConfigTokenChagned, nil
}

type SMSConfigHTTPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID       string `json:"id,omitempty"`
	Endpoint string `json:"endpoint,omitempty"`
}

func NewSMSConfigHTTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id,
	endpoint string,
) *SMSConfigHTTPAddedEvent {
	return &SMSConfigHTTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMSConfigHTTPAddedEventType,
		),
		ID:       id,
		Endpoint: endpoint,
	}
}

func (e *SMSConfigHTTPAddedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigHTTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigHTTPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigAdded := &SMSConfigHTTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
	}

	return smsConfigAdded, nil
}

type SMSConfigHTTPChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID       string  `json:"id,omitempty"`
	Endpoint *string `json:"endpoint,omitempty"`
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
		BaseEvent: *eventstore.NewBaseEventForPush(
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

func ChangeSMSConfigHTTPEndpoint(endpoint string) func(event *SMSConfigHTTPChangedEvent) {
	return func(e *SMSConfigHTTPChangedEvent) {
		e.Endpoint = &endpoint
	}
}

func (e *SMSConfigHTTPChangedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigHTTPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigHTTPChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigChanged := &SMSConfigHTTPChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigChanged)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-smwiR", "unable to unmarshal sms config twilio added")
	}

	return smsConfigChanged, nil
}

type SMSConfigTwilioActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioActivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioActivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigActivated := &SMSConfigTwilioActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigActivated)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio activated changed")
	}

	return smsConfigActivated, nil
}

type SMSConfigActivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMSConfigActivatedEvent(
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

func (e *SMSConfigActivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigActivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigActivated := &SMSConfigActivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigActivated)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio activated changed")
	}

	return smsConfigActivated, nil
}

type SMSConfigTwilioDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigDeactivated := &SMSConfigTwilioDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigDeactivated)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio deactivated changed")
	}

	return smsConfigDeactivated, nil
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

func (e *SMSConfigDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigDeactivated := &SMSConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigDeactivated)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-dn92f", "unable to unmarshal sms config twilio deactivated changed")
	}

	return smsConfigDeactivated, nil
}

type SMSConfigTwilioRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func (e *SMSConfigTwilioRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigTwilioRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigTwilioRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigRemoved := &SMSConfigTwilioRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal sms config removed")
	}

	return smsConfigRemoved, nil
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

func (e *SMSConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMSConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMSConfigRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smsConfigRemoved := &SMSConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smsConfigRemoved)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal sms config removed")
	}

	return smsConfigRemoved, nil
}
