package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	smtpConfigPrefix                   = "smtp.config."
	SMTPConfigAddedEventType           = instanceEventTypePrefix + smtpConfigPrefix + "added"
	SMTPConfigChangedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "changed"
	SMTPConfigPasswordChangedEventType = instanceEventTypePrefix + smtpConfigPrefix + "password.changed"
	SMTPConfigRemovedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "removed"
)

type SMTPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SenderAddress  string              `json:"senderAddress,omitempty"`
	SenderName     string              `json:"senderName,omitempty"`
	ReplyToAddress string              `json:"replyToAddress,omitempty"`
	TLS            bool                `json:"tls,omitempty"`
	Host           string              `json:"host,omitempty"`
	User           string              `json:"user,omitempty"`
	Password       *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMTPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tls bool,
	senderAddress,
	senderName,
	replyToAddress,
	host,
	user string,
	password *crypto.CryptoValue,
) *SMTPConfigAddedEvent {
	return &SMTPConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigAddedEventType,
		),
		TLS:            tls,
		SenderAddress:  senderAddress,
		SenderName:     senderName,
		ReplyToAddress: replyToAddress,
		Host:           host,
		User:           user,
		Password:       password,
	}
}

func (e *SMTPConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMTPConfigAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smtpConfigAdded := &SMTPConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smtpConfigAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-39fks", "unable to unmarshal smtp config added")
	}

	return smtpConfigAdded, nil
}

type SMTPConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FromAddress    *string `json:"senderAddress,omitempty"`
	FromName       *string `json:"senderName,omitempty"`
	ReplyToAddress *string `json:"replyToAddress,omitempty"`
	TLS            *bool   `json:"tls,omitempty"`
	Host           *string `json:"host,omitempty"`
	User           *string `json:"user,omitempty"`
}

func (e *SMTPConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSMTPConfigChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []SMTPConfigChanges,
) (*SMTPConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-o0pWf", "Errors.NoChangesFound")
	}
	changeEvent := &SMTPConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SMTPConfigChanges func(event *SMTPConfigChangedEvent)

func ChangeSMTPConfigTLS(tls bool) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.TLS = &tls
	}
}

func ChangeSMTPConfigFromAddress(senderAddress string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.FromAddress = &senderAddress
	}
}

func ChangeSMTPConfigFromName(senderName string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.FromName = &senderName
	}
}

func ChangeSMTPConfigReplyToAddress(replyToAddress string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.ReplyToAddress = &replyToAddress
	}
}

func ChangeSMTPConfigSMTPHost(smtpHost string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.Host = &smtpHost
	}
}

func ChangeSMTPConfigSMTPUser(smtpUser string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.User = &smtpUser
	}
}

func SMTPConfigChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SMTPConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-m09oo", "unable to unmarshal smtp changed")
	}

	return e, nil
}

type SMTPConfigPasswordChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Password *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMTPConfigPasswordChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	password *crypto.CryptoValue,
) *SMTPConfigPasswordChangedEvent {
	return &SMTPConfigPasswordChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigPasswordChangedEventType,
		),
		Password: password,
	}
}

func (e *SMTPConfigPasswordChangedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigPasswordChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMTPConfigPasswordChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smtpConfigPasswordChagned := &SMTPConfigPasswordChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smtpConfigPasswordChagned)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal smtp config password changed")
	}

	return smtpConfigPasswordChagned, nil
}

type SMTPConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewSMTPConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *SMTPConfigRemovedEvent {
	return &SMTPConfigRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigRemovedEventType,
		),
	}
}

func (e *SMTPConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func SMTPConfigRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	smtpConfigRemoved := &SMTPConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(smtpConfigRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-DVw1s", "unable to unmarshal smtp config removed")
	}

	return smtpConfigRemoved, nil
}
