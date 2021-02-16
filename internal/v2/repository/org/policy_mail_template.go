package org

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/v2/repository/policy"
)

var (
	MailTemplateAddedEventType   = orgEventTypePrefix + policy.MailTemplatePolicyAddedEventType
	MailTemplateChangedEventType = orgEventTypePrefix + policy.MailTemplatePolicyChangedEventType
	MailTemplateRemovedEventType = orgEventTypePrefix + policy.MailTemplatePolicyRemovedEventType
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

type MailTemplateRemovedEvent struct {
	policy.MailTemplateRemovedEvent
}

func NewMailTemplateRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *MailTemplateRemovedEvent {
	return &MailTemplateRemovedEvent{
		MailTemplateRemovedEvent: *policy.NewMailTemplateRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, MailTemplateRemovedEventType),
		),
	}
}

func MailTemplateRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.MailTemplateRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTemplateRemovedEvent{MailTemplateRemovedEvent: *e.(*policy.MailTemplateRemovedEvent)}, nil
}
