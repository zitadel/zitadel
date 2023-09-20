package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
)

type AuthRequestWriteModel struct {
	eventstore.WriteModel
	aggregate *eventstore.Aggregate

	LoginClient      string
	ClientID         string
	RedirectURI      string
	State            string
	Nonce            string
	Scope            []string
	Audience         []string
	ResponseType     domain.OIDCResponseType
	CodeChallenge    *domain.OIDCCodeChallenge
	Prompt           []domain.Prompt
	UILocales        []string
	MaxAge           *time.Duration
	LoginHint        *string
	HintUserID       *string
	SessionID        string
	UserID           string
	AuthTime         time.Time
	AuthMethods      []domain.UserAuthMethodType
	AuthRequestState domain.AuthRequestState
}

func NewAuthRequestWriteModel(ctx context.Context, id string) *AuthRequestWriteModel {
	return &AuthRequestWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: id,
		},
		aggregate: &authrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
	}
}

func (m *AuthRequestWriteModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *authrequest.AddedEvent:
			m.LoginClient = e.LoginClient
			m.ClientID = e.ClientID
			m.RedirectURI = e.RedirectURI
			m.State = e.State
			m.Nonce = e.Nonce
			m.Scope = e.Scope
			m.Audience = e.Audience
			m.ResponseType = e.ResponseType
			m.CodeChallenge = e.CodeChallenge
			m.Prompt = e.Prompt
			m.UILocales = e.UILocales
			m.MaxAge = e.MaxAge
			m.LoginHint = e.LoginHint
			m.HintUserID = e.HintUserID
			m.AuthRequestState = domain.AuthRequestStateAdded
		case *authrequest.SessionLinkedEvent:
			m.SessionID = e.SessionID
			m.UserID = e.UserID
			m.AuthTime = e.AuthTime
			m.AuthMethods = e.AuthMethods
		case *authrequest.CodeAddedEvent:
			m.AuthRequestState = domain.AuthRequestStateCodeAdded
		case *authrequest.FailedEvent:
			m.AuthRequestState = domain.AuthRequestStateFailed
		case *authrequest.CodeExchangedEvent:
			m.AuthRequestState = domain.AuthRequestStateCodeExchanged
		case *authrequest.SucceededEvent:
			m.AuthRequestState = domain.AuthRequestStateSucceeded
		}
	}

	return m.WriteModel.Reduce()
}

func (m *AuthRequestWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(authrequest.AggregateType).
		AggregateIDs(m.AggregateID).
		Builder()
}

// CheckAuthenticated checks that the auth request exists, a session must have been linked
// and in case of a Code Flow the code must have been exchanged
func (m *AuthRequestWriteModel) CheckAuthenticated() error {
	if m.SessionID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "AUTHR-SF2r2", "Errors.AuthRequest.NotAuthenticated")
	}
	// in case of OIDC Code Flow, the code must have been exchanged
	if m.ResponseType == domain.OIDCResponseTypeCode && m.AuthRequestState == domain.AuthRequestStateCodeExchanged {
		return nil
	}
	// in case of OIDC Implicit Flow, check that the requests exists, but has not succeeded yet
	if (m.ResponseType == domain.OIDCResponseTypeIDToken || m.ResponseType == domain.OIDCResponseTypeIDTokenToken) &&
		m.AuthRequestState == domain.AuthRequestStateAdded {
		return nil
	}
	return caos_errs.ThrowPreconditionFailed(nil, "AUTHR-sajk3", "Errors.AuthRequest.NotAuthenticated")
}
