package factors

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/factors"
)

var (
	iamEventPrefix                          = eventstore.EventType("iam.")
	LoginPolicySecondFactorAddedEventType   = iamEventPrefix + factors.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactorRemovedEventType = iamEventPrefix + factors.LoginPolicySecondFactorRemovedEventType

	LoginPolicyMultiFactorAddedEventType   = iamEventPrefix + factors.LoginPolicyMultiFactorAddedEventType
	LoginPolicyMultiFactorRemovedEventType = iamEventPrefix + factors.LoginPolicyMultiFactorRemovedEventType
)

type LoginPolicySecondFactorAddedEvent struct {
	factors.SecondFactorAddedEvent
}

func NewLoginPolicySecondFactorAddedEvent(
	ctx context.Context,
	mfaType factors.SecondFactorType,
) *LoginPolicySecondFactorAddedEvent {
	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *factors.NewSecondFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorAddedEventType),
			mfaType),
	}
}

func SecondFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := factors.SecondFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *e.(*factors.SecondFactorAddedEvent),
	}, nil
}

type LoginPolicySecondFactorRemovedEvent struct {
	factors.SecondFactorRemovedEvent
}

func NewLoginPolicySecondFactorRemovedEvent(
	ctx context.Context,
	mfaType factors.SecondFactorType,
) *LoginPolicySecondFactorRemovedEvent {

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *factors.NewSecondFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorRemovedEventType),
			mfaType),
	}
}

func SecondFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := factors.SecondFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *e.(*factors.SecondFactorRemovedEvent),
	}, nil
}

type LoginPolicyMultiFactorAddedEvent struct {
	factors.MultiFactorAddedEvent
}

func NewLoginPolicyMultiFactorAddedEvent(
	ctx context.Context,
	mfaType factors.MultiFactorType,
) *LoginPolicyMultiFactorAddedEvent {
	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *factors.NewMultiFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorAddedEventType),
			mfaType),
	}
}

func MultiFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := factors.MultiFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *e.(*factors.MultiFactorAddedEvent),
	}, nil
}

type LoginPolicyMultiFactorRemovedEvent struct {
	factors.MultiFactorRemovedEvent
}

func NewLoginPolicyMultiFactorRemovedEvent(
	ctx context.Context,
	mfaType factors.MultiFactorType,
) *LoginPolicyMultiFactorRemovedEvent {

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *factors.NewMultiFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorRemovedEventType),
			mfaType),
	}
}

func MultiFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := factors.MultiFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *e.(*factors.MultiFactorRemovedEvent),
	}, nil
}
