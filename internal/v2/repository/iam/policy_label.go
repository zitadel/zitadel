package iam

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = iamEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = iamEventTypePrefix + policy.LabelPolicyChangedEventType
)

type LabelPolicyReadModel struct{ policy.LabelPolicyReadModel }

func (rm *LabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *policy.LabelPolicyAddedEvent, *policy.LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LabelPolicyAddedEvent struct {
	policy.LabelPolicyAddedEvent
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

func LabelPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LabelPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LabelPolicyChangedEvent{LabelPolicyChangedEvent: *e.(*policy.LabelPolicyChangedEvent)}, nil
}

// func NewLabelPolicyAddedEvent(
// 	ctx context.Context,
// 	primaryColor,
// 	secondaryColor string,
// ) *LabelPolicyAddedEvent {

// 	return &LabelPolicyAddedEvent{
// 		LabelPolicyAddedEvent: *policy.NewLabelPolicyAddedEvent(
// 			ctx,
// 			primaryColor,
// 			secondaryColor,
// 		),
// 	}
// }

// func NewLabelPolicyChangedEvent(
// 	ctx context.Context,
// 	primaryColor,
// 	secondaryColor string,
// ) *MemberChangedEvent {

// 	return &LabelPolicyChangedEvent{
// 		LabelPolicyChangedEvent: *policy.NewLabelPolicyChangedEvent(
// 			ctx,
// 			primaryColor,
// 			secondaryColor,
// 		),
// 	}
// }
