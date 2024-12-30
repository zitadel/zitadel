package group

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	group_member "github.com/zitadel/zitadel/internal/repository/groupmember"
	"github.com/zitadel/zitadel/internal/repository/member"
)

var (
	MemberAddedType          = groupEventTypePrefix + group_member.AddedEventType
	MemberChangedType        = groupEventTypePrefix + group_member.ChangedEventType
	MemberRemovedType        = groupEventTypePrefix + group_member.RemovedEventType
	MemberCascadeRemovedType = groupEventTypePrefix + group_member.CascadeRemovedEventType
)

type GroupMemberAddedEvent struct {
	group_member.GroupMemberAddedEvent
}

func NewGroupMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupMemberAddedEvent {
	return &GroupMemberAddedEvent{
		GroupMemberAddedEvent: *group_member.NewGroupMemberAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberAddedType,
			),
			userID,
		),
	}
}

func MemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.MemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberAddedEvent{GroupMemberAddedEvent: *e.(*group_member.GroupMemberAddedEvent)}, nil
}

type GroupMemberChangedEvent struct {
	group_member.GroupMemberChangedEvent
}

func NewGroupMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupMemberChangedEvent {
	return &GroupMemberChangedEvent{
		GroupMemberChangedEvent: *group_member.NewGroupMemberChangedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberChangedType,
			),
			userID,
		),
	}
}

func MemberChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberChangedEvent{GroupMemberChangedEvent: *e.(*group_member.GroupMemberChangedEvent)}, nil
}

type GroupMemberRemovedEvent struct {
	group_member.GroupMemberRemovedEvent
}

func NewGroupMemberRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupMemberRemovedEvent {

	return &GroupMemberRemovedEvent{
		GroupMemberRemovedEvent: *group_member.NewRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberRemovedType,
			),
			userID,
		),
	}
}

func MemberRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := group_member.RemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberRemovedEvent{GroupMemberRemovedEvent: *e.(*group_member.GroupMemberRemovedEvent)}, nil
}

type GroupMemberCascadeRemovedEvent struct {
	group_member.GroupMemberCascadeRemovedEvent
}

func NewGroupMemberCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userID string,
) *GroupMemberCascadeRemovedEvent {

	return &GroupMemberCascadeRemovedEvent{
		GroupMemberCascadeRemovedEvent: *group_member.NewCascadeRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				MemberCascadeRemovedType,
			),
			userID,
		),
	}
}

func MemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := member.CascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberCascadeRemovedEvent{GroupMemberCascadeRemovedEvent: *e.(*group_member.GroupMemberCascadeRemovedEvent)}, nil
}
