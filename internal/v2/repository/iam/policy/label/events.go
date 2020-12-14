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

type AddedEvent struct {
	label.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	primaryColor,
	secondaryColor string,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *label.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, LabelPolicyAddedEventType),
			primaryColor,
			secondaryColor),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := label.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*label.AddedEvent)}, nil
}

type ChangedEvent struct {
	label.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	primaryColor,
	secondaryColor string,
) (*ChangedEvent, error) {
	event := label.NewChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			LabelPolicyChangedEventType,
		),
		&current.WriteModel,
		primaryColor,
		secondaryColor,
	)
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := label.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*label.ChangedEvent)}, nil
}
