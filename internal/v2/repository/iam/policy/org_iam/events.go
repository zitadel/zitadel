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

type AddedEvent struct {
	org_iam.AddedEvent
}

func NewAddedEvent(
	ctx context.Context,
	userLoginMustBeDomain bool,
) *AddedEvent {
	return &AddedEvent{
		AddedEvent: *org_iam.NewAddedEvent(
			eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyAddedEventType),
			userLoginMustBeDomain,
		),
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := org_iam.AddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &AddedEvent{AddedEvent: *e.(*org_iam.AddedEvent)}, nil
}

type ChangedEvent struct {
	org_iam.ChangedEvent
}

func ChangedEventFromExisting(
	ctx context.Context,
	current *WriteModel,
	userLoginMustBeDomain bool,
) (*ChangedEvent, error) {
	event := org_iam.NewChangedEvent(
		eventstore.NewBaseEventForPush(
			ctx,
			OrgIAMPolicyChangedEventType,
		),
		&current.WriteModel,
		userLoginMustBeDomain,
	)
	return &ChangedEvent{
		*event,
	}, nil
}

func ChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := org_iam.ChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &ChangedEvent{ChangedEvent: *e.(*org_iam.ChangedEvent)}, nil
}
