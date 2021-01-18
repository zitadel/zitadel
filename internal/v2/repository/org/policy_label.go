package org

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + policy.LabelPolicyChangedEventType
	LabelPolicyRemovedEventType = orgEventTypePrefix + policy.LabelPolicyRemovedEventType
)

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
}

func NewLabelPolicyAddedEvent(
	ctx context.Context,
	primaryColor,
	secondaryColor string,
) *LabelPolicyAddedEvent {
	return &LabelPolicyAddedEvent{
		LabelPolicyAddedEvent: *policy.NewLabelPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LabelPolicyAddedEventType),
			primaryColor,
			secondaryColor),
	}
}

func LabelPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LabelPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyAddedEvent{LabelPolicyAddedEvent: *e.(*policy.LabelPolicyAddedEvent)}, nil
}

type LabelPolicyChangedEvent struct {
	policy.LabelPolicyChangedEvent
}

func NewLabelPolicyChangedEvent(
	ctx context.Context,
	changes []policy.LabelPolicyChanges,
) (*LabelPolicyChangedEvent, error) {
	changedEvent, err := policy.NewLabelPolicyChangedEvent(
		eventstore.NewBaseEventForPush(ctx, LabelPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *changedEvent}, nil
}

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LabelPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *e.(*policy.LabelPolicyChangedEvent)}, nil
}

type LabelPolicyRemovedEvent struct {
	policy.LabelPolicyRemovedEvent
}

func NewLabelPolicyRemovedEvent(
	ctx context.Context,
) *LabelPolicyRemovedEvent {
	return &LabelPolicyRemovedEvent{
		LabelPolicyRemovedEvent: *policy.NewLabelPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, LabelPolicyRemovedEventType),
		),
	}
}

func LabelPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LabelPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyRemovedEvent{LabelPolicyRemovedEvent: *e.(*policy.LabelPolicyRemovedEvent)}, nil
}
