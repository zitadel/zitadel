package samlrequest

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	samlRequestEventPrefix = "saml_request."
	AddedType              = samlRequestEventPrefix + "added"
	FailedType             = samlRequestEventPrefix + "failed"
	SessionLinkedType      = samlRequestEventPrefix + "session.linked"
	SucceededType          = samlRequestEventPrefix + "succeeded"
)

type AddedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	LoginClient    string `json:"login_client,omitempty"`
	ApplicationID  string `json:"application_id,omitempty"`
	ACSURL         string `json:"acs_url,omitempty"`
	RelayState     string `json:"relay_state,omitempty"`
	RequestID      string `json:"request_id,omitempty"`
	Binding        string `json:"binding,omitempty"`
	Issuer         string `json:"issuer,omitempty"`
	Destination    string `json:"destination,omitempty"`
	ResponseIssuer string `json:"response_issuer,omitempty"`
}

func (e *AddedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

func (e *AddedEvent) Payload() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	loginClient,
	applicationID string,
	acsURL string,
	relayState string,
	requestID string,
	binding string,
	issuer string,
	destination string,
	responseIssuer string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		LoginClient:    loginClient,
		ApplicationID:  applicationID,
		ACSURL:         acsURL,
		RelayState:     relayState,
		RequestID:      requestID,
		Binding:        binding,
		Issuer:         issuer,
		Destination:    destination,
		ResponseIssuer: responseIssuer,
	}
}

type SessionLinkedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	SessionID   string                      `json:"session_id"`
	UserID      string                      `json:"user_id"`
	AuthTime    time.Time                   `json:"auth_time"`
	AuthMethods []domain.UserAuthMethodType `json:"auth_methods"`
}

func (e *SessionLinkedEvent) Payload() interface{} {
	return e
}

func (e *SessionLinkedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSessionLinkedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
	sessionID,
	userID string,
	authTime time.Time,
	authMethods []domain.UserAuthMethodType,
) *SessionLinkedEvent {
	return &SessionLinkedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SessionLinkedType,
		),
		SessionID:   sessionID,
		UserID:      userID,
		AuthTime:    authTime,
		AuthMethods: authMethods,
	}
}

func (e *SessionLinkedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

type FailedEvent struct {
	*eventstore.BaseEvent `json:"-"`

	Reason domain.SAMLErrorReason `json:"reason,omitempty"`
}

func (e *FailedEvent) Payload() interface{} {
	return e
}

func (e *FailedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewFailedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	reason domain.SAMLErrorReason,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedType,
		),
		Reason: reason,
	}
}

func (e *FailedEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}

type SucceededEvent struct {
	*eventstore.BaseEvent `json:"-"`
}

func (e *SucceededEvent) Payload() interface{} {
	return nil
}

func (e *SucceededEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewSucceededEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
) *SucceededEvent {
	return &SucceededEvent{
		BaseEvent: eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SucceededType,
		),
	}
}

func (e *SucceededEvent) SetBaseEvent(event *eventstore.BaseEvent) {
	e.BaseEvent = event
}
