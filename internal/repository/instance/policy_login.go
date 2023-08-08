package instance

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	LoginPolicyAddedEventType   = instanceEventTypePrefix + policy.LoginPolicyAddedEventType
	LoginPolicyChangedEventType = instanceEventTypePrefix + policy.LoginPolicyChangedEventType
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
	forceMFA,
	forceMFALocalOnly,
	hidePasswordReset,
	ignoreUnknownUsernames,
	allowDomainDiscovery,
	disableLoginWithEmail,
	disableLoginWithPhone bool,
	passwordlessType domain.PasswordlessType,
	defaultRedirectURI string,
	passwordCheckLifetime,
	externalLoginCheckLifetime,
	mfaInitSkipLifetime,
	secondFactorCheckLifetime,
	multiFactorCheckLifetime time.Duration,
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
			forceMFALocalOnly,
			hidePasswordReset,
			ignoreUnknownUsernames,
			allowDomainDiscovery,
			disableLoginWithEmail,
			disableLoginWithPhone,
			passwordlessType,
			defaultRedirectURI,
			passwordCheckLifetime,
			externalLoginCheckLifetime,
			mfaInitSkipLifetime,
			secondFactorCheckLifetime,
			multiFactorCheckLifetime),
	}
}

func LoginPolicyAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
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

func LoginPolicyChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.LoginPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &LoginPolicyChangedEvent{LoginPolicyChangedEvent: *e.(*policy.LoginPolicyChangedEvent)}, nil
}
