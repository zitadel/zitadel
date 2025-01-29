package groupmember

import (
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
	base *eventstore.BaseEvent,
	groupID string,
	roles ...string,
) *GroupMemberAddedEvent {

	return &GroupMemberAddedEvent{
		BaseEvent: *base,
		GroupID:   groupID,
		Roles:     roles,
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
	base *eventstore.BaseEvent,
	groupID string,
	roles ...string,
) *GroupMemberChangedEvent {
	return &GroupMemberChangedEvent{
		BaseEvent: *base,
		GroupID:   groupID,
		Roles:     roles,
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
	base *eventstore.BaseEvent,
	groupID string,
) *GroupMemberRemovedEvent {

	return &GroupMemberRemovedEvent{
		BaseEvent: *base,
		GroupID:   groupID,
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
	base *eventstore.BaseEvent,
	groupID string,
) *GroupMemberCascadeRemovedEvent {

	return &GroupMemberCascadeRemovedEvent{
		BaseEvent: *base,
		GroupID:   groupID,
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
