package password_lockout

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_lockout"
)

var (
	iamEventPrefix                        = eventstore.EventType("iam.")
	PasswordLockoutPolicyAddedEventType   = iamEventPrefix + password_lockout.PasswordLockoutPolicyAddedEventType
	PasswordLockoutPolicyChangedEventType = iamEventPrefix + password_lockout.PasswordLockoutPolicyChangedEventType
)

type PasswordLockoutPolicyAddedEvent struct {
	password_lockout.PasswordLockoutPolicyAddedEvent
}

func NewPasswordLockoutPolicyAddedEvent(
	ctx context.Context,
	maxAttempts uint64,
	showLockoutFailure bool,
) *PasswordLockoutPolicyAddedEvent {
	return &PasswordLockoutPolicyAddedEvent{
		PasswordLockoutPolicyAddedEvent: *password_lockout.NewPasswordLockoutPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordLockoutPolicyAddedEventType),
			maxAttempts,
			showLockoutFailure),
	}
}

func PasswordLockoutPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_lockout.PasswordLockoutPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyAddedEvent{PasswordLockoutPolicyAddedEvent: *e.(*password_lockout.PasswordLockoutPolicyAddedEvent)}, nil
}

type PasswordLockoutPolicyChangedEvent struct {
	password_lockout.PasswordLockoutPolicyChangedEvent
}

func PasswordLockoutPolicyChangedEventFromExisting(
	ctx context.Context,
	current *PasswordLockoutPolicyWriteModel,
	maxAttempts uint64,
	showLockoutFailure bool,
) (*PasswordLockoutPolicyChangedEvent, error) {
	event := password_lockout.NewPasswordLockoutPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			PasswordLockoutPolicyChangedEventType,
		),
		&current.Policy,
		maxAttempts,
		showLockoutFailure,
	)
	return &PasswordLockoutPolicyChangedEvent{
		*event,
	}, nil
}

func PasswordLockoutPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_lockout.PasswordLockoutPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordLockoutPolicyChangedEvent{PasswordLockoutPolicyChangedEvent: *e.(*password_lockout.PasswordLockoutPolicyChangedEvent)}, nil
}
