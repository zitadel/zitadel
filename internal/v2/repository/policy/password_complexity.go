package policy

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
)

const (
	PasswordComplexityPolicyAddedEventType   = "policy.password.complexity.added"
	PasswordComplexityPolicyChangedEventType = "policy.password.complexity.changed"
	PasswordComplexityPolicyRemovedEventType = "policy.password.complexity.removed"
)

type PasswordComplexityPolicyAggregate struct {
	MinLength    int
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

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
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
	minLength int,
) *PasswordComplexityPolicyAddedEvent {

	return &PasswordComplexityPolicyAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordComplexityPolicyAddedEventType,
		),
		HasLowercase: hasLowerCase,
		HasNumber:    hasNumber,
		HasSymbol:    hasSymbol,
		HasUpperCase: hasUpperCase,
		MinLength:    minLength,
	}
}

type PasswordComplexityPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	current *PasswordComplexityPolicyAggregate
	changed *PasswordComplexityPolicyAggregate
}

func (e *PasswordComplexityPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordComplexityPolicyChangedEvent) Data() interface{} {
	changes := map[string]interface{}{}

	if e.current.MinLength != e.changed.MinLength {
		changes["minLength"] = e.changed.MinLength
	}
	if e.current.HasLowercase != e.changed.HasLowercase {
		changes["hasLowercase"] = e.changed.HasLowercase
	}
	if e.current.HasUpperCase != e.changed.HasUpperCase {
		changes["hasUppercase"] = e.changed.HasUpperCase
	}
	if e.current.HasNumber != e.changed.HasNumber {
		changes["hasNumber"] = e.changed.HasNumber
	}
	if e.current.HasSymbol != e.changed.HasSymbol {
		changes["hasSymbol"] = e.changed.HasSymbol
	}

	return changes
}

func NewPasswordComplexityPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordComplexityPolicyAggregate,
) *PasswordComplexityPolicyChangedEvent {

	return &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordComplexityPolicyAddedEventType,
		),
		current: current,
		changed: changed,
	}
}

type PasswordComplexityPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PasswordComplexityPolicyRemovedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordComplexityPolicyRemovedEvent) Data() interface{} {
	return nil
}

func NewPasswordComplexityPolicyRemovedEvent(
	ctx context.Context,
) *PasswordComplexityPolicyRemovedEvent {

	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordComplexityPolicyChangedEventType,
		),
	}
}
