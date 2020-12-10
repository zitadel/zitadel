package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + label.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + label.LabelPolicyChangedEventType
)

type LabelPolicyReadModel struct{ label.ReadModel }

func (rm *LabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.AddedEvent)
		case *LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.ChangedEvent)
		case *label.AddedEvent, *label.ChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

type LabelPolicyAddedEvent struct {
	label.AddedEvent
}

type LabelPolicyChangedEvent struct {
	label.ChangedEvent
}

// func NewAddedEvent(
// 	ctx context.Context,
// 	primaryColor,
// 	secondaryColor string,
// ) *AddedEvent {

// 	return &AddedEvent{
// 		AddedEvent: *policy.NewAddedEvent(
// 			ctx,
// 			primaryColor,
// 			secondaryColor,
// 		),
// 	}
// }

// func NewChangedEvent(
// 	ctx context.Context,
// 	primaryColor,
// 	secondaryColor string,
// ) *MemberChangedEvent {

// 	return &ChangedEvent{
// 		ChangedEvent: *policy.NewChangedEvent(
// 			ctx,
// 			primaryColor,
// 			secondaryColor,
// 		),
// 	}
// }
