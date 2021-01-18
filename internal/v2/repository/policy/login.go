package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
)

const (
	loginPolicyPrefix           = "policy.login."
	LoginPolicyAddedEventType   = loginPolicyPrefix + "added"
	LoginPolicyChangedEventType = loginPolicyPrefix + "changed"
	LoginPolicyRemovedEventType = loginPolicyPrefix + "removed"
)

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool                    `json:"allowRegister,omitempty"`
	AllowExternalIDP      bool                    `json:"allowExternalIdp,omitempty"`
	ForceMFA              bool                    `json:"forceMFA,omitempty"`
	PasswordlessType      domain.PasswordlessType `json:"passwordlessType,omitempty"`
}

func (e *LoginPolicyAddedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyAddedEvent(
	base *eventstore.BaseEvent,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType domain.PasswordlessType,
) *LoginPolicyAddedEvent {
	return &LoginPolicyAddedEvent{
		BaseEvent:             *base,
		AllowExternalIDP:      allowExternalIDP,
		AllowRegister:         allowRegister,
		AllowUserNamePassword: allowUserNamePassword,
		ForceMFA:              forceMFA,
		PasswordlessType:      passwordlessType,
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LoginPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-nWndT", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword *bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister         *bool                    `json:"allowRegister,omitempty"`
	AllowExternalIDP      *bool                    `json:"allowExternalIdp,omitempty"`
	ForceMFA              *bool                    `json:"forceMFA,omitempty"`
	PasswordlessType      *domain.PasswordlessType `json:"passwordlessType,omitempty"`
}

type LoginPolicyEventData struct {
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []LoginPolicyChanges,
) (*LoginPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-ADg34", "Errors.NoChangesFound")
	}
	changeEvent := &LoginPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type LoginPolicyChanges func(*LoginPolicyChangedEvent)

func ChangeAllowUserNamePassword(allowUserNamePassword bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowUserNamePassword = &allowUserNamePassword
	}
}

func ChangeAllowRegister(allowRegister bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowRegister = &allowRegister
	}
}

func ChangeAllowExternalIDP(allowExternalIDP bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.AllowExternalIDP = &allowExternalIDP
	}
}

func ChangeForceMFA(forceMFA bool) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.ForceMFA = &forceMFA
	}
}

func ChangePasswordlessType(passwordlessType domain.PasswordlessType) func(*LoginPolicyChangedEvent) {
	return func(e *LoginPolicyChangedEvent) {
		e.PasswordlessType = &passwordlessType
	}
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &LoginPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type LoginPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *LoginPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewLoginPolicyRemovedEvent(base *eventstore.BaseEvent) *LoginPolicyRemovedEvent {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func LoginPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
