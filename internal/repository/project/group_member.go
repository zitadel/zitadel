package project

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/groupmember"
)

var (
	GroupMemberAddedType          = projectEventTypePrefix + groupmember.GroupAddedEventType
	GroupMemberChangedType        = projectEventTypePrefix + groupmember.GroupChangedEventType
	GroupMemberRemovedType        = projectEventTypePrefix + groupmember.GroupRemovedEventType
	GroupMemberCascadeRemovedType = projectEventTypePrefix + groupmember.GroupCascadeRemovedEventType
)

type GroupMemberAddedEvent struct {
	groupmember.GroupMemberAddedEvent
}

func NewProjectGroupMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
	roles ...string,
) *GroupMemberAddedEvent {
	return &GroupMemberAddedEvent{
		GroupMemberAddedEvent: *groupmember.NewGroupMemberAddedEvent(
			ctx,
			aggregate,
			groupID,
			roles...,
		),
	}
}

func GroupMemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := groupmember.GroupMemberAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberAddedEvent{GroupMemberAddedEvent: *e.(*groupmember.GroupMemberAddedEvent)}, nil
}

type GroupMemberChangedEvent struct {
	groupmember.GroupMemberChangedEvent
}

func GroupNewProjectGroupMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
	roles ...string,
) *GroupMemberChangedEvent {

	return &GroupMemberChangedEvent{
		GroupMemberChangedEvent: *groupmember.NewGroupMemberChangedEvent(
			ctx,
			aggregate,
			groupID,
			roles...,
		),
	}
}

func GroupGroupMemberChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := groupmember.GroupChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberChangedEvent{GroupMemberChangedEvent: *e.(*groupmember.GroupMemberChangedEvent)}, nil
}

type GroupMemberRemovedEvent struct {
	groupmember.GroupMemberRemovedEvent
}

func NewProjectGroupMemberRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
) *GroupMemberRemovedEvent {

	return &GroupMemberRemovedEvent{
		GroupMemberRemovedEvent: *groupmember.NewGroupRemovedEvent(
			ctx,
			aggregate,
			groupID,
		),
	}
}

func GroupMemberRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := groupmember.GroupRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberRemovedEvent{GroupMemberRemovedEvent: *e.(*groupmember.GroupMemberRemovedEvent)}, nil
}

type GroupMemberCascadeRemovedEvent struct {
	groupmember.GroupMemberCascadeRemovedEvent
}

func NewProjectGroupMemberCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
) *GroupMemberCascadeRemovedEvent {

	return &GroupMemberCascadeRemovedEvent{
		GroupMemberCascadeRemovedEvent: *groupmember.NewGroupCascadeRemovedEvent(
			ctx,
			aggregate,
			groupID,
		),
	}
}

func GroupMemberCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := groupmember.GroupCascadeRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &GroupMemberCascadeRemovedEvent{GroupMemberCascadeRemovedEvent: *e.(*groupmember.GroupMemberCascadeRemovedEvent)}, nil
}
