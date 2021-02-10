package iam

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	OrgIAMPolicyAddedEventType   = iamEventTypePrefix + policy.OrgIAMPolicyAddedEventType
	OrgIAMPolicyChangedEventType = iamEventTypePrefix + policy.OrgIAMPolicyChangedEventType
)

type OrgIAMPolicyAddedEvent struct {
	policy.OrgIAMPolicyAddedEvent
}

func NewOrgIAMPolicyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {
	return &OrgIAMPolicyAddedEvent{
		OrgIAMPolicyAddedEvent: *policy.NewOrgIAMPolicyAddedEvent(
			eventstore.NewBaseEventForPush(
				ctx,
				aggregate,
				OrgIAMPolicyAddedEventType),
			userLoginMustBeDomain,
		),
	}
}

func OrgIAMPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.OrgIAMPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyAddedEvent{OrgIAMPolicyAddedEvent: *e.(*policy.OrgIAMPolicyAddedEvent)}, nil
}

type OrgIAMPolicyChangedEvent struct {
	policy.OrgIAMPolicyChangedEvent
}

func NewOrgIAMPolicyChangedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	changes []policy.OrgIAMPolicyChanges,
) (*OrgIAMPolicyChangedEvent, error) {
	changedEvent, err := policy.NewOrgIAMPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			OrgIAMPolicyChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &OrgIAMPolicyChangedEvent{OrgIAMPolicyChangedEvent: *changedEvent}, nil
}

func OrgIAMPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.OrgIAMPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyChangedEvent{OrgIAMPolicyChangedEvent: *e.(*policy.OrgIAMPolicyChangedEvent)}, nil
}
