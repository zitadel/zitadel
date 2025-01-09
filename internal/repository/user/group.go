package user

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	UniqueUserGroupType = "group_id"

	// Group Events
	UserGroupAddedType   = userEventTypePrefix + "groups.added"
	UserGroupRemovedType = userEventTypePrefix + "groups.removed"
	UserGroupChangedType = userEventTypePrefix + "groups.changed"
)

func NewUserGroupUniqueConstraint(groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUserGroupType,
		groupID,
		"Errors.Group.AlreadyExists")
}

func NewRemoveUserGroupUniqueConstraint(groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUserGroupType,
		groupID)
}

type UserGroupAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	GroupID              []string
}

func (e *UserGroupAddedEvent) Payload() interface{} {
	return e
}

func (e *UserGroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewUserGroupUniqueConstraint(e.Aggregate().ID)}
}

func NewUserGroupAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID ...string,
) *UserGroupAddedEvent {
	return &UserGroupAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGroupAddedType,
		),
		GroupID: groupID,
	}
}

type UserGroupChangedEvent struct {
	eventstore.BaseEvent `json:"-"`
	GroupID              []string `json:"groupID,omitempty"`
}

type UserGroupChanges func(event *UserGroupChangedEvent)

type UserGroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID []string
}

func (e *UserGroupRemovedEvent) Payload() interface{} {
	return e
}

func (e *UserGroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserGroupUniqueConstraint(e.Aggregate().ID)}
}

func NewUserGroupChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID ...string,
) *UserGroupChangedEvent {
	return &UserGroupChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGroupChangedType,
		),
		GroupID: groupID,
	}
}

func NewUserGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID ...string,
) *UserGroupRemovedEvent {
	return &UserGroupRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			UserGroupRemovedType,
		),
		GroupID: groupID,
	}
}
