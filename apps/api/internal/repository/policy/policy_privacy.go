package policy

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	PrivacyPolicyAddedEventType   = "policy.privacy.added"
	PrivacyPolicyChangedEventType = "policy.privacy.changed"
	PrivacyPolicyRemovedEventType = "policy.privacy.removed"
)

type PrivacyPolicyAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TOSLink        string              `json:"tosLink,omitempty"`
	PrivacyLink    string              `json:"privacyLink,omitempty"`
	HelpLink       string              `json:"helpLink,omitempty"`
	SupportEmail   domain.EmailAddress `json:"supportEmail,omitempty"`
	DocsLink       string              `json:"docsLink,omitempty"`
	CustomLink     string              `json:"customLink,omitempty"`
	CustomLinkText string              `json:"customLinkText,omitempty"`
}

func (e *PrivacyPolicyAddedEvent) Payload() interface{} {
	return e
}

func (e *PrivacyPolicyAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPrivacyPolicyAddedEvent(
	base *eventstore.BaseEvent,
	tosLink,
	privacyLink,
	helpLink string,
	supportEmail domain.EmailAddress,
	docsLink, customLink, customLinkText string,
) *PrivacyPolicyAddedEvent {
	return &PrivacyPolicyAddedEvent{
		BaseEvent:      *base,
		TOSLink:        tosLink,
		PrivacyLink:    privacyLink,
		HelpLink:       helpLink,
		SupportEmail:   supportEmail,
		DocsLink:       docsLink,
		CustomLink:     customLink,
		CustomLinkText: customLinkText,
	}
}

func PrivacyPolicyAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PrivacyPolicyAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-2k0fs", "unable to unmarshal policy")
	}

	return e, nil
}

type PrivacyPolicyChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	TOSLink        *string              `json:"tosLink,omitempty"`
	PrivacyLink    *string              `json:"privacyLink,omitempty"`
	HelpLink       *string              `json:"helpLink,omitempty"`
	SupportEmail   *domain.EmailAddress `json:"supportEmail,omitempty"`
	DocsLink       *string              `json:"docsLink,omitempty"`
	CustomLink     *string              `json:"customLink,omitempty"`
	CustomLinkText *string              `json:"customLinkText,omitempty"`
}

func (e *PrivacyPolicyChangedEvent) Payload() interface{} {
	return e
}

func (e *PrivacyPolicyChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPrivacyPolicyChangedEvent(
	base *eventstore.BaseEvent,
	changes []PrivacyPolicyChanges,
) (*PrivacyPolicyChangedEvent, error) {
	if len(changes) == 0 {
		return nil, zerrors.ThrowPreconditionFailed(nil, "POLICY-PPo0s", "Errors.NoChangesFound")
	}
	changeEvent := &PrivacyPolicyChangedEvent{
		BaseEvent: *base,
	}
	for _, change := range changes {
		change(changeEvent)
	}
	return changeEvent, nil
}

type PrivacyPolicyChanges func(*PrivacyPolicyChangedEvent)

func ChangeTOSLink(tosLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.TOSLink = &tosLink
	}
}

func ChangePrivacyLink(privacyLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.PrivacyLink = &privacyLink
	}
}

func ChangeHelpLink(helpLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.HelpLink = &helpLink
	}
}

func ChangeSupportEmail(supportEmail domain.EmailAddress) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.SupportEmail = &supportEmail
	}
}

func ChangeDocsLink(docsLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.DocsLink = &docsLink
	}
}

func ChangeCustomLink(customLink string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.CustomLink = &customLink
	}
}

func ChangeCustomLinkText(customLinkText string) func(*PrivacyPolicyChangedEvent) {
	return func(e *PrivacyPolicyChangedEvent) {
		e.CustomLinkText = &customLinkText
	}
}

func PrivacyPolicyChangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &PrivacyPolicyChangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "POLIC-22nf9", "unable to unmarshal policy")
	}

	return e, nil
}

type PrivacyPolicyRemovedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *PrivacyPolicyRemovedEvent) Payload() interface{} {
	return nil
}

func (e *PrivacyPolicyRemovedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewPrivacyPolicyRemovedEvent(base *eventstore.BaseEvent) *PrivacyPolicyRemovedEvent {
	return &PrivacyPolicyRemovedEvent{
		BaseEvent: *base,
	}
}

func PrivacyPolicyRemovedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &PrivacyPolicyRemovedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
