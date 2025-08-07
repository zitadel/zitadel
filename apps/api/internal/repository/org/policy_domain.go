package org

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	DomainPolicyAddedEventType   = orgEventTypePrefix + policy.DomainPolicyAddedEventType
	DomainPolicyChangedEventType = orgEventTypePrefix + policy.DomainPolicyChangedEventType
	DomainPolicyRemovedEventType = orgEventTypePrefix + policy.DomainPolicyRemovedEventType
)

type DomainPolicyAddedEvent struct {
	policy.DomainPolicyAddedEvent
}

func NewDomainPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain,
	validateOrgDomains,
	smtpSenderAddressMatchesInstanceDomain bool,
) *DomainPolicyAddedEvent {
	return &DomainPolicyAddedEvent{
		DomainPolicyAddedEvent: *policy.NewDomainPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DomainPolicyAddedEventType),
			userLoginMustBeDomain,
			validateOrgDomains,
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

type DomainPolicyRemovedEvent struct {
	policy.DomainPolicyRemovedEvent
}

func NewDomainPolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *DomainPolicyRemovedEvent {
	return &DomainPolicyRemovedEvent{
		DomainPolicyRemovedEvent: *policy.NewDomainPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				DomainPolicyRemovedEventType),
		),
	}
}

func DomainPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &DomainPolicyRemovedEvent{DomainPolicyRemovedEvent: *e.(*policy.DomainPolicyRemovedEvent)}, nil
}
