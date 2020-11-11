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
	eventstore.Aggregate

	MinLength    int
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

type PasswordComplexityPolicyReadModel struct {
	eventstore.ReadModel

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

	MinLength    int  `json:"minLength"`
	HasLowercase bool `json:"hasLowercase"`
	HasUpperCase bool `json:"hasUppercase"`
	HasNumber    bool `json:"hasNumber"`
	HasSymbol    bool `json:"hasSymbol"`
}

func (e *PasswordComplexityPolicyChangedEvent) CheckPrevious() bool {
	return true
}

func (e *PasswordComplexityPolicyChangedEvent) Data() interface{} {
	return e
}

func NewPasswordComplexityPolicyChangedEvent(
	ctx context.Context,
	current,
	changed *PasswordComplexityPolicyAggregate,
) *PasswordComplexityPolicyChangedEvent {

	e := &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			PasswordComplexityPolicyChangedEventType,
		),
	}

	if current.MinLength != changed.MinLength {
		e.MinLength = changed.MinLength
	}
	if current.HasLowercase != changed.HasLowercase {
		e.HasLowercase = changed.HasLowercase
	}
	if current.HasUpperCase != changed.HasUpperCase {
		e.HasUpperCase = changed.HasUpperCase
	}
	if current.HasNumber != changed.HasNumber {
		e.HasNumber = changed.HasNumber
	}
	if current.HasSymbol != changed.HasSymbol {
		e.HasSymbol = changed.HasSymbol
	}

	return e
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
			PasswordComplexityPolicyRemovedEventType,
		),
	}
}
