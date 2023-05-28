package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	DomainPolicyAddedEventType   = instanceEventTypePrefix + policy.DomainPolicyAddedEventType
	DomainPolicyChangedEventType = instanceEventTypePrefix + policy.DomainPolicyChangedEventType
)

type DomainPolicyAddedEvent struct {
	policy.DomainPolicyAddedEvent
}

func NewDomainPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomain,
	smtpSenderAddressMatchesInstanceDomain bool,
) *DomainPolicyAddedEvent {
	return &DomainPolicyAddedEvent{
		DomainPolicyAddedEvent: *policy.NewDomainPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DomainPolicyAddedEventType),
			userLoginMustBeDomain,
			validateOrgDomain,
			smtpSenderAddressMatchesInstanceDomain,
		),
	}
}

func DomainPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DomainPolicyAddedEvent{DomainPolicyAddedEvent: *e.(*policy.DomainPolicyAddedEvent)}, nil
}

type DomainPolicyChangedEvent struct {
	policy.DomainPolicyChangedEvent
}

func NewDomainPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.DomainPolicyChanges,
) (*DomainPolicyChangedEvent, error) {
	changedEvent, err := policy.NewDomainPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DomainPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &DomainPolicyChangedEvent{DomainPolicyChangedEvent: *changedEvent}, nil
}

func DomainPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DomainPolicyChangedEvent{DomainPolicyChangedEvent: *e.(*policy.DomainPolicyChangedEvent)}, nil
}
