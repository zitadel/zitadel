package multi_factors

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/login/multi_factors"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

var (
	iamEventPrefix                         = eventstore.EventType("iam.")
	LoginPolicyMultiFactorAddedEventType   = iamEventPrefix + multi_factors.LoginPolicyMultiFactorAddedEventType
	LoginPolicyMultiFactorRemovedEventType = iamEventPrefix + multi_factors.LoginPolicyMultiFactorRemovedEventType
)

type LoginPolicyMultiFactorAddedEvent struct {
	multi_factors.MultiFactorAddedEvent
}

func NewLoginPolicyMultiFactorAddedEvent(
	ctx context.Context,
	mfaType multi_factors.MultiFactorType,
) *LoginPolicyMultiFactorAddedEvent {
	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *multi_factors.NewMultiFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorAddedEventType),
			mfaType),
	}
}

func MultiFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := multi_factors.MultiFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *e.(*multi_factors.MultiFactorAddedEvent),
	}, nil
}

type LoginPolicyMultiFactorRemovedEvent struct {
	multi_factors.MultiFactorRemovedEvent
}

func NewLoginPolicyMultiFactorRemovedEvent(
	ctx context.Context,
	mfaType multi_factors.MultiFactorType,
) *LoginPolicyMultiFactorRemovedEvent {

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *multi_factors.NewMultiFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorRemovedEventType),
			mfaType),
	}
}

func MultiFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := multi_factors.MultiFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *e.(*multi_factors.MultiFactorRemovedEvent),
	}, nil
}
