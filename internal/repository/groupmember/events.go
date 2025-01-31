package groupmember

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	GroupUniqueMember            = "member.group"
	GroupAddedEventType          = "member.group.added"
	GroupChangedEventType        = "member.group.changed"
	GroupRemovedEventType        = "member.group.removed"
	GroupCascadeRemovedEventType = "member.group.cascade.removed"
)

func NewAddGroupMemberUniqueConstraint(aggregateID, groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		GroupUniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, groupID),
		"Errors.GroupMember.AlreadyExists")
}

func NewRemoveGroupMemberUniqueConstraint(aggregateID, groupID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		GroupUniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, groupID),
	)
}

type GroupMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string   `json:"groupId"`
	Roles   []string `json:"roles"`
}

func (e *GroupMemberAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID)}
}

func NewGroupMemberAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	// base *eventstore.BaseEvent,
	groupID string,
	roles ...string,
) *GroupMemberAddedEvent {
	return &GroupMemberAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupAddedEventType,
		),
		GroupID: groupID,
		Roles:   roles,
	}
}

func GroupMemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupMemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUPM-qwqv4", "unable to unmarshal group member")
	}

	return e, nil
}

type GroupMemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string   `json:"groupId,omitempty"`
	Roles   []string `json:"roles,omitempty"`
}

func (e *GroupMemberChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupMemberChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
	roles ...string,
) *GroupMemberChangedEvent {
	return &GroupMemberChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupChangedEventType,
		),
		GroupID: groupID,
		Roles:   roles,
	}
}

func GroupChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupMemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUPM-qzqv4", "unable to unmarshal group member")
	}

	return e, nil
}

type GroupMemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string `json:"groupId"`
}

func (e *GroupMemberRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID)}
}

func NewGroupRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
) *GroupMemberRemovedEvent {

	return &GroupMemberRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupRemovedEventType,
		),
		GroupID: groupID,
	}
}

func GroupRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupMemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "GROUPM-Fp4ip", "unable to unmarshal group member")
	}

	return e, nil
}

type GroupMemberCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	GroupID string `json:"groupId"`
}

func (e *GroupMemberCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupMemberUniqueConstraint(e.Aggregate().ID, e.GroupID)}
}

func NewGroupCascadeRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	groupID string,
) *GroupMemberCascadeRemovedEvent {

	return &GroupMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			GroupCascadeRemovedEventType,
		),
		GroupID: groupID,
	}
}

func GroupCascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MEMBER-3j9sf", "unable to unmarshal label policy")
	}

	return e, nil
}
