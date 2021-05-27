package policy

import (
	"encoding/json"

	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
)

const (
	customTextPrefix           = mailPolicyPrefix + "customtext."
	CustomTextSetEventType     = customTextPrefix + "set"
	CustomTextRemovedEventType = customTextPrefix + "removed"
)

type CustomTextSetEvent struct {
	eventstore.BaseEvent `json:"-"`

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
	key,
	text string,
	language language.Tag,
) *CustomTextSetEvent {
	return &CustomTextSetEvent{
		BaseEvent: *base,
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

	Key      string       `json:"key,omitempty"`
	Language language.Tag `json:"language,omitempty"`
}

func (e *CustomTextRemovedEvent) Data() interface{} {
	return nil
}

func (e *CustomTextRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewCustomTextRemovedEvent(base *eventstore.BaseEvent, key string, language language.Tag) *CustomTextRemovedEvent {
	return &CustomTextRemovedEvent{
		BaseEvent: *base,
		Key:       key,
		Language:  language,
	}
}

func CustomTextRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &CustomTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
