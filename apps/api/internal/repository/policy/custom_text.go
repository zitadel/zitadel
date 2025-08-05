package policy

import (
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	customTextPrefix                   = "customtext."
	CustomTextSetEventType             = customTextPrefix + "set"
	CustomTextRemovedEventType         = customTextPrefix + "removed"
	CustomTextTemplateRemovedEventType = customTextPrefix + "template.removed"
)

type CustomTextSetEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
	Text     string       `json:"text,omitempty"`
}

func (e *CustomTextSetEvent) Payload() interface{} {
	return e
}

func (e *CustomTextSetEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func CustomTextSetEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &CustomTextSetEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TEXT-28dwe", "unable to unmarshal custom text")
	}

	return e, nil
}

type CustomTextRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextRemovedEvent) Payload() interface{} {
	return e
}

func (e *CustomTextRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
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

func CustomTextRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &CustomTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TEXT-28sMf", "unable to unmarshal custom text removed")
	}

	return e, nil
}

type CustomTextTemplateRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Template string       `json:"template,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextTemplateRemovedEvent) Payload() interface{} {
	return e
}

func (e *CustomTextTemplateRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCustomTextTemplateRemovedEvent(base *eventstore.BaseEvent, template string, language language.Tag) *CustomTextTemplateRemovedEvent {
	return &CustomTextTemplateRemovedEvent{
		BaseEvent: *base,
		Template:  template,
		Language:  language,
	}
}

func CustomTextTemplateRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &CustomTextTemplateRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "TEXT-mKKRs", "unable to unmarshal custom text message removed")
	}

	return e, nil
}
