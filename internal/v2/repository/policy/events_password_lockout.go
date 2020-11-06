package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordLockoutPolicyAddedEventType = "policy.password.lockout.added"
)

type PasswordLockoutPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MaxAttempts         int  `json:"maxAttempts"`
	ShowLockOutFailures bool `json:"showLockOutFailures"`
}

func (e *PasswordLockoutPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordLockoutPolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordLockoutPolicyAddedEvent(
	ctx context.Context,
	service string,
	maxAttempts int,
	showLockOutFailures bool,
) *PasswordLockoutPolicyAddedEvent {

	return &PasswordLockoutPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			LabelPolicyAddedEventType,
		),
		MaxAttempts:         maxAttempts,
		ShowLockOutFailures: showLockOutFailures,
	}
}
