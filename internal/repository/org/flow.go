package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/flow"
)

var (
	TriggerActionsSetEventType            = orgEventTypePrefix + flow.TriggerActionsSetEventType
	TriggerActionsCascadeRemovedEventType = orgEventTypePrefix + flow.TriggerActionsCascadeRemovedEventType
	FlowClearedEventType                  = orgEventTypePrefix + flow.FlowClearedEventType
)

type TriggerActionsSetEvent struct {
	flow.TriggerActionsSetEvent
}

func NewTriggerActionsSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	flowType domain.FlowType,
	triggerType domain.TriggerType,
	actionIDs []string,
) *TriggerActionsSetEvent {
	return &TriggerActionsSetEvent{
		TriggerActionsSetEvent: *flow.NewTriggerActionsSetEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				TriggerActionsSetEventType),
			flowType,
			triggerType,
			actionIDs),
	}
}

func TriggerActionsSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := flow.TriggerActionsSetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &TriggerActionsSetEvent{TriggerActionsSetEvent: *e.(*flow.TriggerActionsSetEvent)}, nil
}

type TriggerActionsCascadeRemovedEvent struct {
	flow.TriggerActionsCascadeRemovedEvent
}

func NewTriggerActionsCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	flowType domain.FlowType,
	actionID string,
) *TriggerActionsCascadeRemovedEvent {
	return &TriggerActionsCascadeRemovedEvent{
		TriggerActionsCascadeRemovedEvent: *flow.NewTriggerActionsCascadeRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				TriggerActionsCascadeRemovedEventType),
			flowType,
			actionID),
	}
}

func TriggerActionsCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := flow.TriggerActionsCascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &TriggerActionsCascadeRemovedEvent{TriggerActionsCascadeRemovedEvent: *e.(*flow.TriggerActionsCascadeRemovedEvent)}, nil
}

type FlowClearedEvent struct {
	flow.FlowClearedEvent
}

func NewFlowClearedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	flowType domain.FlowType,
) *FlowClearedEvent {
	return &FlowClearedEvent{
		FlowClearedEvent: *flow.NewFlowClearedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				FlowClearedEventType),
			flowType),
	}
}

func FlowClearedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := flow.FlowClearedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &FlowClearedEvent{FlowClearedEvent: *e.(*flow.FlowClearedEvent)}, nil
}
