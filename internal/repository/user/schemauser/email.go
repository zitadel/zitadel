package schemauser

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	emailEventPrefix            = eventPrefix + "email."
	EmailChangedType            = emailEventPrefix + "changed"
	EmailVerifiedType           = emailEventPrefix + "verified"
	EmailVerificationFailedType = emailEventPrefix + "verification.failed"
	EmailCodeAddedType          = emailEventPrefix + "code.added"
	EmailCodeSentType           = emailEventPrefix + "code.sent"
)

type EmailChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	EmailAddress domain.EmailAddress `json:"email,omitempty"`
}

func (e *EmailChangedEvent) Payload() interface{} {
	return e
}

func (e *EmailChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, emailAddress domain.EmailAddress) *EmailChangedEvent {
	return &EmailChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailChangedType,
		),
		EmailAddress: emailAddress,
	}
}

func EmailChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	emailChangedEvent := &EmailChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(emailChangedEvent)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return emailChangedEvent, nil
}

type EmailVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *EmailVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *EmailVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewEmailVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailVerifiedEvent {
	return &EmailVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailVerifiedType,
		),
	}
}

func HumanVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	emailVerified := &EmailVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsEmailVerified: true,
	}
	return emailVerified, nil
}

type EmailVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *EmailVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *EmailVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailVerificationFailedEvent {
	return &EmailVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailVerificationFailedType,
		),
	}
}

func EmailVerificationFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &EmailVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type EmailCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	URLTemplate       string              `json:"url_template,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
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
		BaseEvent: *eventstore.NewBaseEventForPush(
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

func EmailCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	codeAdded := &EmailCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(codeAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-3M0sd", "unable to unmarshal human email code added")
	}

	return codeAdded, nil
}

type EmailCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *EmailCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *EmailCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *EmailCodeSentEvent {
	return &EmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			EmailCodeSentType,
		),
	}
}

func EmailCodeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &EmailCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
