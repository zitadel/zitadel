package user

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	passwordlessEventPrefix                    = humanEventPrefix + "passwordless.token."
	HumanPasswordlessTokenAddedType            = passwordlessEventPrefix + "added"
	HumanPasswordlessTokenVerifiedType         = passwordlessEventPrefix + "verified"
	HumanPasswordlessTokenSignCountChangedType = passwordlessEventPrefix + "signcount.changed"
	HumanPasswordlessTokenRemovedType          = passwordlessEventPrefix + "removed"
	HumanPasswordlessTokenBeginLoginType       = passwordlessEventPrefix + "begin.login"
	HumanPasswordlessTokenCheckSucceededType   = passwordlessEventPrefix + "check.succeeded"
	HumanPasswordlessTokenCheckFailedType      = passwordlessEventPrefix + "check.failed"
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
	info *AuthRequestInfo,
) *HumanPasswordlessBeginLoginEvent {
	return &HumanPasswordlessBeginLoginEvent{
		HumanWebAuthNBeginLoginEvent: *NewHumanWebAuthNBeginLoginEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				HumanPasswordlessTokenVerifiedType,
			),
			challenge,
			info,
		),
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
