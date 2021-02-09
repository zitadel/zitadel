package policy

import (
	"encoding/json"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

const (
	mailTextPolicyPrefix           = mailPolicyPrefix + "text."
	MailTextPolicyAddedEventType   = mailTextPolicyPrefix + "added"
	MailTextPolicyChangedEventType = mailTextPolicyPrefix + "changed"
	MailTextPolicyRemovedEventType = mailTextPolicyPrefix + "removed"
)

type MailTextAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MailTextType string `json:"mailTextType,omitempty"`
	Language     string `json:"language,omitempty"`
	Title        string `json:"title,omitempty"`
	PreHeader    string `json:"preHeader,omitempty"`
	Subject      string `json:"subject,omitempty"`
	Greeting     string `json:"greeting,omitempty"`
	Text         string `json:"text,omitempty"`
	ButtonText   string `json:"buttonText,omitempty"`
}

func (e *MailTextAddedEvent) Data() interface{} {
	return e
}

func (e *MailTextAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMailTextAddedEvent(
	base *eventstore.BaseEvent,
	mailTextType,
	language,
	title,
	preHeader,
	subject,
	greeting,
	text,
	buttonText string,
) *MailTextAddedEvent {
	return &MailTextAddedEvent{
		BaseEvent:    *base,
		MailTextType: mailTextType,
		Language:     language,
		Title:        title,
		PreHeader:    preHeader,
		Subject:      subject,
		Greeting:     greeting,
		Text:         text,
		ButtonText:   buttonText,
	}
}

func MailTextAddedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MailTextAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-5m9if", "unable to unmarshal mail text policy")
	}

	return e, nil
}

type MailTextChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MailTextType string  `json:"mailTextType,omitempty"`
	Language     string  `json:"language,omitempty"`
	Title        *string `json:"title,omitempty"`
	PreHeader    *string `json:"preHeader,omitempty"`
	Subject      *string `json:"subject,omitempty"`
	Greeting     *string `json:"greeting,omitempty"`
	Text         *string `json:"text,omitempty"`
	ButtonText   *string `json:"buttonText,omitempty"`
}

func (e *MailTextChangedEvent) Data() interface{} {
	return e
}

func (e *MailTextChangedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMailTextChangedEvent(
	base *eventstore.BaseEvent,
	mailTextType,
	language string,
	changes []MailTextChanges,
) (*MailTextChangedEvent, error) {
	if len(changes) == 0 {
		return nil, errors.ThrowPreconditionFailed(nil, "POLICY-m9osd", "Errors.NoChangesFound")
	}
	changeEvent := &MailTextChangedEvent{
		BaseEvent:    *base,
		MailTextType: mailTextType,
		Language:     language,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type MailTextChanges func(*MailTextChangedEvent)

func ChangeTitle(title string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.Title = &title
	}
}

func ChangePreHeader(preHeader string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.PreHeader = &preHeader
	}
}

func ChangeSubject(greeting string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.Subject = &greeting
	}
}

func ChangeGreeting(greeting string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.Greeting = &greeting
	}
}

func ChangeText(text string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.Text = &text
	}
}

func ChangeButtonText(buttonText string) func(*MailTextChangedEvent) {
	return func(e *MailTextChangedEvent) {
		e.ButtonText = &buttonText
	}
}

func MailTextChangedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	e := &MailTextChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := json.Unmarshal(event.Data, e)
	if err != nil {
		return nil, errors.ThrowInternal(err, "POLIC-bn88u", "unable to unmarshal mail text policy")
	}

	return e, nil
}

type MailTextRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *MailTextRemovedEvent) Data() interface{} {
	return nil
}

func (e *MailTextRemovedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
	return nil
}

func NewMailTextRemovedEvent(base *eventstore.BaseEvent) *MailTextRemovedEvent {
	return &MailTextRemovedEvent{
		BaseEvent: *base,
	}
}

func MailTextRemovedEventMapper(event *repository.Event) (eventstore.EventReader, error) {
	return &MailTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
