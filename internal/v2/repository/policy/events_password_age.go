package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordAgePolicyAddedEventType = "policy.password.age.added"
)

type PasswordAgePolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	ExpireWarnDays int `json:"expireWarnDays"`
	MaxAgeDays     int `json:"maxAgeDays"`
}

func (e *PasswordAgePolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordAgePolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	service string,
	expireWarnDays,
	maxAgeDays int,
) *PasswordAgePolicyAddedEvent {
	return &PasswordAgePolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			PasswordAgePolicyAddedEventType,
		),
		ExpireWarnDays: expireWarnDays,
		MaxAgeDays:     maxAgeDays,
	}
}
