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

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Roles  []string `json:"roles,omitempty"`
	UserID string   `json:"userId,omitempty"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewMemberChangedEvent(
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

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-puqv4", "unable to unmarshal label policy")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserID string `json:"userId"`
}

func (e *RemovedEvent) Data() interface{} {
	return e
}

func NewRemovedEvent(
	base *eventstore.BaseEvent,
	userID string,
) *RemovedEvent {

	return &RemovedEvent{
		BaseEvent: *base,
		UserID:    userID,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-Ep4ip", "unable to unmarshal label policy")
	}

	return e, nil
}
