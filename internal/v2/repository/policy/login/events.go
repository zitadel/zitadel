package login

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/idp/provider"
)

const (
	loginPolicyPrefix                      = "policy.login."
	LoginPolicyAddedEventType              = loginPolicyPrefix + "added"
	LoginPolicyChangedEventType            = loginPolicyPrefix + "changed"
	LoginPolicyRemovedEventType            = loginPolicyPrefix + "removed"
	LoginPolicyIDPProviderAddedEventType   = loginPolicyPrefix + provider.AddedEventType
	LoginPolicyIDPProviderRemovedEventType = loginPolicyPrefix + provider.RemovedEventType
)

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool             `json:"allowUsernamePassword"`
	AllowRegister         bool             `json:"allowRegister"`
	AllowExternalIDP      bool             `json:"allowExternalIdp"`
	ForceMFA              bool             `json:"forceMFA"`
	PasswordlessType      PasswordlessType `json:"passwordlessType"`
}

func (e *LoginPolicyAddedEvent) CheckPrevious() bool {
	return true
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
	passwordlessType PasswordlessType,
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

	AllowUserNamePassword bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool             `json:"allowRegister"`
	AllowExternalIDP      bool             `json:"allowExternalIdp"`
	ForceMFA              bool             `json:"forceMFA"`
	PasswordlessType      PasswordlessType `json:"passwordlessType"`
}

func (e *LoginPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyChangedEvent(
	base *eventstore.BaseEvent,
	current *LoginPolicyWriteModel,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType PasswordlessType,
) *LoginPolicyChangedEvent {

	e := &LoginPolicyChangedEvent{
		BaseEvent: *base,
	}

	if current.AllowUserNamePassword != allowUserNamePassword {
		e.AllowUserNamePassword = allowUserNamePassword
	}
	if current.AllowRegister != allowRegister {
		e.AllowRegister = allowRegister
	}
	if current.AllowExternalIDP != allowExternalIDP {
		e.AllowExternalIDP = allowExternalIDP
	}
	if current.ForceMFA != forceMFA {
		e.ForceMFA = forceMFA
	}
	if current.PasswordlessType != passwordlessType {
		e.PasswordlessType = passwordlessType
	}

	return e
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

func (e *LoginPolicyRemovedEvent) CheckPrevious() bool {
	return true
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
