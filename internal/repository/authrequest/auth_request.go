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
	sessionEventPrefix = "auth_request."
	AddedType          = sessionEventPrefix + "added"
)

type AddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	LoginClient        string                    `json:"login_client"`
	ClientID           string                    `json:"client_id"`
	RedirectURI        string                    `json:"redirect_uri"`
	State              string                    `json:"state,omitempty"`
	Nonce              string                    `json:"nonce,omitempty"`
	Scope              []string                  `json:"scope"`
	ResponseType       domain.OIDCResponseType   `json:"response_type,omitempty"`
	CodeChallenge      *domain.OIDCCodeChallenge `json:"code_challenge,omitempty"`
	Prompts            []domain.Prompt           `json:"prompts,omitempty"`
	UILocales          []string                  `json:"ui_locales,omitempty"`
	MaxAge             *time.Duration            `json:"max_age,omitempty"`
	LoginHint          string                    `json:"login_hint,omitempty"`
	IDTokenHintSubject string                    `json:"id_token_hint_subject,omitempty"`
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
	scope []string,
	responseType domain.OIDCResponseType,
	codeChallenge *domain.OIDCCodeChallenge,
	prompts []domain.Prompt,
	uiLocales []string,
	maxAge *time.Duration,
	loginHint,
	idTokenHintSubject string,
) *AddedEvent {
	return &AddedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			ctx,
			aggregate,
			AddedType,
		),
		LoginClient:        loginClient,
		ClientID:           clientID,
		RedirectURI:        redirectURI,
		State:              state,
		Nonce:              nonce,
		Scope:              scope,
		ResponseType:       responseType,
		CodeChallenge:      codeChallenge,
		Prompts:            prompts,
		UILocales:          uiLocales,
		MaxAge:             maxAge,
		LoginHint:          loginHint,
		IDTokenHintSubject: idTokenHintSubject,
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
