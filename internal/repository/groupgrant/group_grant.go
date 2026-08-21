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
	GroupGrantRemovedType        = groupGrantEventTypePrefix + "removed"
	GroupGrantCascadeRemovedType = groupGrantEventTypePrefix + "cascade.removed"
)

func NewAddGroupGrantUniqueConstraint(resourceOwner, groupID, projectID, projectGrantID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueGroupGrant,
		fmt.Sprintf("%s:%s:%s:%s", resourceOwner, groupID, projectID, projectGrantID),
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
	roleKeys []string,
) *GroupGrantAddedEvent {
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
		return nil, zerrors.ThrowInternal(err, "GGRANT-jx5Tlk", "unable to unmarshal group grant")
	}

	return e, nil
}

type GroupGrantChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	RoleKeys []string `json:"roleKeys"`
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
	roleKeys []string,
) *GroupGrantChangedEvent {
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
		return nil, zerrors.ThrowInternal(err, "GGRANT-7pWqVd", "unable to unmarshal group grant")
	}

	return e, nil
}

type GroupGrantRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID        string `json:"groupId,omitempty"`
	ProjectID      string `json:"projectId,omitempty"`
	ProjectGrantID string `json:"grantId,omitempty"`
}

func (e *GroupGrantRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupGrantRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupGrantUniqueConstraint(e.Aggregate().ResourceOwner, e.GroupID, e.ProjectID, e.ProjectGrantID)}
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
		GroupID:        groupID,
		ProjectID:      projectID,
		ProjectGrantID: projectGrantID,
	}
}

func GroupGrantRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &GroupGrantRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type GroupGrantCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	groupID        string `json:"-"`
	projectID      string `json:"-"`
	projectGrantID string `json:"-"`
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
