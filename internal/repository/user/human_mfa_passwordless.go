package user

import (
	"context"
	"encoding/json"
	"time"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	passwordlessEventPrefix                     = humanEventPrefix + "passwordless."
	humanPasswordlessTokenEventPrefix           = passwordlessEventPrefix + "token."
	HumanPasswordlessTokenAddedType             = humanPasswordlessTokenEventPrefix + "added"
	HumanPasswordlessTokenVerifiedType          = humanPasswordlessTokenEventPrefix + "verified"
	HumanPasswordlessTokenSignCountChangedType  = humanPasswordlessTokenEventPrefix + "signcount.changed"
	HumanPasswordlessTokenRemovedType           = humanPasswordlessTokenEventPrefix + "removed"
	HumanPasswordlessTokenBeginLoginType        = humanPasswordlessTokenEventPrefix + "begin.login"
	HumanPasswordlessTokenCheckSucceededType    = humanPasswordlessTokenEventPrefix + "check.succeeded"
	HumanPasswordlessTokenCheckFailedType       = humanPasswordlessTokenEventPrefix + "check.failed"
	humanPasswordlessInitCodePrefix             = passwordlessEventPrefix + "initialization.code."
	HumanPasswordlessInitCodeAddedType          = humanPasswordlessInitCodePrefix + "added"
	HumanPasswordlessInitCodeRequestedType      = humanPasswordlessInitCodePrefix + "requested"
	HumanPasswordlessInitCodeSentType           = humanPasswordlessInitCodePrefix + "sent"
	HumanPasswordlessInitCodeCheckFailedType    = humanPasswordlessInitCodePrefix + "check.failed"
	HumanPasswordlessInitCodeCheckSucceededType = humanPasswordlessInitCodePrefix + "check.succeeded"
)

type HumanPasswordlessAddedEvent struct {
	HumanWebAuthNAddedEvent
}

func NewHumanPasswordlessAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID,
	challenge string,
) *HumanPasswordlessAddedEvent {
	return &HumanPasswordlessAddedEvent{
		HumanWebAuthNAddedEvent: *NewHumanWebAuthNAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenAddedType,
			),
			webAuthNTokenID,
			challenge,
		),
	}
}

func HumanPasswordlessAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessAddedEvent{HumanWebAuthNAddedEvent: *e.(*HumanWebAuthNAddedEvent)}, nil
}

type HumanPasswordlessVerifiedEvent struct {
	HumanWebAuthNVerifiedEvent
}

func NewHumanPasswordlessVerifiedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
	userAgentID string,
) *HumanPasswordlessVerifiedEvent {
	return &HumanPasswordlessVerifiedEvent{
		HumanWebAuthNVerifiedEvent: *NewHumanWebAuthNVerifiedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenVerifiedType,
			),
			webAuthNTokenID,
			webAuthNTokenName,
			attestationType,
			keyID,
			publicKey,
			aaguid,
			signCount,
			userAgentID,
		),
	}
}

func HumanPasswordlessVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNVerifiedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessVerifiedEvent{HumanWebAuthNVerifiedEvent: *e.(*HumanWebAuthNVerifiedEvent)}, nil
}

type HumanPasswordlessSignCountChangedEvent struct {
	HumanWebAuthNSignCountChangedEvent
}

func NewHumanPasswordlessSignCountChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID string,
	signCount uint32,
) *HumanPasswordlessSignCountChangedEvent {
	return &HumanPasswordlessSignCountChangedEvent{
		HumanWebAuthNSignCountChangedEvent: *NewHumanWebAuthNSignCountChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenSignCountChangedType,
			),
			webAuthNTokenID,
			signCount,
		),
	}
}

func HumanPasswordlessSignCountChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNSignCountChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessSignCountChangedEvent{HumanWebAuthNSignCountChangedEvent: *e.(*HumanWebAuthNSignCountChangedEvent)}, nil
}

type HumanPasswordlessRemovedEvent struct {
	HumanWebAuthNRemovedEvent
}

func PrepareHumanPasswordlessRemovedEvent(ctx context.Context, webAuthNTokenID string) func(*eventstore.Aggregate) eventstore.EventPusher {
	return func(a *eventstore.Aggregate) eventstore.EventPusher {
		return NewHumanPasswordlessRemovedEvent(ctx, a, webAuthNTokenID)
	}
}

func NewHumanPasswordlessRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID string,
) *HumanPasswordlessRemovedEvent {
	return &HumanPasswordlessRemovedEvent{
		HumanWebAuthNRemovedEvent: *NewHumanWebAuthNRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenRemovedType,
			),
			webAuthNTokenID,
		),
	}
}

func HumanPasswordlessRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessRemovedEvent{HumanWebAuthNRemovedEvent: *e.(*HumanWebAuthNRemovedEvent)}, nil
}

type HumanPasswordlessBeginLoginEvent struct {
	HumanWebAuthNBeginLoginEvent
}

func NewHumanPasswordlessBeginLoginEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	challenge string,
	allowedCredentialIDs [][]byte,
	userVerification domain.UserVerificationRequirement,
	info *AuthRequestInfo,
) *HumanPasswordlessBeginLoginEvent {
	return &HumanPasswordlessBeginLoginEvent{
		HumanWebAuthNBeginLoginEvent: *NewHumanWebAuthNBeginLoginEvent(eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessTokenBeginLoginType,
		),
			challenge,
			allowedCredentialIDs,
			userVerification,
			info),
	}
}

func HumanPasswordlessBeginLoginEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNBeginLoginEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessBeginLoginEvent{HumanWebAuthNBeginLoginEvent: *e.(*HumanWebAuthNBeginLoginEvent)}, nil
}

type HumanPasswordlessCheckSucceededEvent struct {
	HumanWebAuthNCheckSucceededEvent
}

func NewHumanPasswordlessCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *HumanPasswordlessCheckSucceededEvent {
	return &HumanPasswordlessCheckSucceededEvent{
		HumanWebAuthNCheckSucceededEvent: *NewHumanWebAuthNCheckSucceededEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenCheckSucceededType,
			),
			info,
		),
	}
}

func HumanPasswordlessCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNCheckSucceededEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessCheckSucceededEvent{HumanWebAuthNCheckSucceededEvent: *e.(*HumanWebAuthNCheckSucceededEvent)}, nil
}

type HumanPasswordlessCheckFailedEvent struct {
	HumanWebAuthNCheckFailedEvent
}

func NewHumanPasswordlessCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *HumanPasswordlessCheckFailedEvent {
	return &HumanPasswordlessCheckFailedEvent{
		HumanWebAuthNCheckFailedEvent: *NewHumanWebAuthNCheckFailedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenCheckFailedType,
			),
			info,
		),
	}
}

func HumanPasswordlessCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNCheckFailedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanPasswordlessCheckFailedEvent{HumanWebAuthNCheckFailedEvent: *e.(*HumanWebAuthNCheckFailedEvent)}, nil
}

type HumanPasswordlessInitCodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID     string              `json:"id"`
	Code   *crypto.CryptoValue `json:"code"`
	Expiry time.Duration       `json:"expiry"`
}

func (e *HumanPasswordlessInitCodeAddedEvent) Data() interface{} {
	return e
}

func (e *HumanPasswordlessInitCodeAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanPasswordlessInitCodeAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *HumanPasswordlessInitCodeAddedEvent {
	return &HumanPasswordlessInitCodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessInitCodeAddedType,
		),
		ID:     id,
		Code:   code,
		Expiry: expiry,
	}
}

func HumanPasswordlessInitCodeAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanPasswordlessInitCodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-BDf32", "unable to unmarshal human passwordless code added")
	}
	return webAuthNAdded, nil
}

type HumanPasswordlessInitCodeRequestedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID     string              `json:"id"`
	Code   *crypto.CryptoValue `json:"code"`
	Expiry time.Duration       `json:"expiry"`
}

func (e *HumanPasswordlessInitCodeRequestedEvent) Data() interface{} {
	return e
}

func (e *HumanPasswordlessInitCodeRequestedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanPasswordlessInitCodeRequestedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
	code *crypto.CryptoValue,
	expiry time.Duration,
) *HumanPasswordlessInitCodeRequestedEvent {
	return &HumanPasswordlessInitCodeRequestedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessInitCodeRequestedType,
		),
		ID:     id,
		Code:   code,
		Expiry: expiry,
	}
}

func HumanPasswordlessInitCodeRequestedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanPasswordlessInitCodeRequestedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-VGfg3", "unable to unmarshal human passwordless code delivery added")
	}
	return webAuthNAdded, nil
}

type HumanPasswordlessInitCodeSentEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id"`
}

func (e *HumanPasswordlessInitCodeSentEvent) Data() interface{} {
	return e
}

func (e *HumanPasswordlessInitCodeSentEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanPasswordlessInitCodeSentEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *HumanPasswordlessInitCodeSentEvent {
	return &HumanPasswordlessInitCodeSentEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessInitCodeSentType,
		),
		ID: id,
	}
}

func HumanPasswordlessInitCodeSentEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanPasswordlessInitCodeSentEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Gtg4j", "unable to unmarshal human passwordless code sent")
	}
	return webAuthNAdded, nil
}

type HumanPasswordlessInitCodeCheckFailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id"`
}

func (e *HumanPasswordlessInitCodeCheckFailedEvent) Data() interface{} {
	return e
}

func (e *HumanPasswordlessInitCodeCheckFailedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanPasswordlessInitCodeCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *HumanPasswordlessInitCodeCheckFailedEvent {
	return &HumanPasswordlessInitCodeCheckFailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessInitCodeCheckFailedType,
		),
		ID: id,
	}
}

func HumanPasswordlessInitCodeCodeCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanPasswordlessInitCodeCheckFailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Gtg4j", "unable to unmarshal human passwordless code check failed")
	}
	return webAuthNAdded, nil
}

type HumanPasswordlessInitCodeCheckSucceededEvent struct {
	eventstore.BaseEvent `json:"-"`

	ID string `json:"id"`
}

func (e *HumanPasswordlessInitCodeCheckSucceededEvent) Data() interface{} {
	return e
}

func (e *HumanPasswordlessInitCodeCheckSucceededEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewHumanPasswordlessInitCodeCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	id string,
) *HumanPasswordlessInitCodeCheckSucceededEvent {
	return &HumanPasswordlessInitCodeCheckSucceededEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			HumanPasswordlessInitCodeCheckSucceededType,
		),
		ID: id,
	}
}

func HumanPasswordlessInitCodeCodeCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	webAuthNAdded := &HumanPasswordlessInitCodeCheckSucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, webAuthNAdded)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-Gtg4j", "unable to unmarshal human passwordless code check succeeded")
	}
	return webAuthNAdded, nil
}
