package policy

import (
	"fmt"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UniqueMailText                 = "mail_text"
	mailTextPolicyPrefix           = mailPolicyPrefix + "text."
	MailTextPolicyAddedEventType   = mailTextPolicyPrefix + "added"
	MailTextPolicyChangedEventType = mailTextPolicyPrefix + "changed"
	MailTextPolicyRemovedEventType = mailTextPolicyPrefix + "removed"
)

func NewAddMailTextUniqueConstraint(aggregateID, mailTextType, langugage string) *eventstore.UniqueConstraint {
	return eventstore.NewAddEventUniqueConstraint(
		UniqueMailText,
		fmt.Sprintf("%v:%v:%v", aggregateID, mailTextType, langugage),
		"Errors.Org.AlreadyExists")
}

func NewRemoveMailTextUniqueConstraint(aggregateID, mailTextType, langugage string) *eventstore.UniqueConstraint {
	return eventstore.NewRemoveUniqueConstraint(
		UniqueMailText,
		fmt.Sprintf("%v:%v:%v", aggregateID, mailTextType, langugage))
}

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

func (e *MailTextAddedEvent) Payload() interface{} {
	return e
}

func (e *MailTextAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewAddMailTextUniqueConstraint(e.Aggregate().ResourceOwner, e.MailTextType, e.Language)}
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

func MailTextAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MailTextAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-5m9if", "unable to unmarshal mail text policy")
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

func (e *MailTextChangedEvent) Payload() interface{} {
	return e
}

func (e *MailTextChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewMailTextChangedEvent(
	base *eventstore.BaseEvent,
	mailTextType,
	language string,
	changes []MailTextChanges,
) (*MailTextChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-m9osd", "Errors.NoChangesFound")
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

func MailTextChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &MailTextChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-bn88u", "unable to unmarshal mail text policy")
	}

	return e, nil
}

type MailTextRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`

	MailTextType string `json:"mailTextType,omitempty"`
	Language     string `json:"language,omitempty"`
}

func (e *MailTextRemovedEvent) Payload() interface{} {
	return nil
}

func (e *MailTextRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return []*eventstore.UniqueConstraint{NewRemoveMailTextUniqueConstraint(e.Aggregate().ResourceOwner, e.MailTextType, e.Language)}
}

func NewMailTextRemovedEvent(base *eventstore.BaseEvent, mailTextType, language string) *MailTextRemovedEvent {
	return &MailTextRemovedEvent{
		BaseEvent:    *base,
		MailTextType: mailTextType,
		Language:     language,
	}
}

func MailTextRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &MailTextRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
