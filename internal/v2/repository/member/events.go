package member

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	AddedEventType   = "member.added"
	ChangedEventType = "member.changed"
	RemovedEventType = "member.removed"
)

type MemberAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles"`
	UserID string   `json:"userId"`
}

func (e *MemberAddedEvent) Data() interface{} {
	return e
}

func (e *MemberAddedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
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

func MemberAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MemberAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type MemberChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles,omitempty"`
	UserID string   `json:"userId,omitempty"`
}

func (e *MemberChangedEvent) Data() interface{} {
	return e
}

func (e *MemberChangedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
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

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MemberChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type MemberRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *MemberRemovedEvent) Data() interface{} {
	return e
}

func (e *MemberRemovedEvent) UniqueConstraint() []eventstore.EventUniqueConstraint {
	return nil
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

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MemberRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-Ep4ip", "unable to unmarshal label policy")
	}

	return e, nil
}
