package authrequest

import (
	"context"
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
)

const (
	authRequestEventPrefix = "auth_request."
	AddedType              = authRequestEventPrefix + "added"
	CodeAddedType          = authRequestEventPrefix + "code.added"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	LoginClient   string                    `json:"login_client"`
	ClientID      string                    `json:"client_id"`
	RedirectURI   string                    `json:"redirect_uri"`
	State         string                    `json:"state,omitempty"`
	Nonce         string                    `json:"nonce,omitempty"`
	Scope         []string                  `json:"scope,omitempty"`
	Audience      []string                  `json:"audience,omitempty"`
	ResponseType  domain.OIDCResponseType   `json:"response_type,omitempty"`
	CodeChallenge *domain.OIDCCodeChallenge `json:"code_challenge,omitempty"`
	Prompt        []domain.Prompt           `json:"prompt,omitempty"`
	UILocales     []string                  `json:"ui_locales,omitempty"`
	MaxAge        *time.Duration            `json:"max_age,omitempty"`
	LoginHint     *string                   `json:"login_hint,omitempty"`
	HintUserID    *string                   `json:"hint_user_id,omitempty"`
}

func (e *AddedEvent) Data() interface{} {
	return e
}

func (e *AddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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
	codeChallenge *domain.OIDCCodeChallenge,
	prompt []domain.Prompt,
	uiLocales []string,
	maxAge *time.Duration,
	loginHint,
	hintUserID *string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		LoginClient:   loginClient,
		ClientID:      clientID,
		RedirectURI:   redirectURI,
		State:         state,
		Nonce:         nonce,
		Scope:         scope,
		Audience:      audience,
		ResponseType:  responseType,
		CodeChallenge: codeChallenge,
		Prompt:        prompt,
		UILocales:     uiLocales,
		MaxAge:        maxAge,
		LoginHint:     loginHint,
		HintUserID:    hintUserID,
	}
}

func AddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &AddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "AUTHR-DG4gn", "unable to unmarshal auth request added")
	}

	return added, nil
}

type CodeAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	//TODO: add necessary fields
}

func (e *CodeAddedEvent) Data() interface{} {
	return e
}

func (e *CodeAddedEvent) UniqueConstraints() []*eventstore.EventUniqueConstraint {
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

func CodeAddedEventMapper(event *repository.Event) (eventstore.Event, error) {
	added := &CodeAddedEvent{
		BaseEvent: *eventstore.BaseEventFromRepo(event),
	}
	err := json.Unmarshal(event.Data, added)
	if err != nil {
		return nil, errors.ThrowInternal(err, "AUTHR-Sfe3w", "unable to unmarshal auth request code added")
	}

	return added, nil
}
