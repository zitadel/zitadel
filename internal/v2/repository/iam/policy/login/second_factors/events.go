package second_factors

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/second_factors"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

var (
	iamEventPrefix                          = eventstore.EventType("iam.")
	LoginPolicySecondFactorAddedEventType   = iamEventPrefix + second_factors.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactorRemovedEventType = iamEventPrefix + second_factors.LoginPolicySecondFactorRemovedEventType
)

type LoginPolicySecondFactorAddedEvent struct {
	second_factors.SecondFactorAddedEvent
}

func NewLoginPolicySecondFactorAddedEvent(
	ctx context.Context,
	mfaType second_factors.SecondFactorType,
) *LoginPolicySecondFactorAddedEvent {
	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *second_factors.NewSecondFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorAddedEventType),
			mfaType),
	}
}

func SecondFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := second_factors.SecondFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *e.(*second_factors.SecondFactorAddedEvent),
	}, nil
}

type LoginPolicySecondFactorRemovedEvent struct {
	second_factors.SecondFactorRemovedEvent
}

func NewLoginPolicySecondFactorRemovedEvent(
	ctx context.Context,
	mfaType second_factors.SecondFactorType,
) *LoginPolicySecondFactorRemovedEvent {

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *second_factors.NewSecondFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorRemovedEventType),
			mfaType),
	}
}

func SecondFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := second_factors.SecondFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *e.(*second_factors.SecondFactorRemovedEvent),
	}, nil
}
