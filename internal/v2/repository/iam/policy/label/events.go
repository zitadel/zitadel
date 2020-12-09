package label

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

var (
	iamEventPrefix              = eventstore.EventType("iam.")
	LabelPolicyAddedEventType   = iamEventPrefix + label.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = iamEventPrefix + label.LabelPolicyChangedEventType
)

type LabelPolicyAddedEvent struct {
	label.LabelPolicyAddedEvent
}

func NewLabelPolicyAddedEventEvent(
	ctx context.Context,
	primaryColor,
	secondaryColor string,
) *LabelPolicyAddedEvent {
	return &LabelPolicyAddedEvent{
		LabelPolicyAddedEvent: *label.NewLabelPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LabelPolicyAddedEventType),
			primaryColor,
			secondaryColor),
	}
}

func LabelPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := label.LabelPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyAddedEvent{LabelPolicyAddedEvent: *e.(*label.LabelPolicyAddedEvent)}, nil
}

type LabelPolicyChangedEvent struct {
	label.LabelPolicyChangedEvent
}

func LabelPolicyChangedEventFromExisting(
	ctx context.Context,
	current *LabelPolicyWriteModel,
	primaryColor,
	secondaryColor string,
) (*LabelPolicyChangedEvent, error) {
	event := label.NewLabelPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyChangedEventType,
		),
		&current.Policy,
		primaryColor,
		secondaryColor,
	)
	return &LabelPolicyChangedEvent{
		*event,
	}, nil
}

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := label.LabelPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *e.(*label.LabelPolicyChangedEvent)}, nil
}
