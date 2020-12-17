package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/business/domain"
)

const (
	loginPolicyPrefix           = "policy.login."
	LoginPolicyAddedEventType   = loginPolicyPrefix + "added"
	LoginPolicyChangedEventType = loginPolicyPrefix + "changed"
	LoginPolicyRemovedEventType = loginPolicyPrefix + "removed"
)

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool                    `json:"allowUsernamePassword"`
	AllowRegister         bool                    `json:"allowRegister"`
	AllowExternalIDP      bool                    `json:"allowExternalIdp"`
	ForceMFA              bool                    `json:"forceMFA"`
	PasswordlessType      domain.PasswordlessType `json:"passwordlessType"`
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

	AllowUserNamePassword bool                    `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool                    `json:"allowRegister"`
	AllowExternalIDP      bool                    `json:"allowExternalIdp"`
	ForceMFA              bool                    `json:"forceMFA"`
	PasswordlessType      domain.PasswordlessType `json:"passwordlessType"`
}

type LoginPolicyEventData struct {
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyChangedEvent(
	base *eventstore.BaseEvent,
) *LoginPolicyChangedEvent {
	return &LoginPolicyChangedEvent{
		BaseEvent: *base,
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
