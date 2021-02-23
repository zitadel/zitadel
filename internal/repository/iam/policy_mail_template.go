package iam

import (
	"context"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
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
	aggregate *eventstore.Aggregate,
	template []byte,
) *MailTemplateAddedEvent {
	return &MailTemplateAddedEvent{
		MailTemplateAddedEvent: *policy.NewMailTemplateAddedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, MailTemplateAddedEventType),
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
	aggregate *eventstore.Aggregate,
	changes []policy.MailTemplateChanges,
) (*MailTemplateChangedEvent, error) {
	changedEvent, err := policy.NewMailTemplateChangedEvent(
		eventstore.NewBaseEventForPush(ctx, aggregate, MailTemplateChangedEventType),
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
