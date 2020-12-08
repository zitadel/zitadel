package login

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/login"
)

var (
	iamEventPrefix              = eventstore.EventType("iam.")
	LoginPolicyAddedEventType   = iamEventPrefix + login.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = iamEventPrefix + login.LoginPolicyChangedEventType

	LoginPolicyIDPProviderAddedEventType   = iamEventPrefix + login.LoginPolicyIDPProviderAddedEventType
	LoginPolicyIDPProviderRemovedEventType = iamEventPrefix + login.LoginPolicyIDPProviderRemovedEventType
)

type LoginPolicyAddedEvent struct {
	login.LoginPolicyAddedEvent
}

func NewLoginPolicyAddedEventEvent(
	ctx context.Context,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType login.PasswordlessType,
) *LoginPolicyAddedEvent {
	return &LoginPolicyAddedEvent{
		LoginPolicyAddedEvent: *login.NewLoginPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, login.LoginPolicyAddedEventType),
			allowUsernamePassword,
			allowRegister,
			allowExternalIDP,
			forceMFA,
			passwordlessType),
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := login.LoginPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyAddedEvent{LoginPolicyAddedEvent: *e.(*login.LoginPolicyAddedEvent)}, nil
}

type LoginPolicyChangedEvent struct {
	login.LoginPolicyChangedEvent
}

func LoginPolicyChangedEventFromExisting(
	ctx context.Context,
	current *LoginPolicyWriteModel,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType login.PasswordlessType,
) (*LoginPolicyChangedEvent, error) {

	event := login.NewLoginPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			LoginPolicyChangedEventType,
		),
		&current.Policy,
		allowUsernamePassword,
		allowRegister,
		allowExternalIDP,
		forceMFA,
		passwordlessType,
	)
	return &LoginPolicyChangedEvent{
		*event,
	}, nil
}

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := login.LoginPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyChangedEvent{LoginPolicyChangedEvent: *e.(*login.LoginPolicyChangedEvent)}, nil
}
