package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	phoneEventPrefix                 = humanEventPrefix + "phone."
	HumanPhoneChangedType            = phoneEventPrefix + "changed"
	HumanPhoneRemovedType            = phoneEventPrefix + "removed"
	HumanPhoneVerifiedType           = phoneEventPrefix + "verified"
	HumanPhoneVerificationFailedType = phoneEventPrefix + "verification.failed"
	HumanPhoneCodeAddedType          = phoneEventPrefix + "code.added"
	HumanPhoneCodeSentType           = phoneEventPrefix + "code.sent"
)

type HumanPhoneChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	PhoneNumber domain.PhoneNumber `json:"phone,omitempty"`
}

func (e *HumanPhoneChangedEvent) Payload() interface{} {
	return e
}

func (e *HumanPhoneChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPhoneChangedEvent(ctx context.Context, aggregate *eventstore.Aggregate, phone domain.PhoneNumber) *HumanPhoneChangedEvent {
	return &HumanPhoneChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneChangedType,
		),
		PhoneNumber: phone,
	}
}

func HumanPhoneChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	phoneChangedEvent := &HumanPhoneChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(phoneChangedEvent)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-5M0pd", "unable to unmarshal human phone changed")
	}

	return phoneChangedEvent, nil
}

type HumanPhoneRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPhoneRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanPhoneRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPhoneRemovedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanPhoneRemovedEvent {
	return &HumanPhoneRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneRemovedType,
		),
	}
}

func HumanPhoneRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanPhoneRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPhoneVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsPhoneVerified bool `json:"-"`
}

func (e *HumanPhoneVerifiedEvent) Payload() interface{} {
	return nil
}

func (e *HumanPhoneVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPhoneVerifiedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanPhoneVerifiedEvent {
	return &HumanPhoneVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneVerifiedType,
		),
	}
}

func HumanPhoneVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanPhoneVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsPhoneVerified: true,
	}, nil
}

type HumanPhoneVerificationFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanPhoneVerificationFailedEvent) Payload() interface{} {
	return nil
}

func (e *HumanPhoneVerificationFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPhoneVerificationFailedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *HumanPhoneVerificationFailedEvent {
	return &HumanPhoneVerificationFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneVerificationFailedType,
		),
	}
}

func HumanPhoneVerificationFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanPhoneVerificationFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanPhoneCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	CodeReturned      bool                `json:"code_returned,omitempty"`
	GeneratorID       string              `json:"generatorId,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
}

func (e *HumanPhoneCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanPhoneCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanPhoneCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanPhoneCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	generatorID string,
) *HumanPhoneCodeAddedEvent {
	return NewHumanPhoneCodeAddedEventV2(ctx, aggregate, code, expiry, false, generatorID)
}

func NewHumanPhoneCodeAddedEventV2(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	codeReturned bool,
	generatorID string,
) *HumanPhoneCodeAddedEvent {
	return &HumanPhoneCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		CodeReturned:      codeReturned,
		GeneratorID:       generatorID,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

func HumanPhoneCodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	codeAdded := &HumanPhoneCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(codeAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-6Ms9d", "unable to unmarshal human phone code added")
	}

	return codeAdded, nil
}

type HumanPhoneCodeSentEvent struct {
	*eventstore.BaseEvent `json:"-"`

	GeneratorInfo *senders.CodeGeneratorInfo `json:"generatorInfo,omitempty"`
}

func (e *HumanPhoneCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *HumanPhoneCodeSentEvent) Payload() interface{} {
	return e
}

func (e *HumanPhoneCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanPhoneCodeSentEvent(ctx context.Context, aggregate *eventstore.Aggregate, generatorInfo *senders.CodeGeneratorInfo) *HumanPhoneCodeSentEvent {
	return &HumanPhoneCodeSentEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPhoneCodeSentType,
		),
		GeneratorInfo: generatorInfo,
	}
}
