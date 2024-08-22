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
	phoneEventPrefix            = eventPrefix + "phone."
	PhoneChangedType            = phoneEventPrefix + "changed"
	PhoneRemovedType            = phoneEventPrefix + "removed"
	PhoneVerifiedType           = phoneEventPrefix + "verified"
	PhoneVerificationFailedType = phoneEventPrefix + "verification.failed"
	PhoneCodeAddedType          = phoneEventPrefix + "code.added"
	PhoneCodeSentType           = phoneEventPrefix + "code.sent"
)

type PhoneChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`
}

func (e *PhoneChangedEvent) Payload() interface{} {
	return e
}

func (e *PhoneChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, phone domain.PhoneNumber) *PhoneChangedEvent {
	return &PhoneChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneChangedType,
		),
		PhoneNumber: phone,
	}
}

func PhoneChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	phoneChangedEvent := &PhoneChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(phoneChangedEvent)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal  phone changed")
	}

	return phoneChangedEvent, nil
}

type PhoneRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PhoneRemovedEvent) Payload() interface{} {
	return nil
}

func (e *PhoneRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneRemovedEvent {
	return &PhoneRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneRemovedType,
		),
	}
}

func PhoneRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PhoneRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type PhoneVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *PhoneVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *PhoneVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneVerifiedEvent {
	return &PhoneVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneVerifiedType,
		),
	}
}

func PhoneVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PhoneVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsPhoneVerified: true,
	}, nil
}

type PhoneVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PhoneVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *PhoneVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneVerificationFailedEvent {
	return &PhoneVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneVerificationFailedType,
		),
	}
}

func PhoneVerificationFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PhoneVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type PhoneCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *PhoneCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *PhoneCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PhoneCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewPhoneCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	codeReturned bool,
) *PhoneCodeAddedEvent {
	return &PhoneCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		CodeReturned:      codeReturned,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

func PhoneCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	codeAdded := &PhoneCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(codeAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-6Ms9d", "unable to unmarshal  phone code added")
	}

	return codeAdded, nil
}

type PhoneCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PhoneCodeSentEvent) Payload() interface{} {
	return e
}

func (e *PhoneCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneCodeSentEvent {
	return &PhoneCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneCodeSentType,
		),
	}
}

func PhoneCodeSentEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PhoneCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
