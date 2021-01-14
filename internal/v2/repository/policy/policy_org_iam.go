package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	//TODO: use for org events as suffix (when possible)
	OrgIAMPolicyAddedEventType   = "policy.org.iam.added"
	OrgIAMPolicyChangedEventType = "policy.org.iam.changed"
)

type OrgIAMPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain,omitempty"`
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

	UserLoginMustBeDomain *bool `json:"userLoginMustBeDomain,omitempty"`
}

func (e *OrgIAMPolicyChangedEvent) Data() interface{} {
	return e
}

func NewOrgIAMPolicyChangedEvent(
	base *eventstore.BaseEvent,
) *OrgIAMPolicyChangedEvent {
	return &OrgIAMPolicyChangedEvent{
		BaseEvent: *base,
	}
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

type OrgIAMPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *OrgIAMPolicyRemovedEvent) Data() interface{} {
	return e
}

func NewOrgIAMPolicyRemovedEvent(
	base *eventstore.BaseEvent,
) *OrgIAMPolicyRemovedEvent {
	return &OrgIAMPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func OrgIAMPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &OrgIAMPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-0Pl9d", "unable to unmarshal policy")
	}

	return e, nil
}
