package group

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	uniqueGroupUserType       = "group_user"
	GroupUserAddedEventType   = groupEventTypePrefix + "user.added"
	GroupUserRemovedEventType = groupEventTypePrefix + "user.removed"
)

func NewAddGroupUserUniqueConstraint(groupID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		uniqueGroupUserType,
		fmt.Sprintf("%s:%s", groupID, userID),
		"Errors.Group.User.AlreadyExists",
	)
}

func NewRemoveGroupUserUniqueConstraint(groupID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		uniqueGroupUserType,
		fmt.Sprintf("%s:%s", groupID, userID))
}

type GroupUserAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func NewGroupUserAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, userID string) *GroupUserAddedEvent {
	return &GroupUserAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUserAddedEventType,
		),
		UserID: userID,
	}
}

func (e *GroupUserAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUserAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupUserUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

type GroupUserRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func NewGroupUserRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupUserRemovedEvent {
	return &GroupUserRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUserRemovedEventType,
		),
		UserID: userID,
	}
}

func (e *GroupUserRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUserRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupUserRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupUserUniqueConstraint(e.Aggregate().ID, e.UserID)}
}
