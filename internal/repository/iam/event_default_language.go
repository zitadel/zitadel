package iam

import (
	"context"
	"encoding/json"

	"github.com/caos/zitadel/internal/eventstore"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	DefaultLanguageSetEventType eventstore.EventType = "iam.default.language.set"
)

type DefaultLanguageSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	DefaultLanguage language.Tag `json:"defaultLanguage"`
}

func (e *DefaultLanguageSetEvent) Data() interface{} {
	return e
}

func (e *DefaultLanguageSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewDefaultLanguageSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	defaultLanguage language.Tag,
) *DefaultLanguageSetEvent {
	return &DefaultLanguageSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DefaultLanguageSetEventType,
		),
		DefaultLanguage: defaultLanguage,
	}
}

func DefaultLanguageSetMapper(event *repository.Event) (eventstore.Event, error) {
	e := &DefaultLanguageSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-3j9fs", "unable to unmarshal default language set")
	}

	return e, nil
}
