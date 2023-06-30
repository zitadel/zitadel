package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/authrequest"
)

type AuthRequestWriteModel struct {
	eventstore.WriteModel

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
	LoginHint        string
	HintUserID       string
	AuthRequestState domain.AuthRequestState
	SessionID        string
}

func NewAuthRequestWriteModel(id string) *AuthRequestWriteModel {
	return &AuthRequestWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: id,
		},
	}
}

func (m *AuthRequestWriteModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *authrequest.AddedEvent:
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
		case *authrequest.CodeAddedEvent:
			// TODO: left fold fields ASAP
			m.AuthRequestState = domain.AuthRequestStateCodeAdded
		case *authrequest.SessionLinkedEvent:
			m.SessionID = e.SessionID
		case *authrequest.FailedEvent:
			m.AuthRequestState = domain.AuthRequestStateFailed
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
