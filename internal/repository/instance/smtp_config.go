package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	smtpConfigPrefix                   = "smtp.config."
	httpConfigPrefix                   = "http."
	SMTPConfigAddedEventType           = instanceEventTypePrefix + smtpConfigPrefix + "added"
	SMTPConfigChangedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "changed"
	SMTPConfigPasswordChangedEventType = instanceEventTypePrefix + smtpConfigPrefix + "password.changed"
	SMTPConfigHTTPAddedEventType       = instanceEventTypePrefix + smtpConfigPrefix + httpConfigPrefix + "added"
	SMTPConfigHTTPChangedEventType     = instanceEventTypePrefix + smtpConfigPrefix + httpConfigPrefix + "changed"
	SMTPConfigRemovedEventType         = instanceEventTypePrefix + smtpConfigPrefix + "removed"
	SMTPConfigActivatedEventType       = instanceEventTypePrefix + smtpConfigPrefix + "activated"
	SMTPConfigDeactivatedEventType     = instanceEventTypePrefix + smtpConfigPrefix + "deactivated"
)

type SMTPConfigAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID             string              `json:"id,omitempty"`
	Description    string              `json:"description,omitempty"`
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
	id, description string,
	tls bool,
	senderAddress,
	senderName,
	replyToAddress,
	host,
	user string,
	password *crypto.CryptoValue,
) *SMTPConfigAddedEvent {
	return &SMTPConfigAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigAddedEventType,
		),
		ID:             id,
		Description:    description,
		TLS:            tls,
		SenderAddress:  senderAddress,
		SenderName:     senderName,
		ReplyToAddress: replyToAddress,
		Host:           host,
		User:           user,
		Password:       password,
	}
}
func (e *SMTPConfigAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMTPConfigChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string              `json:"id,omitempty"`
	Description           *string             `json:"description,omitempty"`
	FromAddress           *string             `json:"senderAddress,omitempty"`
	FromName              *string             `json:"senderName,omitempty"`
	ReplyToAddress        *string             `json:"replyToAddress,omitempty"`
	TLS                   *bool               `json:"tls,omitempty"`
	Host                  *string             `json:"host,omitempty"`
	User                  *string             `json:"user,omitempty"`
	Password              *crypto.CryptoValue `json:"password,omitempty"`
}

func (e *SMTPConfigChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
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
	id string,
	changes []SMTPConfigChanges,
) (*SMTPConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-o0pWf", "Errors.NoChangesFound")
	}
	changeEvent := &SMTPConfigChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SMTPConfigChanges func(event *SMTPConfigChangedEvent)

func ChangeSMTPConfigID(id string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.ID = id
	}
}

func ChangeSMTPConfigDescription(description string) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.Description = &description
	}
}

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

func ChangeSMTPConfigSMTPPassword(password *crypto.CryptoValue) func(event *SMTPConfigChangedEvent) {
	return func(e *SMTPConfigChangedEvent) {
		e.Password = password
	}
}

type SMTPConfigPasswordChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string              `json:"id,omitempty"`
	Password              *crypto.CryptoValue `json:"password,omitempty"`
}

func NewSMTPConfigPasswordChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	password *crypto.CryptoValue,
) *SMTPConfigPasswordChangedEvent {
	return &SMTPConfigPasswordChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigPasswordChangedEventType,
		),
		Password: password,
	}
}

func (e *SMTPConfigPasswordChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigPasswordChangedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigPasswordChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMTPConfigHTTPAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID          string              `json:"id,omitempty"`
	Description string              `json:"description,omitempty"`
	Endpoint    string              `json:"endpoint,omitempty"`
	SigningKey  *crypto.CryptoValue `json:"signingKey"`
}

func NewSMTPConfigHTTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id, description string,
	endpoint string,
	signingKey *crypto.CryptoValue,
) *SMTPConfigHTTPAddedEvent {
	return &SMTPConfigHTTPAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigHTTPAddedEventType,
		),
		ID:          id,
		Description: description,
		Endpoint:    endpoint,
		SigningKey:  signingKey,
	}
}

func (e *SMTPConfigHTTPAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigHTTPAddedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigHTTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMTPConfigHTTPChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string              `json:"id,omitempty"`
	Description           *string             `json:"description,omitempty"`
	Endpoint              *string             `json:"endpoint,omitempty"`
	SigningKey            *crypto.CryptoValue `json:"signingKey,omitempty"`
}

func (e *SMTPConfigHTTPChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigHTTPChangedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigHTTPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSMTPConfigHTTPChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []SMTPConfigHTTPChanges,
) (*SMTPConfigHTTPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "IAM-o0pWf", "Errors.NoChangesFound")
	}
	changeEvent := &SMTPConfigHTTPChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigHTTPChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type SMTPConfigHTTPChanges func(event *SMTPConfigHTTPChangedEvent)

func ChangeSMTPConfigHTTPID(id string) func(event *SMTPConfigHTTPChangedEvent) {
	return func(e *SMTPConfigHTTPChangedEvent) {
		e.ID = id
	}
}

func ChangeSMTPConfigHTTPDescription(description string) func(event *SMTPConfigHTTPChangedEvent) {
	return func(e *SMTPConfigHTTPChangedEvent) {
		e.Description = &description
	}
}

func ChangeSMTPConfigHTTPEndpoint(endpoint string) func(event *SMTPConfigHTTPChangedEvent) {
	return func(e *SMTPConfigHTTPChangedEvent) {
		e.Endpoint = &endpoint
	}
}

func ChangeSMTPConfigHTTPSigningKey(signingKey *crypto.CryptoValue) func(event *SMTPConfigHTTPChangedEvent) {
	return func(e *SMTPConfigHTTPChangedEvent) {
		e.SigningKey = signingKey
	}
}

type SMTPConfigActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMTPConfigActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMTPConfigActivatedEvent {
	return &SMTPConfigActivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigActivatedEventType,
		),
		ID: id,
	}
}

func (e *SMTPConfigActivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigActivatedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMTPConfigDeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMTPConfigDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMTPConfigDeactivatedEvent {
	return &SMTPConfigDeactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *SMTPConfigDeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *SMTPConfigDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type SMTPConfigRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewSMTPConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *SMTPConfigRemovedEvent {
	return &SMTPConfigRemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SMTPConfigRemovedEventType,
		),
		ID: id,
	}
}

func (e *SMTPConfigRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}
func (e *SMTPConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *SMTPConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
