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

type AddedEvent struct {
	password_lockout.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	maxAttempts uint64,
	showLockoutFailure bool,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *password_lockout.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordLockoutPolicyAddedEventType),
			maxAttempts,
			showLockoutFailure),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_lockout.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*password_lockout.AddedEvent)}, nil
}

type ChangedEvent struct {
	password_lockout.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	maxAttempts uint64,
	showLockoutFailure bool,
) (*ChangedEvent, error) {
	event := password_lockout.NewChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			PasswordLockoutPolicyChangedEventType,
		),
		&current.WriteModel,
		maxAttempts,
		showLockoutFailure,
	)
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_lockout.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*password_lockout.ChangedEvent)}, nil
}
