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

type AddedEvent struct {
	password_age.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	expireWarnDays,
	maxAgeDays uint64,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *password_age.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, PasswordAgePolicyAddedEventType),
			expireWarnDays,
			maxAgeDays),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_age.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*password_age.AddedEvent)}, nil
}

type ChangedEvent struct {
	password_age.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	expireWarnDays,
	maxAgeDays uint64,
) (*ChangedEvent, error) {
	event := password_age.NewChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			PasswordAgePolicyChangedEventType,
		),
		&current.Policy,
		expireWarnDays,
		maxAgeDays,
	)
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := password_age.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*password_age.ChangedEvent)}, nil
}
