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
	emailWebhookConfigPrefix                   = "emailwebhook.config."
	EmailWebhookConfigAddedEventType           = instanceEventTypePrefix + emailWebhookConfigPrefix + "added"
	EmailWebhookConfigChangedEventType         = instanceEventTypePrefix + emailWebhookConfigPrefix + "changed"
	EmailWebhookConfigPasswordChangedEventType = instanceEventTypePrefix + emailWebhookConfigPrefix + "password.changed"
	EmailWebhookConfigRemovedEventType         = instanceEventTypePrefix + emailWebhookConfigPrefix + "removed"
)

type EmailWebhookConfigAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SenderAddress  string              `json:"senderAddress,omitempty"`
	SenderName     string              `json:"senderName,omitempty"`
	ReplyToAddress string              `json:"replyToAddress,omitempty"`
	TLS            bool                `json:"tls,omitempty"`
	Host           string              `json:"host,omitempty"`
	User           string              `json:"user,omitempty"`
	Password       *crypto.CryptoValue `json:"password,omitempty"`
}

func NewEmailWebhookConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tls bool,
	senderAddress,
	senderName,
	replyToAddress,
	host,
	user string,
	password *crypto.CryptoValue,
) *EmailWebhookConfigAddedEvent {
	return &EmailWebhookConfigAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailWebhookConfigAddedEventType,
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

func (e *EmailWebhookConfigAddedEvent) Data() interface{} {
	return e
}

func (e *EmailWebhookConfigAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func EmailWebhookConfigAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	emailWebhookConfigAdded := &EmailWebhookConfigAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailWebhookConfigAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-39fks", "unable to unmarshal emailWebhook config added")
	}

	return emailWebhookConfigAdded, nil
}

type EmailWebhookConfigChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FromAddress    *string `json:"senderAddress,omitempty"`
	FromName       *string `json:"senderName,omitempty"`
	ReplyToAddress *string `json:"replyToAddress,omitempty"`
	TLS            *bool   `json:"tls,omitempty"`
	Host           *string `json:"host,omitempty"`
	User           *string `json:"user,omitempty"`
}

func (e *EmailWebhookConfigChangedEvent) Data() interface{} {
	return e
}

func (e *EmailWebhookConfigChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewEmailWebhookConfigChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []EmailWebhookConfigChanges,
) (*EmailWebhookConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "IAM-o0pWf", "Errors.NoChangesFound")
	}
	changeEvent := &EmailWebhookConfigChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailWebhookConfigChangedEventType,
		),
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type EmailWebhookConfigChanges func(event *EmailWebhookConfigChangedEvent)

func ChangeEmailWebhookConfigTLS(tls bool) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.TLS = &tls
	}
}

func ChangeEmailWebhookConfigFromAddress(senderAddress string) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.FromAddress = &senderAddress
	}
}

func ChangeEmailWebhookConfigFromName(senderName string) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.FromName = &senderName
	}
}

func ChangeEmailWebhookConfigReplyToAddress(replyToAddress string) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.ReplyToAddress = &replyToAddress
	}
}

func ChangeEmailWebhookConfigEmailWebhookHost(emailWebhookHost string) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.Host = &emailWebhookHost
	}
}

func ChangeEmailWebhookConfigEmailWebhookUser(emailWebhookUser string) func(event *EmailWebhookConfigChangedEvent) {
	return func(e *EmailWebhookConfigChangedEvent) {
		e.User = &emailWebhookUser
	}
}

func EmailWebhookConfigChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &EmailWebhookConfigChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-m09oo", "unable to unmarshal emailWebhook changed")
	}

	return e, nil
}

type EmailWebhookConfigPasswordChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Password *crypto.CryptoValue `json:"password,omitempty"`
}

func NewEmailWebhookConfigPasswordChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	password *crypto.CryptoValue,
) *EmailWebhookConfigPasswordChangedEvent {
	return &EmailWebhookConfigPasswordChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailWebhookConfigPasswordChangedEventType,
		),
		Password: password,
	}
}

func (e *EmailWebhookConfigPasswordChangedEvent) Data() interface{} {
	return e
}

func (e *EmailWebhookConfigPasswordChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func EmailWebhookConfigPasswordChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	emailWebhookConfigPasswordChagned := &EmailWebhookConfigPasswordChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailWebhookConfigPasswordChagned)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-99iNF", "unable to unmarshal emailWebhook config password changed")
	}

	return emailWebhookConfigPasswordChagned, nil
}

type EmailWebhookConfigRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewEmailWebhookConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *EmailWebhookConfigRemovedEvent {
	return &EmailWebhookConfigRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailWebhookConfigRemovedEventType,
		),
	}
}

func (e *EmailWebhookConfigRemovedEvent) Data() interface{} {
	return e
}

func (e *EmailWebhookConfigRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func EmailWebhookConfigRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	emailWebhookConfigRemoved := &EmailWebhookConfigRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailWebhookConfigRemoved)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-DVw1s", "unable to unmarshal emailWebhook config removed")
	}

	return emailWebhookConfigRemoved, nil
}
