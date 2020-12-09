package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + label.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + label.LabelPolicyChangedEventType
)

type LabelPolicyReadModel struct{ label.LabelPolicyReadModel }

func (rm *LabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *label.LabelPolicyAddedEvent, *label.LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LabelPolicyAddedEvent struct {
	label.LabelPolicyAddedEvent
}

type LabelPolicyChangedEvent struct {
	label.LabelPolicyChangedEvent
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
