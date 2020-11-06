package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	LoginPolicyAddedEventType = "policy.login.added"
)

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
	service string,
	allowUserNamePassword,
	allowRegister,
	allowExternalIDP bool,
) *LoginPolicyAddedEvent {

	return &LoginPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			LoginPolicyAddedEventType,
		),
		AllowExternalIDP:      allowExternalIDP,
		AllowRegister:         allowRegister,
		AllowUserNamePassword: allowUserNamePassword,
	}
}
