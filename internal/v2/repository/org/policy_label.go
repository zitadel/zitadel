package org

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/query"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LabelPolicyAddedEventType   = orgEventTypePrefix + policy.LabelPolicyAddedEventType
	LabelPolicyChangedEventType = orgEventTypePrefix + policy.LabelPolicyChangedEventType
)

type LabelPolicyReadModel struct{ query.LabelPolicyReadModel }

func (rm *LabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.LabelPolicyReadModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *LabelPolicyChangedEvent:
			rm.LabelPolicyReadModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *policy.LabelPolicyAddedEvent, *policy.LabelPolicyChangedEvent:
			rm.LabelPolicyReadModel.AppendEvents(e)
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
// ) *MemberAddedEvent {

// 	return &MemberAddedEvent{
// 		MemberAddedEvent: *policy.NewLabelPolicyAddedEvent(
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
