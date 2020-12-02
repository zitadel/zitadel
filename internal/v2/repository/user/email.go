package user

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	emailEventPrefix                 = humanEventPrefix + "email."
	HumanEmailChangedType            = emailEventPrefix + "changed"
	HumanEmailVerifiedType           = emailEventPrefix + "verified"
	HumanEmailVerificationFailedType = emailEventPrefix + "verification.failed"
	HumanEmailCodeAddedType          = emailEventPrefix + "code.added"
	HumanEmailCodeSentType           = emailEventPrefix + "code.sent"
)

type HumanEmailChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	EmailAddress string `json:"email,omitempty"`
}

func (e *HumanEmailChangedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanEmailChangedEvent) Data() interface{} {
	return e
}

func NewHumanHumanEmailChangedEvent(base *eventstore.BaseEvent, emailAddress string) *HumanEmailChangedEvent {
	return &HumanEmailChangedEvent{
		BaseEvent:    *base,
		EmailAddress: emailAddress,
	}
}

func HumanEmailChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailChangedEvent := &HumanEmailChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, emailChangedEvent)
	if err != nil {
		return nil, errors.ThrowInternal(err, "USER-4M0sd", "unable to unmarshal human password changed")
	}

	return emailChangedEvent, nil
}

type HumanEmailVerifiedEvent struct {
	eventstore.BaseEvent `json:"-"`

	IsEmailVerified bool `json:"-"`
}

func (e *HumanEmailVerifiedEvent) CheckPrevious() bool {
	return false
}

func (e *HumanEmailVerifiedEvent) Data() interface{} {
	return nil
}

func NewHumanEmailVerifiedEvent(base *eventstore.BaseEvent) *HumanEmailVerifiedEvent {
	return &HumanEmailVerifiedEvent{
		BaseEvent: *base,
	}
}

func HumanEmailVerifiedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	emailVerified := &HumanEmailVerifiedEvent{
		BaseEvent:       *eventstore.BaseEventFromRepo(event),
		IsEmailVerified: true,
	}
	return emailVerified, nil
}
