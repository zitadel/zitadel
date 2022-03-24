package org

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	OrgDomainPolicyAddedEventType   = orgEventTypePrefix + policy.DomainPolicyAddedEventType
	OrgDomainPolicyChangedEventType = orgEventTypePrefix + policy.DomainPolicyChangedEventType
	OrgDomainPolicyRemovedEventType = orgEventTypePrefix + policy.DomainPolicyRemovedEventType
)

type OrgDomainPolicyAddedEvent struct {
	policy.DomainPolicyAddedEvent
}

func NewOrgDomainPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool,
) *OrgDomainPolicyAddedEvent {
	return &OrgDomainPolicyAddedEvent{
		DomainPolicyAddedEvent: *policy.NewDomainPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OrgDomainPolicyAddedEventType),
			userLoginMustBeDomain,
		),
	}
}

func OrgDomainPolicyAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgDomainPolicyAddedEvent{DomainPolicyAddedEvent: *e.(*policy.DomainPolicyAddedEvent)}, nil
}

type OrgDomainPolicyChangedEvent struct {
	policy.DomainPolicyChangedEvent
}

func NewOrgDomainPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.OrgPolicyChanges,
) (*OrgDomainPolicyChangedEvent, error) {
	changedEvent, err := policy.NewDomainPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgDomainPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OrgDomainPolicyChangedEvent{DomainPolicyChangedEvent: *changedEvent}, nil
}

func OrgDomainPolicyChangedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgDomainPolicyChangedEvent{DomainPolicyChangedEvent: *e.(*policy.DomainPolicyChangedEvent)}, nil
}

type OrgDomainPolicyRemovedEvent struct {
	policy.DomainPolicyRemovedEvent
}

func NewOrgDomainPolicyRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *OrgDomainPolicyRemovedEvent {
	return &OrgDomainPolicyRemovedEvent{
		DomainPolicyRemovedEvent: *policy.NewDomainPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OrgDomainPolicyRemovedEventType),
		),
	}
}

func OrgDomainPolicyRemovedEventMapper(event *repository.Event) (eventstore.Event, error) {
	e, err := policy.DomainPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgDomainPolicyRemovedEvent{DomainPolicyRemovedEvent: *e.(*policy.DomainPolicyRemovedEvent)}, nil
}
