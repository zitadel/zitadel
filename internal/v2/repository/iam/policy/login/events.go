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

type AddedEvent struct {
	login.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType login.PasswordlessType,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *login.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, login.LoginPolicyAddedEventType),
			allowUsernamePassword,
			allowRegister,
			allowExternalIDP,
			forceMFA,
			passwordlessType),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := login.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*login.AddedEvent)}, nil
}

type ChangedEvent struct {
	login.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	allowUsernamePassword,
	allowRegister,
	allowExternalIDP,
	forceMFA bool,
	passwordlessType login.PasswordlessType,
) (*ChangedEvent, error) {

	event := login.NewChangedEvent(
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
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := login.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*login.ChangedEvent)}, nil
}
