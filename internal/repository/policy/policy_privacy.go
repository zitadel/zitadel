package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	PrivacyPolicyAddedEventType   = "policy.privacy.added"
	PrivacyPolicyChangedEventType = "policy.privacy.changed"
	PrivacyPolicyRemovedEventType = "policy.privacy.removed"
)

type PrivacyPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TOSLink     string `json:"tosLink,omitempty"`
	PrivacyLink string `json:"privacyLink,omitempty"`
}

func (e *PrivacyPolicyAddedEvent) Data() interface{} {
	return e
}

func (e *PrivacyPolicyAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPrivacyPolicyAddedEvent(
	base *eventstore.BaseEvent,
	tosLink,
	privacyLink string,
) *PrivacyPolicyAddedEvent {
	return &PrivacyPolicyAddedEvent{
		BaseEvent:   *base,
		TOSLink:     tosLink,
		PrivacyLink: privacyLink,
	}
}

func PrivacyPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PrivacyPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-2k0fs", "unable to unmarshal policy")
	}

	return e, nil
}

type PrivacyPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TOSLink     *string `json:"tosLink,omitempty"`
	PrivacyLink *string `json:"privacyLink,omitempty"`
}

func (e *PrivacyPolicyChangedEvent) Data() interface{} {
	return e
}

func (e *PrivacyPolicyChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPrivacyPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PrivacyPolicyChanges,
) (*PrivacyPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-PPo0s", "Errors.NoChangesFound")
	}
	changeEvent := &PrivacyPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type PrivacyPolicyChanges func(*PrivacyPolicyChangedEvent)

func ChangeTOSLink(tosLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.TOSLink = &tosLink
	}
}

func ChangePrivacyLink(privacyLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.PrivacyLink = &privacyLink
	}
}

func PrivacyPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PrivacyPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-22nf9", "unable to unmarshal policy")
	}

	return e, nil
}

type PrivacyPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PrivacyPolicyRemovedEvent) Data() interface{} {
	return nil
}

func (e *PrivacyPolicyRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewPrivacyPolicyRemovedEvent(base *eventstore.BaseEvent) *PrivacyPolicyRemovedEvent {
	return &PrivacyPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PrivacyPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PrivacyPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
