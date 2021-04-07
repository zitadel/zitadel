package user

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	u2fEventPrefix                    = mfaEventPrefix + "u2f.token."
	HumanU2FTokenAddedType            = u2fEventPrefix + "added"
	HumanU2FTokenVerifiedType         = u2fEventPrefix + "verified"
	HumanU2FTokenSignCountChangedType = u2fEventPrefix + "signcount.changed"
	HumanU2FTokenRemovedType          = u2fEventPrefix + "removed"
	HumanU2FTokenBeginLoginType       = u2fEventPrefix + "begin.login"
	HumanU2FTokenCheckSucceededType   = u2fEventPrefix + "check.succeeded"
	HumanU2FTokenCheckFailedType      = u2fEventPrefix + "check.failed"
)

type HumanU2FAddedEvent struct {
	HumanWebAuthNAddedEvent
}

func NewHumanU2FAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID,
	challenge string,
) *HumanU2FAddedEvent {
	return &HumanU2FAddedEvent{
		HumanWebAuthNAddedEvent: *NewHumanWebAuthNAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenAddedType,
			),
			webAuthNTokenID,
			challenge,
		),
	}
}

func HumanU2FAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FAddedEvent{HumanWebAuthNAddedEvent: *e.(*HumanWebAuthNAddedEvent)}, nil
}

type HumanU2FVerifiedEvent struct {
	HumanWebAuthNVerifiedEvent
}

func NewHumanU2FVerifiedEvent(
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
) *HumanU2FVerifiedEvent {
	return &HumanU2FVerifiedEvent{
		HumanWebAuthNVerifiedEvent: *NewHumanWebAuthNVerifiedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenVerifiedType,
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

func HumanU2FVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNVerifiedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FVerifiedEvent{HumanWebAuthNVerifiedEvent: *e.(*HumanWebAuthNVerifiedEvent)}, nil
}

type HumanU2FSignCountChangedEvent struct {
	HumanWebAuthNSignCountChangedEvent
}

func NewHumanU2FSignCountChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID string,
	signCount uint32,
) *HumanU2FSignCountChangedEvent {
	return &HumanU2FSignCountChangedEvent{
		HumanWebAuthNSignCountChangedEvent: *NewHumanWebAuthNSignCountChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenSignCountChangedType,
			),
			webAuthNTokenID,
			signCount,
		),
	}
}

func HumanU2FSignCountChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNSignCountChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FSignCountChangedEvent{HumanWebAuthNSignCountChangedEvent: *e.(*HumanWebAuthNSignCountChangedEvent)}, nil
}

type HumanU2FRemovedEvent struct {
	HumanWebAuthNRemovedEvent
}

func PrepareHumanU2FRemovedEvent(ctx context.Context, webAuthNTokenID string) func(*eventstore.Aggregate) eventstore.EventPusher {
	return func(a *eventstore.Aggregate) eventstore.EventPusher {
		return NewHumanU2FRemovedEvent(ctx, a, webAuthNTokenID)
	}
}

func NewHumanU2FRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	webAuthNTokenID string,
) *HumanU2FRemovedEvent {
	return &HumanU2FRemovedEvent{
		HumanWebAuthNRemovedEvent: *NewHumanWebAuthNRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenRemovedType,
			),
			webAuthNTokenID,
		),
	}
}

func HumanU2FRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FRemovedEvent{HumanWebAuthNRemovedEvent: *e.(*HumanWebAuthNRemovedEvent)}, nil
}

type HumanU2FBeginLoginEvent struct {
	HumanWebAuthNBeginLoginEvent
}

func NewHumanU2FBeginLoginEvent(ctx context.Context, aggregate *eventstore.Aggregate, challenge string, allowedCredentialIDs [][]byte, userVerification domain.UserVerificationRequirement, info *AuthRequestInfo) *HumanU2FBeginLoginEvent {
	return &HumanU2FBeginLoginEvent{
		HumanWebAuthNBeginLoginEvent: *NewHumanWebAuthNBeginLoginEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenBeginLoginType,
			),
			challenge,
			allowedCredentialIDs,
			userVerification,
			info,
		),
	}
}

func HumanU2FBeginLoginEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNBeginLoginEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FBeginLoginEvent{HumanWebAuthNBeginLoginEvent: *e.(*HumanWebAuthNBeginLoginEvent)}, nil
}

type HumanU2FCheckSucceededEvent struct {
	HumanWebAuthNCheckSucceededEvent
}

func NewHumanU2FCheckSucceededEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *HumanU2FCheckSucceededEvent {
	return &HumanU2FCheckSucceededEvent{
		HumanWebAuthNCheckSucceededEvent: *NewHumanWebAuthNCheckSucceededEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenCheckSucceededType,
			),
			info,
		),
	}
}

func HumanU2FCheckSucceededEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNCheckSucceededEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FCheckSucceededEvent{HumanWebAuthNCheckSucceededEvent: *e.(*HumanWebAuthNCheckSucceededEvent)}, nil
}

type HumanU2FCheckFailedEvent struct {
	HumanWebAuthNCheckFailedEvent
}

func NewHumanU2FCheckFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	info *AuthRequestInfo) *HumanU2FCheckFailedEvent {
	return &HumanU2FCheckFailedEvent{
		HumanWebAuthNCheckFailedEvent: *NewHumanWebAuthNCheckFailedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanU2FTokenCheckFailedType,
			),
			info,
		),
	}
}

func HumanU2FCheckFailedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := HumanWebAuthNCheckFailedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &HumanU2FCheckFailedEvent{HumanWebAuthNCheckFailedEvent: *e.(*HumanWebAuthNCheckFailedEvent)}, nil
}
