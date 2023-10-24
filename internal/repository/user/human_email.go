package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	emailEventPrefix                 = humanEventPrefix + "email."
	HumanEmailChangedType            = emailEventPrefix + "changed"
	HumanEmailVerifiedType           = emailEventPrefix + "verified"
	HumanEmailVerificationFailedType = emailEventPrefix + "verification.failed"
	HumanEmailCodeAddedType          = emailEventPrefix + "code.added"
	HumanEmailCodeSentType           = emailEventPrefix + "code.sent"
	HumanEmailConfirmURLAddedType    = emailEventPrefix + "confirm_url.added"
)

type HumanEmailChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	EmailAddress domain.EmailAddress `json:"email,omitempty"`
}

func (e *HumanEmailChangedEvent) Payload() interface{} {
	return e
}

func (e *HumanEmailChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, emailAddress domain.EmailAddress) *HumanEmailChangedEvent {
	return &HumanEmailChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailChangedType,
		),
		EmailAddress: emailAddress,
	}
}

func HumanEmailChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	emailChangedEvent := &HumanEmailChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(emailChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return emailChangedEvent, nil
}

type HumanEmailVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *HumanEmailVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *HumanEmailVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailVerifiedEvent {
	return &HumanEmailVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailVerifiedType,
		),
	}
}

func HumanEmailVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	emailVerified := &HumanEmailVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsEmailVerified: true,
	}
	return emailVerified, nil
}

type HumanEmailVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanEmailVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *HumanEmailVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailVerificationFailedEvent {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailVerificationFailedType,
		),
	}
}

func HumanEmailVerificationFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanEmailVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanEmailCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	URLTemplate       string              `json:"url_template,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *HumanEmailCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanEmailCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanEmailCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanEmailCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *HumanEmailCodeAddedEvent {
	return NewHumanEmailCodeAddedEventV2(ctx, aggregate, code, expiry, "", false)
}

func NewHumanEmailCodeAddedEventV2(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	urlTemplate string,
	codeReturned bool,
) *HumanEmailCodeAddedEvent {
	return &HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		URLTemplate:       urlTemplate,
		CodeReturned:      codeReturned,
		TriggeredAtOrigin: http.ComposedOrigin(ctx),
	}
}

func HumanEmailCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	codeAdded := &HumanEmailCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(codeAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-3M0sd", "unable to unmarshal human email code added")
	}

	return codeAdded, nil
}

type HumanEmailCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanEmailCodeSentEvent) Payload() interface{} {
	return nil
}

func (e *HumanEmailCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanEmailCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanEmailCodeSentEvent {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanEmailCodeSentType,
		),
	}
}

func HumanEmailCodeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanEmailCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
