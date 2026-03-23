package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	smtpConfigPrefix                      = "smtp.config."
	httpConfigPrefix                      = "http."
	OrgSMTPConfigAddedEventType           = orgEventTypePrefix + smtpConfigPrefix + "added"
	OrgSMTPConfigChangedEventType         = orgEventTypePrefix + smtpConfigPrefix + "changed"
	OrgSMTPConfigPasswordChangedEventType = orgEventTypePrefix + smtpConfigPrefix + "password.changed"
	OrgSMTPConfigHTTPAddedEventType       = orgEventTypePrefix + smtpConfigPrefix + httpConfigPrefix + "added"
	OrgSMTPConfigHTTPChangedEventType     = orgEventTypePrefix + smtpConfigPrefix + httpConfigPrefix + "changed"
	OrgSMTPConfigRemovedEventType         = orgEventTypePrefix + smtpConfigPrefix + "removed"
	OrgSMTPConfigActivatedEventType       = orgEventTypePrefix + smtpConfigPrefix + "activated"
	OrgSMTPConfigDeactivatedEventType     = orgEventTypePrefix + smtpConfigPrefix + "deactivated"
)

// OrgSMTPConfigAddedEvent is the org-level equivalent of instance.SMTPConfigAddedEvent.
type OrgSMTPConfigAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID             string                `json:"id,omitempty"`
	Description    string                `json:"description,omitempty"`
	SenderAddress  string                `json:"senderAddress,omitempty"`
	SenderName     string                `json:"senderName,omitempty"`
	ReplyToAddress string                `json:"replyToAddress,omitempty"`
	TLS            bool                  `json:"tls,omitempty"`
	Host           string                `json:"host,omitempty"`
	User           string                `json:"user,omitempty"`
	Password       *crypto.CryptoValue   `json:"password,omitempty"`
	PlainAuth      *instance.PlainAuth   `json:"plainAuth,omitempty"`
	XOAuth2Auth    *instance.XOAuth2Auth `json:"xoauth2Auth,omitempty"`
}

func NewOrgSMTPConfigAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id, description string,
	tls bool,
	senderAddress,
	senderName,
	replyToAddress,
	host string,
	user string,
	plainAuth *instance.PlainAuth,
	xoauth2Auth *instance.XOAuth2Auth,
) *OrgSMTPConfigAddedEvent {
	return &OrgSMTPConfigAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigAddedEventType,
		),
		ID:             id,
		Description:    description,
		TLS:            tls,
		SenderAddress:  senderAddress,
		SenderName:     senderName,
		ReplyToAddress: replyToAddress,
		Host:           host,
		User:           user,
		PlainAuth:      plainAuth,
		XOAuth2Auth:    xoauth2Auth,
	}
}

func (e *OrgSMTPConfigAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigAddedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// OrgSMTPConfigChangedEvent is the org-level equivalent of instance.SMTPConfigChangedEvent.
type OrgSMTPConfigChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string                      `json:"id,omitempty"`
	Description           *string                     `json:"description,omitempty"`
	FromAddress           *string                     `json:"senderAddress,omitempty"`
	FromName              *string                     `json:"senderName,omitempty"`
	ReplyToAddress        *string                     `json:"replyToAddress,omitempty"`
	TLS                   *bool                       `json:"tls,omitempty"`
	Host                  *string                     `json:"host,omitempty"`
	User                  *string                     `json:"user,omitempty"`
	Password              *crypto.CryptoValue         `json:"password,omitempty"`
	PlainAuth             instance.PlainAuthChanged   `json:"plainAuth,omitempty"`
	XOAuth2Auth           instance.XOAuth2AuthChanged `json:"xoauth2Auth,omitempty"`
}

func (e *OrgSMTPConfigChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigChangedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type OrgSMTPConfigChanges func(event *OrgSMTPConfigChangedEvent)

func NewOrgSMTPConfigChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []OrgSMTPConfigChanges,
) (*OrgSMTPConfigChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-Wkf8a", "Errors.NoChangesFound")
	}
	changeEvent := &OrgSMTPConfigChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

