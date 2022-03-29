package policy

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	//TODO: use for org events as suffix (when possible)
	DomainPolicyAddedEventType   = "policy.domain.added"
	DomainPolicyChangedEventType = "policy.domain.changed"
	DomainPolicyRemovedEventType = "policy.domain.removed"
)

type DomainPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain bool `json:"userLoginMustBeDomain,omitempty"`
}

func (e *DomainPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *DomainPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDomainPolicyAddedEvent(
	base *eventstore.BaseEvent,
	userLoginMustBeDomain bool,
) *DomainPolicyAddedEvent {

	return &DomainPolicyAddedEvent{
		BaseEvent:             *base,
		UserLoginMustBeDomain: userLoginMustBeDomain,
	}
}

func DomainPolicyAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DomainPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-TvSmA", "unable to unmarshal policy")
	}

	return e, nil
}

type DomainPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	UserLoginMustBeDomain *bool `json:"userLoginMustBeDomain,omitempty"`
}

func (e *DomainPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *DomainPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDomainPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []OrgPolicyChanges,
) (*DomainPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-DAf3h", "Errors.NoChangesFound")
	}
	changeEvent := &DomainPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type OrgPolicyChanges func(*DomainPolicyChangedEvent)

func ChangeUserLoginMustBeDomain(userLoginMustBeDomain bool) func(*DomainPolicyChangedEvent) {
	return func(e *DomainPolicyChangedEvent) {
		e.UserLoginMustBeDomain = &userLoginMustBeDomain
	}
}

func DomainPolicyChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DomainPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-0Pl9d", "unable to unmarshal policy")
	}

	return e, nil
}

type DomainPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *DomainPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *DomainPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDomainPolicyRemovedEvent(base *eventstore.BaseEvent) *DomainPolicyRemovedEvent {
	return &DomainPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func DomainPolicyRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	return &DomainPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
