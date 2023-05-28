package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	PrivacyPolicyAddedEventType   = orgEventTypePrefix + policy.PrivacyPolicyAddedEventType
	PrivacyPolicyChangedEventType = orgEventTypePrefix + policy.PrivacyPolicyChangedEventType
	PrivacyPolicyRemovedEventType = orgEventTypePrefix + policy.PrivacyPolicyRemovedEventType
)

type PrivacyPolicyAddedEvent struct {
	policy.PrivacyPolicyAddedEvent
}

func NewPrivacyPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tosLink,
	privacyLink,
	helpLink string,
	supportEmail domain.EmailAddress,
) *PrivacyPolicyAddedEvent {
	return &PrivacyPolicyAddedEvent{
		PrivacyPolicyAddedEvent: *policy.NewPrivacyPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PrivacyPolicyAddedEventType),
			tosLink,
			privacyLink,
			helpLink,
			supportEmail),
	}
}

func PrivacyPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PrivacyPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PrivacyPolicyAddedEvent{PrivacyPolicyAddedEvent: *e.(*policy.PrivacyPolicyAddedEvent)}, nil
}

type PrivacyPolicyChangedEvent struct {
	policy.PrivacyPolicyChangedEvent
}

func NewPrivacyPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.PrivacyPolicyChanges,
) (*PrivacyPolicyChangedEvent, error) {
	changedEvent, err := policy.NewPrivacyPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PrivacyPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PrivacyPolicyChangedEvent{PrivacyPolicyChangedEvent: *changedEvent}, nil
}

func PrivacyPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PrivacyPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PrivacyPolicyChangedEvent{PrivacyPolicyChangedEvent: *e.(*policy.PrivacyPolicyChangedEvent)}, nil
}

type PrivacyPolicyRemovedEvent struct {
	policy.PrivacyPolicyRemovedEvent
}

func NewPrivacyPolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *PrivacyPolicyRemovedEvent {
	return &PrivacyPolicyRemovedEvent{
		PrivacyPolicyRemovedEvent: *policy.NewPrivacyPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PrivacyPolicyRemovedEventType),
		),
	}
}

func PrivacyPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PrivacyPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PrivacyPolicyRemovedEvent{PrivacyPolicyRemovedEvent: *e.(*policy.PrivacyPolicyRemovedEvent)}, nil
}
