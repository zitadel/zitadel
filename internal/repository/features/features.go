package features

import (
	"encoding/json"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	featuresPrefix           = "features."
	FeaturesSetEventType     = featuresPrefix + "set"
	FeaturesRemovedEventType = featuresPrefix + "removed"
)

type FeaturesSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	TierName                 *string
	TierDescription          *string
	TierState                *domain.TierState
	TierStateDescription     *string
	LoginPolicyFactors       *bool
	LoginPolicyIDP           *bool
	LoginPolicyPasswordless  *bool
	LoginPolicyRegistration  *bool
	LoginPolicyUsernameLogin *bool
}

func (e *FeaturesSetEvent) Data() interface{} {
	return e
}

func (e *FeaturesSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewFeaturesSetEvent(
	base *eventstore.BaseEvent,
	changes []FeaturesChanges,
) (*FeaturesSetEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "FEATURES-d34F4", "Errors.NoChangesFound")
	}
	changeEvent := &FeaturesSetEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type FeaturesChanges func(*FeaturesSetEvent)

func ChangeTierName(tierName string) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.TierName = &tierName
	}
}

func ChangeTierDescription(tierDescription string) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.TierDescription = &tierDescription
	}
}

func ChangeTierState(tierState domain.TierState) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.TierState = &tierState
	}
}

func ChangeTierStateDescription(statusDescription string) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.TierStateDescription = &statusDescription
	}
}

func ChangeLoginPolicyFactors(loginPolicyFactors bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.LoginPolicyFactors = &loginPolicyFactors
	}
}

func ChangeLoginPolicyIDP(loginPolicyIDP bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.LoginPolicyIDP = &loginPolicyIDP
	}
}

func ChangeLoginPolicyPasswordless(loginPolicyPasswordless bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.LoginPolicyPasswordless = &loginPolicyPasswordless
	}
}

func ChangeLoginPolicyRegistration(loginPolicyRegistration bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.LoginPolicyRegistration = &loginPolicyRegistration
	}
}

func ChangeLoginPolicyUsernameLogin(loginPolicyUsernameLogin bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.LoginPolicyUsernameLogin = &loginPolicyUsernameLogin
	}
}

func FeaturesSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &FeaturesSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "FEATURES-fdgDg", "unable to unmarshal features")
	}

	return e, nil
}

type FeaturesRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *FeaturesRemovedEvent) Data() interface{} {
	return nil
}

func (e *FeaturesRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewFeaturesRemovedEvent(base *eventstore.BaseEvent) *FeaturesRemovedEvent {
	return &FeaturesRemovedEvent{
		BaseEvent: *base,
	}
}

func FeaturesRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &FeaturesRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
