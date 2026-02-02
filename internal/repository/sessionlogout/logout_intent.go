package sessionlogout

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	federatedLogoutEventTypePrefix = "session.logout."
	StartedEventType               = federatedLogoutEventTypePrefix + "started"
	SAMLRequestCreatedEventType    = federatedLogoutEventTypePrefix + "saml.request.created"
	SAMLResponseReceivedEventType  = federatedLogoutEventTypePrefix + "saml.response.received"
	CompletedEventType             = federatedLogoutEventTypePrefix + "completed"
	FailedEventType                = federatedLogoutEventTypePrefix + "failed"
)

// StartedEvent is emitted when a federated logout is initiated
type StartedEvent struct {
	eventstore.BaseEvent `json:"-"`

	SessionID             string `json:"sessionId"`
	IDPID                 string `json:"idpId"`
	UserID                string `json:"userId"`
	PostLogoutRedirectURI string `json:"postLogoutRedirectUri"`
}

func NewStartedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	sessionID string,
	idpID string,
	userID string,
	postLogoutRedirectURI string,
) *StartedEvent {
	return &StartedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			StartedEventType,
		),
		SessionID:             sessionID,
		IDPID:                 idpID,
		UserID:                userID,
		PostLogoutRedirectURI: postLogoutRedirectURI,
	}
}

func (e *StartedEvent) Payload() any {
	return e
}

func (e *StartedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *StartedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func StartedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &StartedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "LOGOUT-Sf3g2", "unable to unmarshal event")
	}

	return e, nil
}

// SAMLRequestCreatedEvent is emitted when a SAML logout request is created
type SAMLRequestCreatedEvent struct {
	eventstore.BaseEvent `json:"-"`

	RequestID   string `json:"requestId"`
	BindingType string `json:"bindingType"` // "redirect" or "post"
	RedirectURL string `json:"redirectUrl,omitempty"`
	PostURL     string `json:"postUrl,omitempty"`
	SAMLRequest string `json:"samlRequest,omitempty"`
	RelayState  string `json:"relayState,omitempty"`
}

func NewSAMLRequestCreatedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	requestID string,
	bindingType string,
	redirectURL string,
	postURL string,
	samlRequest string,
	relayState string,
) *SAMLRequestCreatedEvent {
	return &SAMLRequestCreatedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLRequestCreatedEventType,
		),
		RequestID:   requestID,
		BindingType: bindingType,
		RedirectURL: redirectURL,
		PostURL:     postURL,
		SAMLRequest: samlRequest,
		RelayState:  relayState,
	}
}

func (e *SAMLRequestCreatedEvent) Payload() any {
	return e
}

func (e *SAMLRequestCreatedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SAMLRequestCreatedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func SAMLRequestCreatedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLRequestCreatedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "LOGOUT-kje3f", "unable to unmarshal event")
	}

	return e, nil
}

// SAMLResponseReceivedEvent is emitted when the IdP responds to the logout request
type SAMLResponseReceivedEvent struct {
	eventstore.BaseEvent `json:"-"`

	RequestID string `json:"requestId"`
}

func NewSAMLResponseReceivedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	requestID string,
) *SAMLResponseReceivedEvent {
	return &SAMLResponseReceivedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SAMLResponseReceivedEventType,
		),
		RequestID: requestID,
	}
}

func (e *SAMLResponseReceivedEvent) Payload() any {
	return e
}

func (e *SAMLResponseReceivedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *SAMLResponseReceivedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func SAMLResponseReceivedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &SAMLResponseReceivedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "LOGOUT-mke4f", "unable to unmarshal event")
	}

	return e, nil
}

// CompletedEvent is emitted when the logout is completed
type CompletedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewCompletedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
) *CompletedEvent {
	return &CompletedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CompletedEventType,
		),
	}
}

func (e *CompletedEvent) Payload() any {
	return e
}

func (e *CompletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *CompletedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func CompletedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &CompletedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "LOGOUT-nke5f", "unable to unmarshal event")
	}

	return e, nil
}

// FailedEvent is emitted when the logout fails
type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Reason string `json:"reason"`
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	reason string,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedEventType,
		),
		Reason: reason,
	}
}

func (e *FailedEvent) Payload() any {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *FailedEvent) SetBaseEvent(b *eventstore.BaseEvent) {
	e.BaseEvent = *b
}

func FailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	e := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}

	err := event.Unmarshal(e)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "LOGOUT-pke6f", "unable to unmarshal event")
	}

	return e, nil
}
