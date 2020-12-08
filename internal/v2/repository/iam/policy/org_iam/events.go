package org_iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy/org_iam"
)

var (
	iamEventPrefix               = eventstore.EventType("iam.")
	OrgIAMPolicyAddedEventType   = iamEventPrefix + org_iam.OrgIAMPolicyAddedEventType
	OrgIAMPolicyChangedEventType = iamEventPrefix + org_iam.OrgIAMPolicyChangedEventType
)

type OrgIAMPolicyAddedEvent struct {
	org_iam.OrgIAMPolicyAddedEvent
}

func NewOrgIAMPolicyAddedEventEvent(
	ctx context.Context,
	userLoginMustBeDomain bool,
) *OrgIAMPolicyAddedEvent {
	return &OrgIAMPolicyAddedEvent{
		OrgIAMPolicyAddedEvent: *org_iam.NewOrgIAMPolicyAddedEvent(
			eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyAddedEventType),
			userLoginMustBeDomain,
		),
	}
}

func OrgIAMPolicyAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := org_iam.OrgIAMPolicyAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyAddedEvent{OrgIAMPolicyAddedEvent: *e.(*org_iam.OrgIAMPolicyAddedEvent)}, nil
}

type OrgIAMPolicyChangedEvent struct {
	org_iam.OrgIAMPolicyChangedEvent
}

func OrgIAMPolicyChangedEventFromExisting(
	ctx context.Context,
	current *OrgIAMPolicyWriteModel,
	userLoginMustBeDomain bool,
) (*OrgIAMPolicyChangedEvent, error) {
	event := org_iam.NewOrgIAMPolicyChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			OrgIAMPolicyChangedEventType,
		),
		&current.Policy,
		userLoginMustBeDomain,
	)
	return &OrgIAMPolicyChangedEvent{
		*event,
	}, nil
}

func OrgIAMPolicyChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := org_iam.OrgIAMPolicyChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyChangedEvent{OrgIAMPolicyChangedEvent: *e.(*org_iam.OrgIAMPolicyChangedEvent)}, nil
}
