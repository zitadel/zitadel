package schemauser

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	emailEventPrefix            = eventPrefix + "email."
	EmailUpdatedType            = emailEventPrefix + "updated"
	EmailVerifiedType           = emailEventPrefix + "verified"
	EmailVerificationFailedType = emailEventPrefix + "verification.failed"
	EmailCodeAddedType          = emailEventPrefix + "code.added"
	EmailCodeSentType           = emailEventPrefix + "code.sent"
)

type EmailUpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	EmailAddress domain.EmailAddress `json:"email,omitempty"`
}

func (e *EmailUpdatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *EmailUpdatedEvent) Payload() interface{} {
	return e
}

func (e *EmailUpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailUpdatedEvent(ctx context.Context, aggregate *eventstore.Aggregate, emailAddress domain.EmailAddress) *EmailUpdatedEvent {
	return &EmailUpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailUpdatedType,
		),
		EmailAddress: emailAddress,
	}
}

type EmailVerifiedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *EmailVerifiedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *EmailVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *EmailVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailVerifiedEvent {
	return &EmailVerifiedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailVerifiedType,
		),
	}
}

type EmailVerificationFailedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *EmailVerificationFailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}
func (e *EmailVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *EmailVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailVerificationFailedEvent {
	return &EmailVerificationFailedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailVerificationFailedType,
		),
	}
}

type EmailCodeAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	URLTemplate       string              `json:"url_template,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *EmailCodeAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *EmailCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *EmailCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *EmailCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewEmailCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	urlTemplate string,
	codeReturned bool,
) *EmailCodeAddedEvent {
	return &EmailCodeAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		URLTemplate:       urlTemplate,
		CodeReturned:      codeReturned,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type EmailCodeSentEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *EmailCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}
func (e *EmailCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *EmailCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailCodeSentEvent {
	return &EmailCodeSentEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailCodeSentType,
		),
	}
}
