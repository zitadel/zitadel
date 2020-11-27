package policy

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	PasswordComplexityPolicyAddedEventType   = "policy.password.complexity.added"
	PasswordComplexityPolicyChangedEventType = "policy.password.complexity.changed"
	PasswordComplexityPolicyRemovedEventType = "policy.password.complexity.removed"
)

type PasswordComplexityPolicyAggregate struct {
	eventstore.Aggregate

	MinLength    uint8
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

type PasswordComplexityPolicyReadModel struct {
	eventstore.ReadModel

	MinLength    uint8
	HasLowercase bool
	HasUpperCase bool
	HasNumber    bool
	HasSymbol    bool
}

func (rm *PasswordComplexityPolicyReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *PasswordComplexityPolicyAddedEvent:
			rm.MinLength = e.MinLength
			rm.HasLowercase = e.HasLowercase
			rm.HasUpperCase = e.HasUpperCase
			rm.HasNumber = e.HasNumber
			rm.HasSymbol = e.HasSymbol
		case *PasswordComplexityPolicyChangedEvent:
			rm.MinLength = e.MinLength
			rm.HasLowercase = e.HasLowercase
			rm.HasUpperCase = e.HasUpperCase
			rm.HasNumber = e.HasNumber
			rm.HasSymbol = e.HasSymbol
		}
	}
	return rm.ReadModel.Reduce()
}

type PasswordComplexityPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint8 `json:"minLength,omitempty"`
	HasLowercase bool  `json:"hasLowercase"`
	HasUpperCase bool  `json:"hasUppercase"`
	HasNumber    bool  `json:"hasNumber"`
	HasSymbol    bool  `json:"hasSymbol"`
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
	minLength uint8,
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

func PasswordComplexityPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordComplexityPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-wYxlM", "unable to unmarshal policy")
	}

	return e, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MinLength    uint8 `json:"minLength"`
	HasLowercase bool  `json:"hasLowercase"`
	HasUpperCase bool  `json:"hasUppercase"`
	HasNumber    bool  `json:"hasNumber"`
	HasSymbol    bool  `json:"hasSymbol"`
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

func PasswordComplexityPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &PasswordComplexityPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-zBGB0", "unable to unmarshal policy")
	}

	return e, nil
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

func PasswordComplexityPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &PasswordComplexityPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
