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
	smtpConfigPrefix           = "smtp.config"
	SMTPConfigAddedEventType   = iamEventTypePrefix + smtpConfigPrefix + "added"
	SMTPConfigChangedEventType = iamEventTypePrefix + smtpConfigPrefix + "changed"
)

type SMTPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FromAddress  string              `json:"fromAddress,omitempty"`
	FromName     string              `json:"fromName,omitempty"`
	TLS          bool                `json:"tls,omitempty"`
	SMTPHost     string              `json:"host,omitempty"`
	SMTPUser     string              `json:"user,omitempty"`
	SMTPPassword *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMTPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tls bool,
	fromAddress,
	fromName,
	smtpHost,
	smtpUser string,
	smtpPassword *crypto.CryptoValue,
) *SMTPConfigAddedEvent {
	return &SMTPConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigAddedEventType,
		),
		TLS:          tls,
		FromAddress:  fromAddress,
		FromName:     fromName,
		SMTPHost:     smtpHost,
		SMTPUser:     smtpUser,
		SMTPPassword: smtpPassword,
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

	FromAddress  *string             `json:"fromAddress,omitempty"`
	FromName     *string             `json:"fromName,omitempty"`
	TLS          *bool               `json:"tls,omitempty"`
	SMTPHost     *string             `json:"host,omitempty"`
	SMTPUser     *string             `json:"user,omitempty"`
	SMTPPassword *crypto.CryptoValue `json:"password,omitempty"`
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

func ChangeSMTPConfigFromAddress(fromAddress string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.FromAddress = &fromAddress
	}
}

func ChangeSMTPConfigFromName(fromName string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.FromName = &fromName
	}
}

func ChangeSMTPConfigSMTPHost(smtpHost string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.SMTPHost = &smtpHost
	}
}

func ChangeSMTPConfigSMTPUser(smtpUser string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.SMTPUser = &smtpUser
	}
}

func ChangeSMTPConfigSMTPPassword(smtpPassword *crypto.CryptoValue) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.SMTPPassword = smtpPassword
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
