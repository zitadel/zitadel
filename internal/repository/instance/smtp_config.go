package instance

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	smtpConfigPrefix                   = "smtp.config."
	SMTPConfigAddedEventType           = instanceEventTypePrefix + smtpConfigPrefix + "added"
	SMTPConfigChangedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "changed"
	SMTPConfigPasswordChangedEventType = instanceEventTypePrefix + smtpConfigPrefix + "password.changed"
	SMTPConfigRemovedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "removed"
	SMTPConfigDeactivatedEventType     = instanceEventTypePrefix + smtpConfigPrefix + "deactivated"
)

type SMTPConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID             string              `json:"id,omitempty"`
	SenderAddress  string              `json:"senderAddress,omitempty"`
	SenderName     string              `json:"senderName,omitempty"`
	ReplyToAddress string              `json:"replyToAddress,omitempty"`
	TLS            bool                `json:"tls,omitempty"`
	Host           string              `json:"host,omitempty"`
	User           string              `json:"user,omitempty"`
	Password       *crypto.CryptoValue `json:"password,omitempty"`
	IsActive       bool                `json:"isActive,omitempty"`
	ProviderType   uint32              `json:"providerType,omitempty"`
}

func NewSMTPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	tls bool,
	senderAddress,
	senderName,
	replyToAddress,
	host,
	user string,
	password *crypto.CryptoValue,
	isActive bool,
	providerType uint32,
) *SMTPConfigAddedEvent {
	return &SMTPConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigAddedEventType,
		),
		ID:             id,
		TLS:            tls,
		SenderAddress:  senderAddress,
		SenderName:     senderName,
		ReplyToAddress: replyToAddress,
		Host:           host,
		User:           user,
		Password:       password,
		IsActive:       isActive,
		ProviderType:   providerType,
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

	FromAddress    *string             `json:"senderAddress,omitempty"`
	FromName       *string             `json:"senderName,omitempty"`
	ReplyToAddress *string             `json:"replyToAddress,omitempty"`
	TLS            *bool               `json:"tls,omitempty"`
	Host           *string             `json:"host,omitempty"`
	User           *string             `json:"user,omitempty"`
	Password       *crypto.CryptoValue `json:"password,omitempty"`
	IsActive       *bool               `json:"isActive,omitempty"`
	ProviderType   *uint32             `json:"providerType,omitempty"`
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

func ChangeSMTPConfigIsActive(isActive bool) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.IsActive = &isActive
	}
}

func ChangeSMTPConfigProviderType(providerType uint32) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.ProviderType = &providerType
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

type SMTPConfigDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMTPConfigDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMTPConfigDeactivatedEvent {
	return &SMTPConfigDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *SMTPConfigDeactivatedEvent) Data() interface{} {
	return e
}

func (e *SMTPConfigDeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMTPConfigDeactivatedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smtpConfigDeactivated := &SMTPConfigDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smtpConfigDeactivated)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-KPr5t", "unable to unmarshal smtp config removed")
	}

	return smtpConfigDeactivated, nil
}

type SMTPConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	ID                   string `json:"id,omitempty"`
}

func NewSMTPConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMTPConfigRemovedEvent {
	return &SMTPConfigRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigRemovedEventType,
		),
		ID: id,
	}
}

func (e *SMTPConfigRemovedEvent) Data() interface{} {
	return e
}

func (e *SMTPConfigRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func SMTPConfigRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	smtpConfigRemoved := &SMTPConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, smtpConfigRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-DVw1s", "unable to unmarshal smtp config removed")
	}

	return smtpConfigRemoved, nil
}
