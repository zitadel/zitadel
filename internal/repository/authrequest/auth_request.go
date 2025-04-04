package authrequest

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	authRequestEventPrefix = "auth_request."
	AddedType              = authRequestEventPrefix + "added"
	FailedType             = authRequestEventPrefix + "failed"
	CodeAddedType          = authRequestEventPrefix + "code.added"
	SessionLinkedType      = authRequestEventPrefix + "session.linked"
	CodeExchangedType      = authRequestEventPrefix + "code.exchanged"
	SucceededType          = authRequestEventPrefix + "succeeded"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	LoginClient      string                    `json:"login_client"`
	ClientID         string                    `json:"client_id"`
	RedirectURI      string                    `json:"redirect_uri"`
	State            string                    `json:"state,omitempty"`
	Nonce            string                    `json:"nonce,omitempty"`
	Scope            []string                  `json:"scope,omitempty"`
	Audience         []string                  `json:"audience,omitempty"`
	ResponseType     domain.OIDCResponseType   `json:"response_type,omitempty"`
	ResponseMode     domain.OIDCResponseMode   `json:"response_mode,omitempty"`
	CodeChallenge    *domain.OIDCCodeChallenge `json:"code_challenge,omitempty"`
	Prompt           []domain.Prompt           `json:"prompt,omitempty"`
	UILocales        []string                  `json:"ui_locales,omitempty"`
	MaxAge           *time.Duration            `json:"max_age,omitempty"`
	LoginHint        *string                   `json:"login_hint,omitempty"`
	HintUserID       *string                   `json:"hint_user_id,omitempty"`
	NeedRefreshToken bool                      `json:"need_refresh_token,omitempty"`
	Issuer           string                    `json:"issuer,omitempty"`
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
	clientID,
	redirectURI,
	state,
	nonce string,
	scope,
	audience []string,
	responseType domain.OIDCResponseType,
	responseMode domain.OIDCResponseMode,
	codeChallenge *domain.OIDCCodeChallenge,
	prompt []domain.Prompt,
	uiLocales []string,
	maxAge *time.Duration,
	loginHint,
	hintUserID *string,
	needRefreshToken bool,
	issuer string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		LoginClient:      loginClient,
		ClientID:         clientID,
		RedirectURI:      redirectURI,
		State:            state,
		Nonce:            nonce,
		Scope:            scope,
		Audience:         audience,
		ResponseType:     responseType,
		ResponseMode:     responseMode,
		CodeChallenge:    codeChallenge,
		Prompt:           prompt,
		UILocales:        uiLocales,
		MaxAge:           maxAge,
		LoginHint:        loginHint,
		HintUserID:       hintUserID,
		NeedRefreshToken: needRefreshToken,
		Issuer:           issuer,
	}
}

func AddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "AUTHR-DG4gn", "unable to unmarshal auth request added")
	}

	return added, nil
}

type SessionLinkedEvent struct {
	eventstore.BaseEvent `json:"-"`

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
		BaseEvent: *eventstore.NewBaseEventForPush(
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

func SessionLinkedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &SessionLinkedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "AUTHR-Sfe3w", "unable to unmarshal auth request session linked")
	}

	return added, nil
}

type FailedEvent struct {
	eventstore.BaseEvent `json:"-"`

	Reason domain.OIDCErrorReason `json:"reason,omitempty"`
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
	reason domain.OIDCErrorReason,
) *FailedEvent {
	return &FailedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			FailedType,
		),
		Reason: reason,
	}
}

func FailedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &FailedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "AUTHR-Sfe3w", "unable to unmarshal auth request session linked")
	}

	return added, nil
}

type CodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CodeAddedEvent) Payload() interface{} {
	return e
}

func (e *CodeAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCodeAddedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
) *CodeAddedEvent {
	return &CodeAddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CodeAddedType,
		),
	}
}

func CodeAddedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	added := &CodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := event.Unmarshal(added)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "AUTHR-Sfe3w", "unable to unmarshal auth request code added")
	}

	return added, nil
}

type CodeExchangedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func (e *CodeExchangedEvent) Payload() interface{} {
	return nil
}

func (e *CodeExchangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func NewCodeExchangedEvent(ctx context.Context,
	aggregate *eventstore.Aggregate,
) *CodeExchangedEvent {
	return &CodeExchangedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			CodeExchangedType,
		),
	}
}

func CodeExchangedEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &CodeExchangedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}

type SucceededEvent struct {
	eventstore.BaseEvent `json:"-"`
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
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			SucceededType,
		),
	}
}

func SucceededEventMapper(event eventstore.Event) (eventstore.Event, error) {
	return &SucceededEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}, nil
}
