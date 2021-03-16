package features

import (
	"encoding/json"
	"time"

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

	TierName                 *string               `json:"tier_name,omitempty"`
	TierDescription          *string               `json:"tier_description,omitempty"`
	State                    *domain.FeaturesState `json:"state,omitempty"`
	StateDescription         *string               `json:"state_description,omitempty"`
	AuditLogRetention        *time.Duration        `json:"audit_log_retention,omitempty"`
	LoginPolicyFactors       *bool                 `json:"login_policy_factors,omitempty"`
	LoginPolicyIDP           *bool                 `json:"login_policy_idp,omitempty"`
	LoginPolicyPasswordless  *bool                 `json:"login_policy_passwordless,omitempty"`
	LoginPolicyRegistration  *bool                 `json:"login_policy_registration,omitempty"`
	LoginPolicyUsernameLogin *bool                 `json:"login_policy_username_login,omitempty"`
	PasswordComplexityPolicy *bool                 `json:"password_complexity_policy,omitempty"`
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

func ChangeState(State domain.FeaturesState) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.State = &State
	}
}

func ChangeStateDescription(statusDescription string) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.StateDescription = &statusDescription
	}
}

func ChangeAuditLogRetention(retention time.Duration) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.AuditLogRetention = &retention
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

func ChangePasswordComplexityPolicy(passwordComplexityPolicy bool) func(event *FeaturesSetEvent) {
	return func(e *FeaturesSetEvent) {
		e.PasswordComplexityPolicy = &passwordComplexityPolicy
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
