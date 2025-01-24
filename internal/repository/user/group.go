package user

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueUserGroupType = "user_group_id"

	// Group Events
	UserGroupAddedType          = userEventTypePrefix + "groups.added"
	UserGroupRemovedType        = userEventTypePrefix + "groups.removed"
	UserGroupCascadeRemovedType = userEventTypePrefix + "groups.cascade.removed"
)

func NewUserGroupUniqueConstraint(aggregateID, groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueUserGroupType,
		fmt.Sprintf("%s:%s", aggregateID, groupID),
		"Errors.Group.Member.AlreadyExists")
}

func NewRemoveUserGroupUniqueConstraint(aggregateID, groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueUserGroupType,
		fmt.Sprintf("%s:%s", aggregateID, groupID))
}

type UserGroupAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
	GroupID              string
}

func (e *UserGroupAddedEvent) Payload() interface{} {
	return e
}

func (e *UserGroupAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewUserGroupUniqueConstraint(e.Aggregate().ID, e.GroupID)}
}

func NewUserGroupAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
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

func UserGroupAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserGroupRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-mL0vs", "unable to unmarshal user group member")
	}

	return e, nil
}

type UserGroupRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string
}

func (e *UserGroupRemovedEvent) Payload() interface{} {
	return e
}

func (e *UserGroupRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveUserGroupUniqueConstraint(e.Aggregate().ID, e.GroupID)}
}

func NewUserGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
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

func UserGroupRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &UserGroupAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUP-mL0vs", "unable to unmarshal user group member")
	}

	return e, nil
}

/*
	type UserGroupChangedEvent struct {
		eventstore.BaseEvent `json:"-"`
		GroupID              []string `json:"groupID,omitempty"`
	}

type UserGroupChanges func(event *UserGroupChangedEvent)

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
*/
