package flow

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	eventTypePrefix            = eventstore.EventType("flow.")
	triggerActionsPrefix       = eventTypePrefix + "trigger_actions."
	TriggerActionsSetEventType = triggerActionsPrefix + "set"
	//ChangedEventType     = triggerActionsPrefix + "changed"
	//DeactivatedEventType = eventTypePrefix + "deactivated"
	//ReactivatedEventType = eventTypePrefix + "reactivated"
	//RemovedEventType     = eventTypePrefix + "removed"
)

//
//func NewAddActionNameUniqueConstraint(actionName, resourceOwner string) *eventstore.EventUniqueConstraint {
//	return eventstore.NewAddEventUniqueConstraint(
//		UniqueActionNameType,
//		actionName+":"+resourceOwner,
//		"Errors.Action.AlreadyExists")
//}
//
//func NewRemoveActionNameUniqueConstraint(actionName, resourceOwner string) *eventstore.EventUniqueConstraint {
//	return eventstore.NewRemoveEventUniqueConstraint(
//		UniqueActionNameType,
//		actionName+":"+resourceOwner)
//}

type TriggerActionsSetEvent struct {
	eventstore.BaseEvent

	FlowType    domain.FlowType
	TriggerType domain.TriggerType
	ActionIDs   map[int32]string
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
	actionIDs map[int32]string,
) *TriggerActionsSetEvent {
	return &TriggerActionsSetEvent{
		BaseEvent:   *base,
		FlowType:    flowType,
		TriggerType: triggerType,
		ActionIDs:   actionIDs,
	}
}

func TriggerActionsSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &TriggerActionsSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "FLOW-4n8vs", "unable to unmarshal trigger actions")
	}

	return e, nil
}

//
//type ChangedEvent struct {
//	eventstore.BaseEvent `json:"-"`
//
//	Name          *string        `json:"name,omitempty"`
//	Script        *string        `json:"script,omitempty"`
//	Timeout       *time.Duration `json:"timeout,omitempty"`
//	AllowedToFail *bool          `json:"allowedToFail,omitempty"`
//	oldName       string
//}
//
//func (e *ChangedEvent) Data() interface{} {
//	return e
//}
//
//func (e *ChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
//	if e.oldName == "" {
//		return nil
//	}
//	return []*eventstore.EventUniqueConstraint{
//		NewRemoveActionNameUniqueConstraint(e.oldName, e.Aggregate().ResourceOwner),
//		NewAddActionNameUniqueConstraint(*e.Name, e.Aggregate().ResourceOwner),
//	}
//}
//
//func NewChangedEvent(
//	ctx context.Context,
//	aggregate *eventstore.Aggregate,
//	changes []ActionChanges,
//) (*ChangedEvent, error) {
//	if len(changes) == 0 {
//		return nil, errors.ThrowPreconditionFailed(nil, "ACTION-dg4t2", "Errors.NoChangesFound")
//	}
//	changeEvent := &ChangedEvent{
//		BaseEvent: *eventstore.NewBaseEventForPush(
//			ctx,
//			aggregate,
//			ChangedEventType,
//		),
//	}
//	for _, change := range changes {
//		change(changeEvent)
//	}
//	return changeEvent, nil
//}
//
//type ActionChanges func(event *ChangedEvent)
//
//func ChangeName(name, oldName string) func(event *ChangedEvent) {
//	return func(e *ChangedEvent) {
//		e.Name = &name
//		e.oldName = oldName
//	}
//}
//
//func ChangeScript(script string) func(event *ChangedEvent) {
//	return func(e *ChangedEvent) {
//		e.Script = &script
//	}
//}
//
//func ChangeTimeout(timeout time.Duration) func(event *ChangedEvent) {
//	return func(e *ChangedEvent) {
//		e.Timeout = &timeout
//	}
//}
//
//func ChangeAllowedToFail(allowedToFail bool) func(event *ChangedEvent) {
//	return func(e *ChangedEvent) {
//		e.AllowedToFail = &allowedToFail
//	}
//}
//
//func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	e := &ChangedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}
//
//	err := json.Unmarshal(event.Data, e)
//	if err != nil {
//		return nil, errors.ThrowInternal(err, "ACTION-4n8vs", "unable to unmarshal action changed")
//	}
//
//	return e, nil
//}
//
//type DeactivatedEvent struct {
//	eventstore.BaseEvent `json:"-"`
//}
//
//func (e *DeactivatedEvent) Data() interface{} {
//	return nil
//}
//
//func (e *DeactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
//	return nil
//}
//
//func NewDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *DeactivatedEvent {
//	return &DeactivatedEvent{
//		BaseEvent: *eventstore.NewBaseEventForPush(
//			ctx,
//			aggregate,
//			DeactivatedEventType,
//		),
//	}
//}
//
//func DeactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	return &DeactivatedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}, nil
//}
//
//type ReactivatedEvent struct {
//	eventstore.BaseEvent `json:"-"`
//}
//
//func (e *ReactivatedEvent) Data() interface{} {
//	return nil
//}
//
//func (e *ReactivatedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
//	return nil
//}
//
//func NewReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *ReactivatedEvent {
//	return &ReactivatedEvent{
//		BaseEvent: *eventstore.NewBaseEventForPush(
//			ctx,
//			aggregate,
//			ReactivatedEventType,
//		),
//	}
//}
//
//func ReactivatedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	return &ReactivatedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}, nil
//}
//
//type RemovedEvent struct {
//	eventstore.BaseEvent `json:"-"`
//
//	name string
//}
//
//func (e *RemovedEvent) Data() interface{} {
//	return e
//}
//
//func (e *RemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
//	return []*eventstore.EventUniqueConstraint{NewRemoveActionNameUniqueConstraint(e.name, e.Aggregate().ResourceOwner)}
//}
//
//func NewRemovedEvent(
//	ctx context.Context,
//	aggregate *eventstore.Aggregate,
//	name string,
//) *RemovedEvent {
//	return &RemovedEvent{
//		BaseEvent: *eventstore.NewBaseEventForPush(
//			ctx,
//			aggregate,
//			RemovedEventType,
//		),
//		name: name,
//	}
//}
//
//func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
//	return &RemovedEvent{
//		BaseEvent: *eventstore.BaseEventFromRepo(event),
//	}, nil
//}
