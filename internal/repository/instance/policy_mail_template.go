package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

var (
	MailTemplateAddedEventType   = instanceEventTypePrefix + policy.MailTemplatePolicyAddedEventType
	MailTemplateChangedEventType = instanceEventTypePrefix + policy.MailTemplatePolicyChangedEventType
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

func MailTemplateAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
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

func MailTemplateChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e, err := policy.MailTemplateChangedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &MailTemplateChangedEvent{MailTemplateChangedEvent: *e.(*policy.MailTemplateChangedEvent)}, nil
}
