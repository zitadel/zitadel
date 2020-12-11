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

type AddedEvent struct {
	password_complexity.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	minLength uint64,
	hasLowercase,
	hasUppercase,
	hasNumber,
	hasSymbol bool,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *password_complexity.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordComplexityPolicyAddedEventType),
			minLength,
			hasLowercase,
			hasUppercase,
			hasNumber,
			hasSymbol),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_complexity.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*password_complexity.AddedEvent)}, nil
}

type ChangedEvent struct {
	password_complexity.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	minLength uint64,
	hasLowerCase,
	hasUpperCase,
	hasNumber,
	hasSymbol bool,
) (*ChangedEvent, error) {
	event := password_complexity.NewChangedEvent(
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
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_complexity.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*password_complexity.ChangedEvent)}, nil
}
