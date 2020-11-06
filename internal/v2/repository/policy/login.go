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

	current *LoginPolicyAggregate
	changed *LoginPolicyAggregate
}

func (e *LoginPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *LoginPolicyChangedEvent) Data() interface{} {
	changes := map[string]interface{}{}
	if e.current.AllowExternalIDP != e.changed.AllowExternalIDP {
		changes["allowUsernamePassword"] = e.changed.AllowExternalIDP
	}
	if e.current.AllowRegister != e.changed.AllowRegister {
		changes["allowRegister"] = e.changed.AllowExternalIDP
	}
	if e.current.AllowExternalIDP != e.changed.AllowExternalIDP {
		changes["allowExternalIdp"] = e.changed.AllowExternalIDP
	}

	return changes
}

func NewLoginPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *LoginPolicyAggregate,
) *LoginPolicyChangedEvent {

	return &LoginPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			LoginPolicyChangedEventType,
		),
	}
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
