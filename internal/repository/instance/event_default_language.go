package instance

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	DefaultLanguageSetEventType eventstore.EventType = "instance.default.language.set"
)

type DefaultLanguageSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Language language.Tag `json:"language"`
}

func (e *DefaultLanguageSetEvent) Payload() interface{} {
	return e
}

func (e *DefaultLanguageSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewDefaultLanguageSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	language language.Tag,
) *DefaultLanguageSetEvent {
	return &DefaultLanguageSetEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			DefaultLanguageSetEventType,
		),
		Language: language,
	}
}

func DefaultLanguageSetMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &DefaultLanguageSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "IAM-3j9fs", "unable to unmarshal default language set")
	}

	return e, nil
}
