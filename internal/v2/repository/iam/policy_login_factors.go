package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicySecondFactorAddedEventType   = iamEventTypePrefix + policy.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactorRemovedEventType = iamEventTypePrefix + policy.LoginPolicySecondFactorRemovedEventType

	LoginPolicyMultiFactorAddedEventType   = iamEventTypePrefix + policy.LoginPolicyMultiFactorAddedEventType
	LoginPolicyMultiFactorRemovedEventType = iamEventTypePrefix + policy.LoginPolicyMultiFactorRemovedEventType
)

type LoginPolicySecondFactorAddedEvent struct {
	policy.SecondFactorAddedEvent
}

func NewLoginPolicySecondFactorAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	mfaType domain.SecondFactorType,
) *LoginPolicySecondFactorAddedEvent {
	return &LoginPolicySecondFactorAddedEvent{
		SecondFactorAddedEvent: *policy.NewSecondFactorAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicySecondFactorAddedEventType),
			mfaType),
	}
}

func SecondFactorAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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
	aggregate *eventstore.Aggregate,
	mfaType domain.SecondFactorType,
) *LoginPolicySecondFactorRemovedEvent {

	return &LoginPolicySecondFactorRemovedEvent{
		SecondFactorRemovedEvent: *policy.NewSecondFactorRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicySecondFactorRemovedEventType),
			mfaType),
	}
}

func SecondFactorRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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
	aggregate *eventstore.Aggregate,
	mfaType domain.MultiFactorType,
) *LoginPolicyMultiFactorAddedEvent {
	return &LoginPolicyMultiFactorAddedEvent{
		MultiFactorAddedEvent: *policy.NewMultiFactorAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicyMultiFactorAddedEventType),
			mfaType),
	}
}

func MultiFactorAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
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
	aggregate *eventstore.Aggregate,
	mfaType domain.MultiFactorType,
) *LoginPolicyMultiFactorRemovedEvent {

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *policy.NewMultiFactorRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicyMultiFactorRemovedEventType),
			mfaType),
	}
}

func MultiFactorRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MultiFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *e.(*policy.MultiFactorRemovedEvent),
	}, nil
}
