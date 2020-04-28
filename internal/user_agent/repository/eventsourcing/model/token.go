package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user_agent/model"
)

type Token struct {
	es_models.ObjectRoot

	TokenID       string        `json:"tokenID"`
	UserSessionID string        `json:"userSessionID"`
	AuthSessionID string        `json:"authSessionID"`
	Expiry        time.Duration `json:"expiry,omitempty"`
}

func TokenFromModel(token *model.Token) *Token {
	return &Token{
		ObjectRoot:    token.ObjectRoot,
		TokenID:       token.TokenID,
		UserSessionID: token.UserSessionID,
		AuthSessionID: token.AuthSessionID,
	}
}

func TokenToModel(token *Token) *model.Token {
	return &model.Token{
		ObjectRoot:    token.ObjectRoot,
		TokenID:       token.TokenID,
		UserSessionID: token.UserSessionID,
		AuthSessionID: token.AuthSessionID,
	}
}

func (a *UserAgent) appendTokenAddedEvent(event *es_models.Event) error {
	token := new(Token)
	token.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, token); err != nil {
		logging.Log("MODEL-452wa").WithError(err).Debug("could not unmarshal event data")
		return err
	}
	if _, userSession := GetUserSession(a.UserSessions, token.UserSessionID); userSession != nil {
		if _, authSession := GetAuthSession(userSession.AuthSessions, token.AuthSessionID); authSession != nil {
			authSession.Token = token
		}
	}
	return nil
}
