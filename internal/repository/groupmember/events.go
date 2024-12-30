package member

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueMember            = "groupmember"
	AddedEventType          = "groupmember.added"
	ChangedEventType        = "groupmember.changed"
	RemovedEventType        = "groupmember.removed"
	CascadeRemovedEventType = "groupmember.cascade.removed"
)

func NewAddGroupMemberUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, userID),
		"Errors.GroupMember.AlreadyExists")
}

func NewRemoveGroupMemberUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, userID),
	)
}

type GroupMemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *GroupMemberAddedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddGroupMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewGroupMemberAddedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *GroupMemberAddedEvent {

	return &GroupMemberAddedEvent{
		BaseEvent: *base,
		UserID:    userID,
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

	UserID string `json:"userId,omitempty"`
}

func (e *GroupMemberChangedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewGroupMemberChangedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *GroupMemberChangedEvent {
	return &GroupMemberChangedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func ChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
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

	UserID string `json:"userId"`
}

func (e *GroupMemberRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *GroupMemberRemovedEvent {

	return &GroupMemberRemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
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

	UserID string `json:"userId"`
}

func (e *GroupMemberCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *GroupMemberCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveGroupMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewCascadeRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *GroupMemberCascadeRemovedEvent {

	return &GroupMemberCascadeRemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func CascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &GroupMemberCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MEMBER-3j9sf", "unable to unmarshal label policy")
	}

	return e, nil
}
