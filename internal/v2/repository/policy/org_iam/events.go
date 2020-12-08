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

type OrgIAMPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain"`
}

func (e *OrgIAMPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *OrgIAMPolicyAddedEvent) Data() interface{} {
	return e
}

func NewOrgIAMPolicyAddedEvent(
	base *eventstore.BaseEvent,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {

	return &OrgIAMPolicyAddedEvent{
		BaseEvent:             *base,
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func OrgIAMPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OrgIAMPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-TvSmA", "unable to unmarshal policy")
	}

	return e, nil
}

type OrgIAMPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain"`
}

func (e *OrgIAMPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *OrgIAMPolicyChangedEvent) Data() interface{} {
	return e
}

func NewOrgIAMPolicyChangedEvent(
	base *eventstore.BaseEvent,
	current *OrgIAMPolicyWriteModel,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyChangedEvent {
	e := &OrgIAMPolicyChangedEvent{
		BaseEvent: *base,
	}
	if current.UserLoginMustBeDomain != userLoginMustBeDomain {
		e.UserLoginMustBeDomain = userLoginMustBeDomain
	}
	return e
}

func OrgIAMPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OrgIAMPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-0Pl9d", "unable to unmarshal policy")
	}

	return e, nil
}
