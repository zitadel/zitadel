package member

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueMember            = "member"
	AddedEventType          = "member.added"
	ChangedEventType        = "member.changed"
	RemovedEventType        = "member.removed"
	CascadeRemovedEventType = "member.cascade.removed"
)

func NewAddMemberUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, userID),
		"Errors.Member.AlreadyExists")
}

func NewRemoveMemberUniqueConstraint(aggregateID, userID string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueMember,
		fmt.Sprintf("%s:%s", aggregateID, userID),
	)
}

type MemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *MemberAddedEvent) Payload() interface{} {
	return e
}

func (e *MemberAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewMemberAddedEvent(
	base *eventstore.BaseEvent,
	userID string,
	roles ...string,
) *MemberAddedEvent {

	return &MemberAddedEvent{
		BaseEvent: *base,
		Roles:     roles,
		UserID:    userID,
	}
}

func MemberAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type MemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles,omitempty"`
	UserID string   `json:"userId,omitempty"`
}

func (e *MemberChangedEvent) Payload() interface{} {
	return e
}

func (e *MemberChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMemberChangedEvent(
	base *eventstore.BaseEvent,
	userID string,
	roles ...string,
) *MemberChangedEvent {
	return &MemberChangedEvent{
		BaseEvent: *base,
		Roles:     roles,
		UserID:    userID,
	}
}

func ChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type MemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *MemberRemovedEvent) Payload() interface{} {
	return e
}

func (e *MemberRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *MemberRemovedEvent {

	return &MemberRemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func RemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MEMBER-Ep4ip", "unable to unmarshal label policy")
	}

	return e, nil
}

type MemberCascadeRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *MemberCascadeRemovedEvent) Payload() interface{} {
	return e
}

func (e *MemberCascadeRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveMemberUniqueConstraint(e.Aggregate().ID, e.UserID)}
}

func NewCascadeRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *MemberCascadeRemovedEvent {

	return &MemberCascadeRemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func CascadeRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MemberCascadeRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "MEMBER-3j9sf", "unable to unmarshal label policy")
	}

	return e, nil
}
