package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

const (
	PrivacyPolicyAddedEventType   = instanceEventTypePrefix + policy.PrivacyPolicyAddedEventType
	PrivacyPolicyChangedEventType = instanceEventTypePrefix + policy.PrivacyPolicyChangedEventType
)

type PrivacyPolicyAddedEvent struct {
	policy.PrivacyPolicyAddedEvent
}

func NewPrivacyPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	tosLink,
	privacyLink,
	helpLink string,
	supportEmail domain.EmailAddress,
	docsLink, customLink, customLinkText string,
) *PrivacyPolicyAddedEvent {
	return &PrivacyPolicyAddedEvent{
		PrivacyPolicyAddedEvent: *policy.NewPrivacyPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				PrivacyPolicyAddedEventType),
			tosLink,
			privacyLink,
			helpLink,
			supportEmail,
			docsLink,
			customLink,
			customLinkText),
	}
}

func PrivacyPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PrivacyPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PrivacyPolicyAddedEvent{PrivacyPolicyAddedEvent: *e.(*policy.PrivacyPolicyAddedEvent)}, nil
}

type PrivacyPolicyChangedEvent struct {
	policy.PrivacyPolicyChangedEvent
}

func NewPrivacyPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.PrivacyPolicyChanges,
) (*PrivacyPolicyChangedEvent, error) {
	changedEvent, err := policy.NewPrivacyPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			PrivacyPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &PrivacyPolicyChangedEvent{PrivacyPolicyChangedEvent: *changedEvent}, nil
}

func PrivacyPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.PrivacyPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &PrivacyPolicyChangedEvent{PrivacyPolicyChangedEvent: *e.(*policy.PrivacyPolicyChangedEvent)}, nil
}
