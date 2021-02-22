package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	LoginPolicyAddedEventType   = iamEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = iamEventTypePrefix + policy.LoginPolicyChangedEventType
)

type LoginPolicyAddedEvent struct {
	policy.LoginPolicyAddedEvent
}

func NewLoginPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType domain.PasswordlessType,
) *LoginPolicyAddedEvent {
	return &LoginPolicyAddedEvent{
		LoginPolicyAddedEvent: *policy.NewLoginPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				LoginPolicyAddedEventType),
			allowUsernamePassword,
			allowRegister,
			allowExternalIDP,
			forceMFA,
			passwordlessType),
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LoginPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyAddedEvent{LoginPolicyAddedEvent: *e.(*policy.LoginPolicyAddedEvent)}, nil
}

type LoginPolicyChangedEvent struct {
	policy.LoginPolicyChangedEvent
}

func NewLoginPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.LoginPolicyChanges,
) (*LoginPolicyChangedEvent, error) {
	changedEvent, err := policy.NewLoginPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			LoginPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &LoginPolicyChangedEvent{LoginPolicyChangedEvent: *changedEvent}, nil
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.LoginPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyChangedEvent{LoginPolicyChangedEvent: *e.(*policy.LoginPolicyChangedEvent)}, nil
}
