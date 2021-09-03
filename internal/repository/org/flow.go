package org

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/flow"
)

var (
	TriggerActionsSetEventType = orgEventTypePrefix + flow.TriggerActionsSetEventType
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

func TriggerActionsSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := flow.TriggerActionsSetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &TriggerActionsSetEvent{TriggerActionsSetEvent: *e.(*flow.TriggerActionsSetEvent)}, nil
}
