package user

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	otpEventPrefix                  = mfaEventPrefix + "otp."
	HumanMFAOTPAddedType            = otpEventPrefix + "added"
	HumanMFAOTPVerifiedType         = otpEventPrefix + "verified"
	HumanMFAOTPRemovedType          = otpEventPrefix + "removed"
	HumanMFAOTPCheckSucceededType   = otpEventPrefix + "check.succeeded"
	HumanMFAOTPCheckFailedType      = otpEventPrefix + "check.failed"
	otpSMSEventPrefix               = otpEventPrefix + "sms."
	HumanOTPSMSAddedType            = otpSMSEventPrefix + "added"
	HumanOTPSMSRemovedType          = otpSMSEventPrefix + "removed"
	HumanOTPSMSCodeAddedType        = otpSMSEventPrefix + "code.added"
	HumanOTPSMSCodeSentType         = otpSMSEventPrefix + "code.sent"
	HumanOTPSMSCheckSucceededType   = otpSMSEventPrefix + "check.succeeded"
	HumanOTPSMSCheckFailedType      = otpSMSEventPrefix + "check.failed"
	otpEmailEventPrefix             = otpEventPrefix + "email."
	HumanOTPEmailAddedType          = otpEmailEventPrefix + "added"
	HumanOTPEmailRemovedType        = otpEmailEventPrefix + "removed"
	HumanOTPEmailCodeAddedType      = otpEmailEventPrefix + "code.added"
	HumanOTPEmailCodeSentType       = otpEmailEventPrefix + "code.sent"
	HumanOTPEmailCheckSucceededType = otpEmailEventPrefix + "check.succeeded"
	HumanOTPEmailCheckFailedType    = otpEmailEventPrefix + "check.failed"
)

type HumanOTPAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Secret *crypto.CryptoValue `json:"otpSecret,omitempty"`
}

func (e *HumanOTPAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	secret *crypto.CryptoValue,
) *HumanOTPAddedEvent {
	return &HumanOTPAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPAddedType,
		),
		Secret: secret,
	}
}

func HumanOTPAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp added")
	}
	return otpAdded, nil
}

type HumanOTPVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`
	UserAgentID          string `json:"userAgentID,omitempty"`
}

func (e *HumanOTPVerifiedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPVerifiedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPVerifiedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userAgentID string,
) *HumanOTPVerifiedEvent {
	return &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPVerifiedType,
		),
		UserAgentID: userAgentID,
	}
}

func HumanOTPVerifiedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	out := &HumanOTPVerifiedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	return out, nil
}

type HumanOTPRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPRemovedEvent {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPRemovedType,
		),
	}
}

func HumanOTPRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &HumanOTPRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type HumanOTPCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPCheckSucceededEvent {
	return &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

func HumanOTPCheckSucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check succeeded")
	}
	return otpAdded, nil
}

type HumanOTPCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewHumanOTPCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPCheckFailedEvent {
	return &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanMFAOTPCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}

func HumanOTPCheckFailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	otpAdded := &HumanOTPCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(otpAdded)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "USER-Ns9df", "unable to unmarshal human otp check failed")
	}
	return otpAdded, nil
}

type HumanOTPSMSAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPSMSAddedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPSMSAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPSMSAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPSMSAddedEvent {
	return &HumanOTPSMSAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSAddedType,
		),
	}
}

type HumanOTPSMSRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPSMSRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPSMSRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPSMSRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPSMSRemovedEvent {
	return &HumanOTPSMSRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSRemovedType,
		),
	}
}

type HumanOTPSMSCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
	GeneratorID       string              `json:"generatorId,omitempty"`
	*AuthRequestInfo
}

func (e *HumanOTPSMSCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPSMSCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSCodeAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *HumanOTPSMSCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanOTPSMSCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	info *AuthRequestInfo,
	generatorID string,
) *HumanOTPSMSCodeAddedEvent {
	return &HumanOTPSMSCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
		AuthRequestInfo:   info,
		GeneratorID:       generatorID,
	}
}

type HumanOTPSMSCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`

	GeneratorInfo *senders.CodeGeneratorInfo `json:"generatorInfo,omitempty"`
}

func (e *HumanOTPSMSCodeSentEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPSMSCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPSMSCodeSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	generatorInfo *senders.CodeGeneratorInfo,
) *HumanOTPSMSCodeSentEvent {
	return &HumanOTPSMSCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSCodeSentType,
		),
		GeneratorInfo: generatorInfo,
	}
}

type HumanOTPSMSCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPSMSCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPSMSCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSCheckSucceededEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPSMSCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPSMSCheckSucceededEvent {
	return &HumanOTPSMSCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

type HumanOTPSMSCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPSMSCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPSMSCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPSMSCheckFailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPSMSCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPSMSCheckFailedEvent {
	return &HumanOTPSMSCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPSMSCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}

type HumanOTPEmailAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPEmailAddedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPEmailAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPEmailAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPEmailAddedEvent {
	return &HumanOTPEmailAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailAddedType,
		),
	}
}

type HumanOTPEmailRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *HumanOTPEmailRemovedEvent) Payload() interface{} {
	return nil
}

func (e *HumanOTPEmailRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPEmailRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPEmailRemovedEvent {
	return &HumanOTPEmailRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailRemovedType,
		),
	}
}

type HumanOTPEmailCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code              *crypto.CryptoValue `json:"code,omitempty"`
	Expiry            time.Duration       `json:"expiry,omitempty"`
	TriggeredAtOrigin string              `json:"triggerOrigin,omitempty"`
	*AuthRequestInfo
}

func (e *HumanOTPEmailCodeAddedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPEmailCodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailCodeAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *HumanOTPEmailCodeAddedEvent) TriggerOrigin() string {
	return e.TriggeredAtOrigin
}

func NewHumanOTPEmailCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	code *crypto.CryptoValue,
	expiry time.Duration,
	info *AuthRequestInfo,
) *HumanOTPEmailCodeAddedEvent {
	return &HumanOTPEmailCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailCodeAddedType,
		),
		Code:              code,
		Expiry:            expiry,
		AuthRequestInfo:   info,
		TriggeredAtOrigin: http.DomainContext(ctx).Origin(),
	}
}

type HumanOTPEmailCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
	*AuthRequestInfo
}

func (e *HumanOTPEmailCodeSentEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPEmailCodeSentEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailCodeSentEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPEmailCodeSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *HumanOTPEmailCodeSentEvent {
	return &HumanOTPEmailCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailCodeSentType,
		),
	}
}

type HumanOTPEmailCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPEmailCheckSucceededEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPEmailCheckSucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailCheckSucceededEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPEmailCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPEmailCheckSucceededEvent {
	return &HumanOTPEmailCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailCheckSucceededType,
		),
		AuthRequestInfo: info,
	}
}

type HumanOTPEmailCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`
	*AuthRequestInfo
}

func (e *HumanOTPEmailCheckFailedEvent) Payload() interface{} {
	return e
}

func (e *HumanOTPEmailCheckFailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HumanOTPEmailCheckFailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func NewHumanOTPEmailCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo,
) *HumanOTPEmailCheckFailedEvent {
	return &HumanOTPEmailCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanOTPEmailCheckFailedType,
		),
		AuthRequestInfo: info,
	}
}
