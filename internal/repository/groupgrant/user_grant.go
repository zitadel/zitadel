package groupgrant

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueGroupGrant             = "group_grant"
	groupGrantEventTypePrefix    = eventstore.EventType("group.grant.")
	GroupGrantAddedType          = groupGrantEventTypePrefix + "added"
	GroupGrantChangedType        = groupGrantEventTypePrefix + "changed"
	GroupGrantCascadeChangedType = groupGrantEventTypePrefix + "cascade.changed" // Why --> No Idea as of now
	GroupGrantRemovedType        = groupGrantEventTypePrefix + "removed"
	GroupGrantCascadeRemovedType = groupGrantEventTypePrefix + "cascade.removed" // Why --> No Idea as of now
	GroupGrantDeactivatedType    = groupGrantEventTypePrefix + "deactivated"
	GroupGrantReactivatedType    = groupGrantEventTypePrefix + "reactivated"
)

func NewAddGroupGrantUniqueConstraint(resourceOwner, groupID, projectID, projectGrantID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueGroupGrant,
		fmt.Sprintf("%s:%s:%s:%v", resourceOwner, groupID, projectID, projectGrantID),
		"Errors.GroupGrant.AlreadyExists")
}

func NewRemoveGroupGrantUniqueConstraint(resourceOwner, groupID, projectID, projectGrantID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueGroupGrant,
		fmt.Sprintf("%s:%s:%s:%s", resourceOwner, groupID, projectID, projectGrantID))
}

type GroupGrantAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID        string   `json:"groupId,omitempty"`
	ProjectID      string   `json:"projectId,omitempty"`
	ProjectGrantID string   `json:"grantId,omitempty"`
	RoleKeys       []string `json:"roleKeys,omitempty"`
}

func (e *GroupGrantAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupGrantAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.GroupID, e.ProjectID, e.ProjectGrantID)}
}

func NewGroupGrantAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	projectID,
	projectGrantID string,
	roleKeys []string) *GroupGrantAddedEvent {
	return &GroupGrantAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantAddedType,
		),
		GroupID:        groupID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
		RoleKeys:       roleKeys,
	}
}

func GroupGrantAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupGrantAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GGRANT-1q9ol", "unable to unmarshal group grant")
	}

	return e, nil
}

type GroupGrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys"`
}

func (e *GroupGrantChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupGrantChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupGrantChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	roleKeys []string) *GroupGrantChangedEvent {
	return &GroupGrantChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func GroupGrantChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GGRANT-3N1se", "unable to unmarshal group grant")
	}

	return e, nil
}

type GroupGrantCascadeChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	RoleKeys             []string `json:"roleKeys,omitempty"`
}

func (e *GroupGrantCascadeChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupGrantCascadeChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupGrantCascadeChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	roleKeys []string) *GroupGrantCascadeChangedEvent {
	return &GroupGrantCascadeChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantCascadeChangedType,
		),
		RoleKeys: roleKeys,
	}
}

func GroupGrantCascadeChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupGrantChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GGRANT-Ua1df", "unable to unmarshal group grant")
	}

	return e, nil
}

type GroupGrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	groupID              string `json:"-"`
	projectID            string `json:"-"`
	projectGrantID       string `json:"-"`
}

func (e *GroupGrantRemovedEvent) Payload() interface{} {
	return nil
}

func (e *GroupGrantRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.groupID, e.projectID, e.projectGrantID)}
}

func NewGroupGrantRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	projectID,
	projectGrantID string,
) *GroupGrantRemovedEvent {
	return &GroupGrantRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantRemovedType,
		),
		groupID:        groupID,
		projectID:      projectID,
		projectGrantID: projectGrantID,
	}
}

func GroupGrantRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupGrantCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
	groupID              string `json:"-"`
	projectID            string `json:"-"`
	projectGrantID       string `json:"-"`
}

func (e *GroupGrantCascadeRemovedEvent) Payload() interface{} {
	return nil
}

func (e *GroupGrantCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.groupID, e.projectID, e.projectGrantID)}
}

func NewGroupGrantCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID,
	projectID,
	projectGrantID string,
) *GroupGrantCascadeRemovedEvent {
	return &GroupGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantCascadeRemovedType,
		),
		groupID:        groupID,
		projectID:      projectID,
		projectGrantID: projectGrantID,
	}
}

func GroupGrantCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupGrantCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupGrantDeactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *GroupGrantDeactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupGrantDeactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupGrantDeactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupGrantDeactivatedEvent {
	return &GroupGrantDeactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantDeactivatedType,
		),
	}
}

func GroupGrantDeactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupGrantDeactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupGrantReactivatedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *GroupGrantReactivatedEvent) Payload() interface{} {
	return nil
}

func (e *GroupGrantReactivatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupGrantReactivatedEvent(ctx context.Context, aggregate *eventstore.Aggregate) *GroupGrantReactivatedEvent {
	return &GroupGrantReactivatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupGrantReactivatedType,
		),
	}
}

func GroupGrantReactivatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupGrantReactivatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
