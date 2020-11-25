package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + policy.LabelPolicyChangedEventType
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

type LabelPolicyChangedEvent struct {
	policy.LabelPolicyChangedEvent
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
