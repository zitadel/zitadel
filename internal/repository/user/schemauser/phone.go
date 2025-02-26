package schemauser

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
)

const (
	phoneEventPrefix            = eventPrefix + "phone."
	PhoneUpdatedType            = phoneEventPrefix + "updated"
	PhoneVerifiedType           = phoneEventPrefix + "verified"
	PhoneVerificationFailedType = phoneEventPrefix + "verification.failed"
	PhoneCodeAddedType          = phoneEventPrefix + "code.added"
	PhoneCodeSentType           = phoneEventPrefix + "code.sent"
)

type PhoneUpdatedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`
}

func (e *PhoneUpdatedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PhoneUpdatedEvent) Payload() interface{} {
	return e
}

func (e *PhoneUpdatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneUpdatedEvent(ctx context.Context, aggregate *eventstore.Aggregate, phone domain.PhoneNumber) *PhoneUpdatedEvent {
	return &PhoneUpdatedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneUpdatedType,
		),
		PhoneNumber: phone,
	}
}

type PhoneVerifiedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *PhoneVerifiedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}
func (e *PhoneVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *PhoneVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneVerifiedEvent {
	return &PhoneVerifiedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneVerifiedType,
		),
	}
}

type PhoneVerificationFailedEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *PhoneVerificationFailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *PhoneVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *PhoneVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPhoneVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *PhoneVerificationFailedEvent {
	return &PhoneVerificationFailedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneVerificationFailedType,
		),
	}
}

type PhoneCodeAddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	GeneratorID       string              `json:"generatorId,omitempty"`
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

func (e *PhoneCodeAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func NewPhoneCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	codeReturned bool,
	generatorID string,
) *PhoneCodeAddedEvent {
	return &PhoneCodeAddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		CodeReturned:      codeReturned,
		GeneratorID:       generatorID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type PhoneCodeSentEvent struct {
	*eventstore.BaseEvent `json:"-"`

	GeneratorInfo *senders.CodeGeneratorInfo `json:"generatorInfo,omitempty"`
}

func (e *PhoneCodeSentEvent) Payload() interface{} {
	return e
}

func (e *PhoneCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *PhoneCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func NewPhoneCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate, generatorInfo *senders.CodeGeneratorInfo) *PhoneCodeSentEvent {
	return &PhoneCodeSentEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PhoneCodeSentType,
		),
		GeneratorInfo: generatorInfo,
	}
}
