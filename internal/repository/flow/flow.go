package flow

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	eventTypePrefix                       = eventstore.EventType("flow.")
	triggerActionsPrefix                  = eventTypePrefix + "trigger_actions."
	TriggerActionsSetEventType            = triggerActionsPrefix + "set"
	TriggerActionsCascadeRemovedEventType = triggerActionsPrefix + "cascade.removed"
	FlowClearedEventType                  = eventTypePrefix + "cleared"
)

type TriggerActionsSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	FlowType    domain.FlowType    `json:"flowType"`
	TriggerType domain.TriggerType `json:"triggerType"`
	ActionIDs   []string           `json:"actionIDs"`
}

func (e *TriggerActionsSetEvent) Payload() interface{} {
	return e
}

func (e *TriggerActionsSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func TriggerActionsSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &TriggerActionsSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "FLOW-4n8vs", "unable to unmarshal trigger actions")
	}

	return e, nil
}

type TriggerActionsCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FlowType    domain.FlowType    `json:"flowType"`
	TriggerType domain.TriggerType `json:"triggerType"`
	ActionID    string             `json:"actionID"`
}

func (e *TriggerActionsCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *TriggerActionsCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func TriggerActionsCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &TriggerActionsCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "FLOW-4n8vs", "unable to unmarshal trigger actions")
	}

	return e, nil
}

type FlowClearedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FlowType domain.FlowType `json:"flowType"`
}

func (e *FlowClearedEvent) Payload() interface{} {
	return e
}

func (e *FlowClearedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func FlowClearedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &FlowClearedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "FLOW-BHfg2", "unable to unmarshal flow cleared")
	}

	return e, nil
}
