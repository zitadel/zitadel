package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = IamEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = IamEventTypePrefix + policy.LabelPolicyChangedEventType
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
	primaryColor,
	secondaryColor string,
) *LabelPolicyChangedEvent {
	return &LabelPolicyChangedEvent{
		*policy.NewLabelPolicyChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				LabelPolicyChangedEventType,
			),
			primaryColor,
			secondaryColor,
		),
	}
}

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LabelPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *e.(*policy.LabelPolicyChangedEvent)}, nil
}
