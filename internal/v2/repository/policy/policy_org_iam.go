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

func (e *OrgIAMPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
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

func (e *OrgIAMPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOrgIAMPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []OrgIAMPolicyChanges,
) (*OrgIAMPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-DAf3h", "Errors.NoChangesFound")
	}
	changeEvent := &OrgIAMPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OrgIAMPolicyChanges func(*OrgIAMPolicyChangedEvent)

func ChangeUserLoginMustBeDomain(userLoginMustBeDomain bool) func(*OrgIAMPolicyChangedEvent) {
	return func(e *OrgIAMPolicyChangedEvent) {
		e.UserLoginMustBeDomain = &userLoginMustBeDomain
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
	return nil
}

func (e *OrgIAMPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewOrgIAMPolicyRemovedEvent(base *eventstore.BaseEvent) *OrgIAMPolicyRemovedEvent {
	return &OrgIAMPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func OrgIAMPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &OrgIAMPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
