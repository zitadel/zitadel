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

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool             `json:"allowUsernamePassword"`
	AllowRegister         bool             `json:"allowRegister"`
	AllowExternalIDP      bool             `json:"allowExternalIdp"`
	ForceMFA              bool             `json:"forceMFA"`
	PasswordlessType      PasswordlessType `json:"passwordlessType"`
}

func (e *AddedEvent) CheckPrevious() bool {
	return true
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func NewAddedEvent(
	base *eventstore.BaseEvent,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType PasswordlessType,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent:             *base,
		AllowExternalIDP:      allowExternalIDP,
		AllowRegister:         allowRegister,
		AllowUserNamePassword: allowUserNamePassword,
		ForceMFA:              forceMFA,
		PasswordlessType:      passwordlessType,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-nWndT", "unable to unmarshal policy")
	}

	return e, nil
}

type ChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool             `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool             `json:"allowRegister"`
	AllowExternalIDP      bool             `json:"allowExternalIdp"`
	ForceMFA              bool             `json:"forceMFA"`
	PasswordlessType      PasswordlessType `json:"passwordlessType"`
}

func (e *ChangedEvent) CheckPrevious() bool {
	return true
}

func (e *ChangedEvent) Data() interface{} {
	return e
}

func NewChangedEvent(
	base *eventstore.BaseEvent,
	current *WriteModel,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType PasswordlessType,
) *ChangedEvent {

	e := &ChangedEvent{
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

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &ChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-ehssl", "unable to unmarshal policy")
	}

	return e, nil
}

type RemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *RemovedEvent) CheckPrevious() bool {
	return true
}

func (e *RemovedEvent) Data() interface{} {
	return nil
}

func NewRemovedEvent(base *eventstore.BaseEvent) *RemovedEvent {
	return &RemovedEvent{
		BaseEvent: *base,
	}
}

func RemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &RemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
