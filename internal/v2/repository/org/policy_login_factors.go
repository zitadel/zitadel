package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicySecondFactorAddedEventType   = orgEventTypePrefix + policy.LoginPolicySecondFactorAddedEventType
	LoginPolicySecondFactorRemovedEventType = orgEventTypePrefix + policy.LoginPolicySecondFactorRemovedEventType

	LoginPolicyMultiFactorAddedEventType   = orgEventTypePrefix + policy.LoginPolicyMultiFactorAddedEventType
	LoginPolicyMultiFactorRemovedEventType = orgEventTypePrefix + policy.LoginPolicyMultiFactorRemovedEventType
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

func MultiFactorRemovedEventEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MultiFactorRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyMultiFactorRemovedEvent{
		MultiFactorRemovedEvent: *e.(*policy.MultiFactorRemovedEvent),
	}, nil
}
