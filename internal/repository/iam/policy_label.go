package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	LabelPolicyAddedEventType   = iamEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = iamEventTypePrefix + policy.LabelPolicyChangedEventType
)

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
}

func NewLabelPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	primaryColor,
	secondaryColor string,
	hideLoginNameSuffix bool,
) *LabelPolicyAddedEvent {
	return &LabelPolicyAddedEvent{
		LabelPolicyAddedEvent: *policy.NewLabelPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LabelPolicyAddedEventType),
			primaryColor,
			secondaryColor,
			hideLoginNameSuffix),
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
	aggregate *eventstore.Aggregate,
	changes []policy.LabelPolicyChanges,
) (*LabelPolicyChangedEvent, error) {
	changedEvent, err := policy.NewLabelPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LabelPolicyChangedEventType),
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
