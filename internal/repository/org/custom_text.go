package org

import (
	"context"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/repository/policy"
)

var (
	CustomTextSetEventType             = orgEventTypePrefix + policy.CustomTextSetEventType
	CustomTextRemovedEventType         = orgEventTypePrefix + policy.CustomTextRemovedEventType
	CustomTextTemplateRemovedEventType = orgEventTypePrefix + policy.CustomTextTemplateRemovedEventType
)

type CustomTextSetEvent struct {
	policy.CustomTextSetEvent
}

func NewCustomTextSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template,
	key,
	text string,
	language language.Tag,
) *CustomTextSetEvent {
	return &CustomTextSetEvent{
		CustomTextSetEvent: *policy.NewCustomTextSetEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, CustomTextSetEventType),
			template,
			key,
			text,
			language),
	}
}

func CustomTextSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.CustomTextSetEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &CustomTextSetEvent{CustomTextSetEvent: *e.(*policy.CustomTextSetEvent)}, nil
}

type CustomTextRemovedEvent struct {
	policy.CustomTextRemovedEvent
}

func NewCustomTextRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template,
	key string,
	language language.Tag,
) *CustomTextRemovedEvent {
	return &CustomTextRemovedEvent{
		CustomTextRemovedEvent: *policy.NewCustomTextRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, CustomTextRemovedEventType),
			template,
			key,
			language,
		),
	}
}

func CustomTextRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.CustomTextRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &CustomTextRemovedEvent{CustomTextRemovedEvent: *e.(*policy.CustomTextRemovedEvent)}, nil
}

type CustomTextTemplateRemovedEvent struct {
	policy.CustomTextTemplateRemovedEvent
}

func NewCustomTextTemplateRemovedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	template string,
	language language.Tag,
) *CustomTextTemplateRemovedEvent {
	return &CustomTextTemplateRemovedEvent{
		CustomTextTemplateRemovedEvent: *policy.NewCustomTextTemplateRemovedEvent(
			eventstore.NewBaseEventForPush(ctx, aggregate, CustomTextTemplateRemovedEventType),
			template,
			language,
		),
	}
}

func CustomTextTemplateRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e, err := policy.CustomTextTemplateRemovedEventMapper(event)
	if err != nil {
		return nil, err
	}

	return &CustomTextTemplateRemovedEvent{CustomTextTemplateRemovedEvent: *e.(*policy.CustomTextTemplateRemovedEvent)}, nil
}
