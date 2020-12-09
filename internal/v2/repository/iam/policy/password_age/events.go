package password_age

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/password_age"
)

var (
	iamEventPrefix                    = eventstore.EventType("iam.")
	PasswordAgePolicyAddedEventType   = iamEventPrefix + password_age.PasswordAgePolicyAddedEventType
	PasswordAgePolicyChangedEventType = iamEventPrefix + password_age.PasswordAgePolicyChangedEventType
)

type PasswordAgePolicyAddedEvent struct {
	password_age.PasswordAgePolicyAddedEvent
}

func NewPasswordAgePolicyAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays uint64,
) *PasswordAgePolicyAddedEvent {
	return &PasswordAgePolicyAddedEvent{
		PasswordAgePolicyAddedEvent: *password_age.NewPasswordAgePolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordAgePolicyAddedEventType),
			expireWarnDays,
			maxAgeDays),
	}
}

func PasswordAgePolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_age.PasswordAgePolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyAddedEvent{PasswordAgePolicyAddedEvent: *e.(*password_age.PasswordAgePolicyAddedEvent)}, nil
}

type PasswordAgePolicyChangedEvent struct {
	password_age.PasswordAgePolicyChangedEvent
}

func PasswordAgePolicyChangedEventFromExisting(
	ctx context.Context,
	current *PasswordAgePolicyWriteModel,
	expireWarnDays,
	maxAgeDays uint64,
) (*PasswordAgePolicyChangedEvent, error) {
	event := password_age.NewPasswordAgePolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
		&current.Policy,
		expireWarnDays,
		maxAgeDays,
	)
	return &PasswordAgePolicyChangedEvent{
		*event,
	}, nil
}

func PasswordAgePolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_age.PasswordAgePolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PasswordAgePolicyChangedEvent{PasswordAgePolicyChangedEvent: *e.(*password_age.PasswordAgePolicyChangedEvent)}, nil
}
