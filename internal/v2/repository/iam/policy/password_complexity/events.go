package password_complexity

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_complexity"
)

const (
	iamEventPrefix                           = eventstore.EventType("iam.")
	PasswordComplexityPolicyAddedEventType   = iamEventPrefix + password_complexity.PasswordComplexityPolicyAddedEventType
	PasswordComplexityPolicyChangedEventType = iamEventPrefix + password_complexity.PasswordComplexityPolicyChangedEventType
)

type PasswordComplexityPolicyAddedEvent struct {
	password_complexity.PasswordComplexityPolicyAddedEvent
}

func NewPasswordComplexityPolicyAddedEventEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) *PasswordComplexityPolicyAddedEvent {
	return &PasswordComplexityPolicyAddedEvent{
		PasswordComplexityPolicyAddedEvent: *password_complexity.NewPasswordComplexityPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyAddedEventType),
			minLength,
			hasLowercase,
			hasUppercase,
			hasNumber,
			hasSymbol),
	}
}

func PasswordComplexityPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_complexity.PasswordComplexityPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyAddedEvent{PasswordComplexityPolicyAddedEvent: *e.(*password_complexity.PasswordComplexityPolicyAddedEvent)}, nil
}

type PasswordComplexityPolicyChangedEvent struct {
	password_complexity.PasswordComplexityPolicyChangedEvent
}

func PasswordComplexityPolicyChangedEventFromExisting(
	ctx context.Context,
	current *PasswordComplexityPolicyWriteModel,
	minLength uint64,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
) (*PasswordComplexityPolicyChangedEvent, error) {
	event := password_complexity.NewPasswordComplexityPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			PasswordComplexityPolicyChangedEventType,
		),
		&current.Policy,
		minLength,
		hasLowerCase,
		hasUpperCase,
		hasNumber,
		hasSymbol,
	)
	return &PasswordComplexityPolicyChangedEvent{
		*event,
	}, nil
}

func PasswordComplexityPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_complexity.PasswordComplexityPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordComplexityPolicyChangedEvent{PasswordComplexityPolicyChangedEvent: *e.(*password_complexity.PasswordComplexityPolicyChangedEvent)}, nil
}