// OrgSMTPConfigPasswordChangedEvent is the org-level equivalent of instance.SMTPConfigPasswordChangedEvent.
type OrgSMTPConfigPasswordChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string              `json:"id,omitempty"`
	Password              *crypto.CryptoValue `json:"password,omitempty"`
}

func NewOrgSMTPConfigPasswordChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	password *crypto.CryptoValue,
) *OrgSMTPConfigPasswordChangedEvent {
	return &OrgSMTPConfigPasswordChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigPasswordChangedEventType,
		),
		ID:       id,
		Password: password,
	}
}

func (e *OrgSMTPConfigPasswordChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigPasswordChangedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigPasswordChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// OrgSMTPConfigHTTPAddedEvent is the org-level equivalent of instance.SMTPConfigHTTPAddedEvent.
type OrgSMTPConfigHTTPAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	ID          string              `json:"id,omitempty"`
	Description string              `json:"description,omitempty"`
	Endpoint    string              `json:"endpoint,omitempty"`
	SigningKey  *crypto.CryptoValue `json:"signingKey"`
}

func NewOrgSMTPConfigHTTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id, description string,
	endpoint string,
	signingKey *crypto.CryptoValue,
) *OrgSMTPConfigHTTPAddedEvent {
	return &OrgSMTPConfigHTTPAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigHTTPAddedEventType,
		),
		ID:          id,
		Description: description,
		Endpoint:    endpoint,
		SigningKey:  signingKey,
	}
}

func (e *OrgSMTPConfigHTTPAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigHTTPAddedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigHTTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// OrgSMTPConfigHTTPChangedEvent is the org-level equivalent of instance.SMTPConfigHTTPChangedEvent.
type OrgSMTPConfigHTTPChangedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string              `json:"id,omitempty"`
	Description           *string             `json:"description,omitempty"`
	Endpoint              *string             `json:"endpoint,omitempty"`
	SigningKey            *crypto.CryptoValue `json:"signingKey,omitempty"`
}

func (e *OrgSMTPConfigHTTPChangedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigHTTPChangedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigHTTPChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type OrgSMTPConfigHTTPChanges func(event *OrgSMTPConfigHTTPChangedEvent)

func NewOrgSMTPConfigHTTPChangeEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	changes []OrgSMTPConfigHTTPChanges,
) (*OrgSMTPConfigHTTPChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "ORG-h0pWf", "Errors.NoChangesFound")
	}
	changeEvent := &OrgSMTPConfigHTTPChangedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigHTTPChangedEventType,
		),
		ID: id,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

// OrgSMTPConfigActivatedEvent is the org-level equivalent of instance.SMTPConfigActivatedEvent.
type OrgSMTPConfigActivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewOrgSMTPConfigActivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *OrgSMTPConfigActivatedEvent {
	return &OrgSMTPConfigActivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigActivatedEventType,
		),
		ID: id,
	}
}

func (e *OrgSMTPConfigActivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigActivatedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigActivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// OrgSMTPConfigDeactivatedEvent is the org-level equivalent of instance.SMTPConfigDeactivatedEvent.
type OrgSMTPConfigDeactivatedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewOrgSMTPConfigDeactivatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *OrgSMTPConfigDeactivatedEvent {
	return &OrgSMTPConfigDeactivatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigDeactivatedEventType,
		),
		ID: id,
	}
}

func (e *OrgSMTPConfigDeactivatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigDeactivatedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

// OrgSMTPConfigRemovedEvent is the org-level equivalent of instance.SMTPConfigRemovedEvent.
type OrgSMTPConfigRemovedEvent struct {
	*eventstore.BaseEvent `json:"-"`
	ID                    string `json:"id,omitempty"`
}

func NewOrgSMTPConfigRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *OrgSMTPConfigRemovedEvent {
	return &OrgSMTPConfigRemovedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgSMTPConfigRemovedEventType,
		),
		ID: id,
	}
}

func (e *OrgSMTPConfigRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *OrgSMTPConfigRemovedEvent) Payload() interface{} {
	return e
}

func (e *OrgSMTPConfigRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
