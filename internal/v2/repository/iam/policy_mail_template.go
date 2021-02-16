package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	MailTemplateAddedEventType   = iamEventTypePrefix + policy.MailTemplatePolicyAddedEventType
	MailTemplateChangedEventType = iamEventTypePrefix + policy.MailTemplatePolicyChangedEventType
)

type MailTemplateAddedEvent struct {
	policy.MailTemplateAddedEvent
}

func NewMailTemplateAddedEvent(
	ctx context.Context,
	template []byte,
) *MailTemplateAddedEvent {
	return &MailTemplateAddedEvent{
		MailTemplateAddedEvent: *policy.NewMailTemplateAddedEvent(
			eventstore.NewBaseEventForPush(ctx, MailTemplateAddedEventType),
			template),
	}
}

func MailTemplateAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MailTemplateAddedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTemplateAddedEvent{MailTemplateAddedEvent: *e.(*policy.MailTemplateAddedEvent)}, nil
}

type MailTemplateChangedEvent struct {
	policy.MailTemplateChangedEvent
}

func NewMailTemplateChangedEvent(
	ctx context.Context,
	changes []policy.MailTemplateChanges,
) (*MailTemplateChangedEvent, error) {
	changedEvent, err := policy.NewMailTemplateChangedEvent(
		eventstore.NewBaseEventForPush(ctx, MailTemplateChangedEventType),
		changes,
	)
	if err != nil {
		return nil, err
	}
	return &MailTemplateChangedEvent{MailTemplateChangedEvent: *changedEvent}, nil
}

func MailTemplateChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MailTemplateChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTemplateChangedEvent{MailTemplateChangedEvent: *e.(*policy.MailTemplateChangedEvent)}, nil
}
