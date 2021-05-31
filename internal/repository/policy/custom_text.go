package policy

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	customTextPrefix                  = "customtext."
	CustomTextSetEventType            = customTextPrefix + "set"
	CustomTextRemovedEventType        = customTextPrefix + "removed"
	CustomTextMessageRemovedEventType = customTextPrefix + "message.removed"
)

type CustomTextSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
	Text     string       `json:"text,omitempty"`
}

func (e *CustomTextSetEvent) Data() interface{} {
	return e
}

func (e *CustomTextSetEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextSetEvent(
	base *eventstore.BaseEvent,
	template,
	key,
	text string,
	language language.Tag,
) *CustomTextSetEvent {
	return &CustomTextSetEvent{
		BaseEvent: *base,
		Template:  template,
		Key:       key,
		Language:  language,
		Text:      text,
	}
}

func CustomTextSetEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &CustomTextSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-28dwe", "unable to unmarshal custom text")
	}

	return e, nil
}

type CustomTextRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextRemovedEvent) Data() interface{} {
	return nil
}

func (e *CustomTextRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextRemovedEvent(base *eventstore.BaseEvent, template, key string, language language.Tag) *CustomTextRemovedEvent {
	return &CustomTextRemovedEvent{
		BaseEvent: *base,
		Template:  template,
		Key:       key,
		Language:  language,
	}
}

func CustomTextRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CustomTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type CustomTextMessageRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextMessageRemovedEvent) Data() interface{} {
	return nil
}

func (e *CustomTextMessageRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextMessageRemovedEvent(base *eventstore.BaseEvent, template string, language language.Tag) *CustomTextMessageRemovedEvent {
	return &CustomTextMessageRemovedEvent{
		BaseEvent: *base,
		Template:  template,
		Language:  language,
	}
}

func CustomTextMessageRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CustomTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
