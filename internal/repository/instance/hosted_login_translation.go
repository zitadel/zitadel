package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	HostedLoginTranslationSet = instanceEventTypePrefix + "hosted_login_translation.set"
)

type HostedLoginTranslationSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Translation map[string]any `json:"translation,omitempty"`
	Language    string         `json:"language,omitempty"`
	Level       string         `json:"level,omitempty"`
	LevelID     string         `json:"level_id,omitempty"`
}

func NewHostedLoginTranslationSetEvent(ctx context.Context, aggregate *eventstore.Aggregate, translation map[string]any, language string) *HostedLoginTranslationSetEvent {
	return &HostedLoginTranslationSetEvent{
		BaseEvent:   *eventstore.NewBaseEventForPush(ctx, aggregate, HostedLoginTranslationSet),
		Translation: translation,
		Language:    language,
		Level:       string(aggregate.Type),
		LevelID:     aggregate.ID,
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
		return nil, zerrors.ThrowInternal(err, "INST-lOxtJJ", "unable to unmarshal hosted login translation set event")
	}

	return translationSet, nil
}
