package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	LoginPolicyAddedEventType   = "policy.login.added"
	LoginPolicyChangedEventType = "policy.login.changed"
	LoginPolicyRemovedEventType = "policy.login.removed"
)

type LoginPolicyAggregate struct {
	eventstore.Aggregate

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
}

type LoginPolicyReadModel struct {
	eventstore.ReadModel

	AllowUserNamePassword bool
	AllowRegister         bool
	AllowExternalIDP      bool
}

type LoginPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool `json:"allowUsernamePassword"`
	AllowRegister         bool `json:"allowRegister"`
	AllowExternalIDP      bool `json:"allowExternalIdp"`
	// TODO: IDPProviders
}

func (e *LoginPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyAddedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyAddedEvent(
	ctx context.Context,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP bool,
) *LoginPolicyAddedEvent {

	return &LoginPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LoginPolicyAddedEventType,
		),
		AllowExternalIDP:      allowExternalIDP,
		AllowRegister:         allowRegister,
		AllowUserNamePassword: allowUserNamePassword,
	}
}

type LoginPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	AllowUserNamePassword bool `json:"allowUsernamePassword,omitempty"`
	AllowRegister         bool `json:"allowRegister"`
	AllowExternalIDP      bool `json:"allowExternalIdp"`
}

func (e *LoginPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	return e
}

func NewLoginPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *LoginPolicyAggregate,
) *LoginPolicyChangedEvent {

	e := &LoginPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LoginPolicyChangedEventType,
		),
	}

	if current.AllowUserNamePassword != changed.AllowUserNamePassword {
		e.AllowUserNamePassword = changed.AllowUserNamePassword
	}
	if current.AllowRegister != changed.AllowRegister {
		e.AllowRegister = changed.AllowRegister
	}
	if current.AllowExternalIDP != changed.AllowExternalIDP {
		e.AllowExternalIDP = changed.AllowExternalIDP
	}

	return e
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

func NewLoginPolicyRemovedEvent(ctx context.Context) *LoginPolicyRemovedEvent {
	return &LoginPolicyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LoginPolicyRemovedEventType,
		),
	}
}
