package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicySecondFactorAddedEventType   = IAMEventTypePrefix + policy.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactorRemovedEventType = IAMEventTypePrefix + policy.LoginPolicySecondFactorRemovedEventType

	LoginPolicyMultiFactorAddedEventType   = IAMEventTypePrefix + policy.LoginPolicyMultiFactorAddedEventType
	LoginPolicyMultiFactorRemovedEventType = IAMEventTypePrefix + policy.LoginPolicyMultiFactorRemovedEventType
)

type LoginPolicySecondFactorAddedEvent struct {
	policy.SecondFactorAddedEvent
}

func NewLoginPolicySecondFactorAddedEvent(
	ctx context.Context,
	mfaType domain.SecondFactorType,
) *LoginPolicySecondFactorAddedEvent {
	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *policy.NewSecondFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorAddedEventType),
			mfaType),
	}
}

func SecondFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.SecondFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *e.(*policy.SecondFactorAddedEvent),
	}, nil
}

type LoginPolicySecondFactorRemovedEvent struct {
	policy.SecondFactorRemovedEvent
}

func NewLoginPolicySecondFactorRemovedEvent(
	ctx context.Context,
	mfaType domain.SecondFactorType,
) *LoginPolicySecondFactorRemovedEvent {

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *policy.NewSecondFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicySecondFactorRemovedEventType),
			mfaType),
	}
}

func SecondFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.SecondFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *e.(*policy.SecondFactorRemovedEvent),
	}, nil
}

type LoginPolicyMultiFactorAddedEvent struct {
	policy.MultiFactorAddedEvent
}

func NewLoginPolicyMultiFactorAddedEvent(
	ctx context.Context,
	mfaType domain.MultiFactorType,
) *LoginPolicyMultiFactorAddedEvent {
	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *policy.NewMultiFactorAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorAddedEventType),
			mfaType),
	}
}

func MultiFactorAddedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MultiFactorAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *e.(*policy.MultiFactorAddedEvent),
	}, nil
}

type LoginPolicyMultiFactorRemovedEvent struct {
	policy.MultiFactorRemovedEvent
}

func NewLoginPolicyMultiFactorRemovedEvent(
	ctx context.Context,
	mfaType domain.MultiFactorType,
) *LoginPolicyMultiFactorRemovedEvent {

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *policy.NewMultiFactorRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LoginPolicyMultiFactorRemovedEventType),
			mfaType),
	}
}

func MultiFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MultiFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *e.(*policy.MultiFactorRemovedEvent),
	}, nil
}
