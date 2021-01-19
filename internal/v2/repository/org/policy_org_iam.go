package org

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	//TODO: enable when possible
	//OrgIAMPolicyAddedEventType   = orgEventTypePrefix + policy.OrgIAMPolicyAddedEventType
	//OrgIAMPolicyChangedEventType = orgEventTypePrefix + policy.OrgIAMPolicyChangedEventType
	OrgIAMPolicyAddedEventType   = orgEventTypePrefix + "iam.policy.added"
	OrgIAMPolicyChangedEventType = orgEventTypePrefix + "iam.policy.changed"
	OrgIAMPolicyRemovedEventType = orgEventTypePrefix + "iam.policy.removed"
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
	changes []policy.OrgIAMPolicyChanges,
) (*OrgIAMPolicyChangedEvent, error) {
	changedEvent, err := policy.NewOrgIAMPolicyChangedEvent(
		eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyChangedEventType),
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

type OrgIAMPolicyRemovedEvent struct {
	policy.OrgIAMPolicyRemovedEvent
}

func NewOrgIAMPolicyRemovedEvent(
	ctx context.Context,
) *OrgIAMPolicyRemovedEvent {
	return &OrgIAMPolicyRemovedEvent{
		OrgIAMPolicyRemovedEvent: *policy.NewOrgIAMPolicyRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, OrgIAMPolicyRemovedEventType),
		),
	}
}

func OrgIAMPolicyRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.OrgIAMPolicyRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &OrgIAMPolicyRemovedEvent{OrgIAMPolicyRemovedEvent: *e.(*policy.OrgIAMPolicyRemovedEvent)}, nil
}
