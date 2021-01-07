package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	OrgIAMPolicyAddedEventType   = orgEventTypePrefix + policy.OrgIAMPolicyAddedEventType
	OrgIAMPolicyChangedEventType = orgEventTypePrefix + policy.OrgIAMPolicyChangedEventType
)

type OrgIAMPolicyAddedEvent struct {
	policy.OrgIAMPolicyAddedEvent
}

func NewOrgIAMPolicyAddedEvent(
	ctx context.Context,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {
	return &OrgIAMPolicyAddedEvent{
		OrgIAMPolicyAddedEvent: *policy.NewOrgIAMPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyAddedEventType),
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
) *OrgIAMPolicyChangedEvent {
	return &OrgIAMPolicyChangedEvent{
		OrgIAMPolicyChangedEvent: *policy.NewOrgIAMPolicyChangedEvent(
			eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyChangedEventType),
		),
	}
}

func OrgIAMPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.OrgIAMPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyChangedEvent{OrgIAMPolicyChangedEvent: *e.(*policy.OrgIAMPolicyChangedEvent)}, nil
}
