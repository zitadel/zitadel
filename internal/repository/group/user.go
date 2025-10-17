package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	GroupUsersAddedEventType   = groupEventTypePrefix + "users.added"
	GroupUsersRemovedEventType = groupEventTypePrefix + "users.removed"
)

type GroupUsersAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserIDs []string `json:"userId"`
}

func NewGroupUsersAddedEvent(ctx context.Context, aggregate *eventstore.Aggregate, userIDs []string) *GroupUsersAddedEvent {
	return &GroupUsersAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUsersAddedEventType,
		),
		UserIDs: userIDs,
	}
}

func (e *GroupUsersAddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUsersAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupUsersAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

type GroupUsersRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserIDs []string `json:"userId"`
}

func NewGroupUsersRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userIDs []string,
) *GroupUsersRemovedEvent {
	return &GroupUsersRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupUsersRemovedEventType,
		),
		UserIDs: userIDs,
	}
}

func (e *GroupUsersRemovedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = *event
}

func (e *GroupUsersRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupUsersRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}
