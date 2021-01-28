package user

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
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
	webAuthNTokenID,
	challenge string,
) *HumanPasswordlessAddedEvent {
	return &HumanPasswordlessAddedEvent{
		HumanWebAuthNAddedEvent: *NewHumanWebAuthNAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
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
	webAuthNTokenID,
	webAuthNTokenName,
	attestationType string,
	keyID,
	publicKey,
	aaguid []byte,
	signCount uint32,
) *HumanPasswordlessVerifiedEvent {
	return &HumanPasswordlessVerifiedEvent{
		HumanWebAuthNVerifiedEvent: *NewHumanWebAuthNVerifiedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				HumanPasswordlessTokenVerifiedType,
			),
			webAuthNTokenID,
			webAuthNTokenName,
			attestationType,
			keyID,
			publicKey,
			aaguid,
			signCount,
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
	webAuthNTokenID string,
	signCount uint32,
) *HumanPasswordlessSignCountChangedEvent {
	return &HumanPasswordlessSignCountChangedEvent{
		HumanWebAuthNSignCountChangedEvent: *NewHumanWebAuthNSignCountChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
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

func NewHumanPasswordlessRemovedEvent(
	ctx context.Context,
	webAuthNTokenID string,
) *HumanPasswordlessRemovedEvent {
	return &HumanPasswordlessRemovedEvent{
		HumanWebAuthNRemovedEvent: *NewHumanWebAuthNRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
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
	webAuthNTokenID,
	challenge string,
) *HumanPasswordlessBeginLoginEvent {
	return &HumanPasswordlessBeginLoginEvent{
		HumanWebAuthNBeginLoginEvent: *NewHumanWebAuthNBeginLoginEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				HumanPasswordlessTokenVerifiedType,
			),
			webAuthNTokenID,
			challenge,
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

func NewHumanPasswordlessCheckSucceededEvent(ctx context.Context) *HumanPasswordlessCheckSucceededEvent {
	return &HumanPasswordlessCheckSucceededEvent{
		HumanWebAuthNCheckSucceededEvent: *NewHumanWebAuthNCheckSucceededEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				HumanPasswordlessTokenCheckSucceededType,
			),
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

func NewHumanPasswordlessCheckFailedEvent(ctx context.Context) *HumanPasswordlessCheckFailedEvent {
	return &HumanPasswordlessCheckFailedEvent{
		HumanWebAuthNCheckFailedEvent: *NewHumanWebAuthNCheckFailedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				HumanPasswordlessTokenCheckFailedType,
			),
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
