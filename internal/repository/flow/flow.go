package flow

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	eventTypePrefix                       = eventstore.EventType("flow.")
	triggerActionsPrefix                  = eventTypePrefix + "trigger_actions."
	TriggerActionsSetEventType            = triggerActionsPrefix + "set"
	TriggerActionsCascadeRemovedEventType = triggerActionsPrefix + "cascade.removed"
	FlowClearedEventType                  = eventTypePrefix + "cleared"
)

type TriggerActionsSetEvent struct {
	eventstore.BaseEvent

	FlowType    domain.FlowType    `json:"flowType"`
	TriggerType domain.TriggerType `json:"triggerType"`
	ActionIDs   []string           `json:"actionIDs"`
}

func (e *TriggerActionsSetEvent) Data() interface{} {
	return e
}

func (e *TriggerActionsSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewTriggerActionsSetEvent(
	base *eventstore.BaseEvent,
	flowType domain.FlowType,
	triggerType domain.TriggerType,
	actionIDs []string,
) *TriggerActionsSetEvent {
	return &TriggerActionsSetEvent{
		BaseEvent:   *base,
		FlowType:    flowType,
		TriggerType: triggerType,
		ActionIDs:   actionIDs,
	}
}

func TriggerActionsSetEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &TriggerActionsSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "FLOW-4n8vs", "unable to unmarshal trigger actions")
	}

	return e, nil
}

type TriggerActionsCascadeRemovedEvent struct {
	eventstore.BaseEvent

	FlowType    domain.FlowType    `json:"flowType"`
	TriggerType domain.TriggerType `json:"triggerType"`
	ActionID    string             `json:"actionID"`
}

func (e *TriggerActionsCascadeRemovedEvent) Data() interface{} {
	return e
}

func (e *TriggerActionsCascadeRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewTriggerActionsCascadeRemovedEvent(
	base *eventstore.BaseEvent,
	flowType domain.FlowType,
	actionID string,
) *TriggerActionsCascadeRemovedEvent {
	return &TriggerActionsCascadeRemovedEvent{
		BaseEvent: *base,
		FlowType:  flowType,
		ActionID:  actionID,
	}
}

func TriggerActionsCascadeRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &TriggerActionsCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "FLOW-4n8vs", "unable to unmarshal trigger actions")
	}

	return e, nil
}

type FlowClearedEvent struct {
	eventstore.BaseEvent

	FlowType domain.FlowType `json:"flowType"`
}

func (e *FlowClearedEvent) Data() interface{} {
	return e
}

func (e *FlowClearedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewFlowClearedEvent(
	base *eventstore.BaseEvent,
	flowType domain.FlowType,
) *FlowClearedEvent {
	return &FlowClearedEvent{
		BaseEvent: *base,
		FlowType:  flowType,
	}
}

func FlowClearedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &FlowClearedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "FLOW-BHfg2", "unable to unmarshal flow cleared")
	}

	return e, nil
}
