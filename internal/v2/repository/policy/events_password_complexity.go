package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordComplexityPolicyAddedEventType = "policy.password.complexity.added"
)

type PasswordComplexityPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    int  `json:"minLength"`
	HasLowercase bool `json:"hasLowercase"`
	HasUpperCase bool `json:"hasUppercase"`
	HasNumber    bool `json:"hasNumber"`
	HasSymbol    bool `json:"hasSymbol"`
}

func (e *PasswordComplexityPolicyAddedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordComplexityPolicyAddedEvent) Data() interface{} {
	return e
}

func NewPasswordComplexityPolicyAddedEvent(
	ctx context.Context,
	service string,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
	minLength int,
) *PasswordComplexityPolicyAddedEvent {

	return &PasswordComplexityPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			service,
			PasswordComplexityPolicyAddedEventType,
		),
		HasLowercase: hasLowerCase,
		HasNumber:    hasNumber,
		HasSymbol:    hasSymbol,
		HasUpperCase: hasUpperCase,
		MinLength:    minLength,
	}
}
