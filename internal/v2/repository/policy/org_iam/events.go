package org_iam

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	OrgIAMPolicyAddedEventType   = "policy.org.iam.added"
	OrgIAMPolicyChangedEventType = "policy.org.iam.changed"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	userLoginMustBeDomain bool,
) *AddedEvent {

	return &AddedEvent{
		BaseEvent:             *base,
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-TvSmA", "unable to unmarshal policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain"`
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	userLoginMustBeDomain bool,
) *ChangedEvent {
	e := &ChangedEvent{
		BaseEvent: *base,
	}
	if current.UserLoginMustBeDomain != userLoginMustBeDomain {
		e.UserLoginMustBeDomain = userLoginMustBeDomain
	}
	return e
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-0Pl9d", "unable to unmarshal policy")
	}

	return e, nil
}
