package instance

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	smtpConfigPrefix                   = "smtp.config"
	SMTPConfigAddedEventType           = instanceEventTypePrefix + smtpConfigPrefix + "added"
	SMTPConfigChangedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "changed"
	SMTPConfigPasswordChangedEventType = instanceEventTypePrefix + smtpConfigPrefix + "password.changed"
)

type SMTPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SenderAddress string              `json:"senderAddress,omitempty"`
	SenderName    string              `json:"senderName,omitempty"`
	TLS           bool                `json:"tls,omitempty"`
	Host          string              `json:"host,omitempty"`
	User          string              `json:"user,omitempty"`
	Password      *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMTPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tls bool,
	senderAddress,
	senderName,
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
		TLS:           tls,
		SenderAddress: senderAddress,
		SenderName:    senderName,
		Host:          host,
		User:          user,
		Password:      password,
	}
}

func (e *SMTPConfigAddedEvent) Data() interface{} {
	return e
}

func (e *SMTPConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMTPConfigAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smtpConfigAdded := &SMTPConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smtpConfigAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-39fks", "unable to unmarshal smtp config added")
	}

	return smtpConfigAdded, nil
}

type SMTPConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FromAddress *string `json:"senderAddress,omitempty"`
	FromName    *string `json:"senderName,omitempty"`
	TLS         *bool   `json:"tls,omitempty"`
	Host        *string `json:"host,omitempty"`
	User        *string `json:"user,omitempty"`
}

func (e *SMTPConfigChangedEvent) Data() interface{} {
	return e
}

func (e *SMTPConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func SMTPConfigChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &SMTPConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
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

func (e *SMTPConfigPasswordChangedEvent) Data() interface{} {
	return e
}

func (e *SMTPConfigPasswordChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMTPConfigPasswordChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smtpConfigPasswordChagned := &SMTPConfigPasswordChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smtpConfigPasswordChagned)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal smtp config password changed")
	}

	return smtpConfigPasswordChagned, nil
}
