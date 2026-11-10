package org

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	HostedLoginTranslationSet = orgEventTypePrefix + "hosted_login_translation.set"
)

type HostedLoginTranslationSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Translation map[string]any `json:"translation,omitempty"`
	Language    language.Tag   `json:"language,omitempty"`
	Level       string         `json:"level,omitempty"`
}

func NewHostedLoginTranslationSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, translation map[string]any, language language.Tag) *HostedLoginTranslationSetEvent {
	return &HostedLoginTranslationSetEvent{
		BaseEvent:   *eventstore.NewBaseEventForPush(ctx, aggregate, HostedLoginTranslationSet),
		Translation: translation,
		Language:    language,
		Level:       string(aggregate.Type),
	}
}

func (e *HostedLoginTranslationSetEvent) Payload() any {
	return e
}

func (e *HostedLoginTranslationSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *HostedLoginTranslationSetEvent) Fields() []*eventstore.FieldOperation {
	return nil
}

func HostedLoginTranslationSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	translationSet := &HostedLoginTranslationSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(translationSet)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "ORG-BH82Eb", "unable to unmarshal hosted login translation set event")
	}

	return translationSet, nil
}
