package instance

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	InstanceDomainPolicyAddedEventType   = instanceEventTypePrefix + policy.DomainPolicyAddedEventType
	InstanceDomainPolicyChangedEventType = instanceEventTypePrefix + policy.DomainPolicyChangedEventType
)

type InstanceDomainPolicyAddedEvent struct {
	policy.DomainPolicyAddedEvent
}

func NewInstnaceDomainPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool,
) *InstanceDomainPolicyAddedEvent {
	return &InstanceDomainPolicyAddedEvent{
		DomainPolicyAddedEvent: *policy.NewDomainPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				InstanceDomainPolicyAddedEventType),
			userLoginMustBeDomain,
		),
	}
}

func InstanceDomainPolicyAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &InstanceDomainPolicyAddedEvent{DomainPolicyAddedEvent: *e.(*policy.DomainPolicyAddedEvent)}, nil
}

type InstanceDomainPolicyChangedEvent struct {
	policy.DomainPolicyChangedEvent
}

func NewInstanceDomainPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.OrgPolicyChanges,
) (*InstanceDomainPolicyChangedEvent, error) {
	changedEvent, err := policy.NewDomainPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			InstanceDomainPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &InstanceDomainPolicyChangedEvent{DomainPolicyChangedEvent: *changedEvent}, nil
}

func InstanceDomainPolicyChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &InstanceDomainPolicyChangedEvent{DomainPolicyChangedEvent: *e.(*policy.DomainPolicyChangedEvent)}, nil
}
