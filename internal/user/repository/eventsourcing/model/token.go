package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
)

type Token struct {
	es_models.ObjectRoot

	TokenID           string    `json:"tokenId" gorm:"column:token_id"`
	ApplicationID     string    `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string    `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          []string  `json:"audience" gorm:"column:audience"`
	Scopes            []string  `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time `json:"expiration" gorm:"column:expiration"`
	PreferredLanguage string    `json:"preferredLanguage" gorm:"column:preferred_language"`
}

func TokenFromModel(token *model.Token) *Token {
	return &Token{
		ObjectRoot:        token.ObjectRoot,
		TokenID:           token.TokenID,
		ApplicationID:     token.ApplicationID,
		UserAgentID:       token.UserAgentID,
		Audience:          token.Audience,
		Scopes:            token.Scopes,
		Expiration:        token.Expiration,
		PreferredLanguage: token.PreferredLanguage,
	}
}

func TokenToModel(token *Token) *model.Token {
	return &model.Token{
		ObjectRoot:        token.ObjectRoot,
		TokenID:           token.TokenID,
		ApplicationID:     token.ApplicationID,
		UserAgentID:       token.UserAgentID,
		Audience:          token.Audience,
		Scopes:            token.Scopes,
		Expiration:        token.Expiration,
		PreferredLanguage: token.PreferredLanguage,
	}
}

func (t *Token) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := t.AppendEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (t *Token) AppendEvent(event *es_models.Event) error {
	switch event.Type {
	case UserTokenAdded:
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreationDate
	}
	return nil
}

func (t *Token) setData(event *es_models.Event) error {
	t.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, t); err != nil {
		logging.Log("EVEN-4Fm9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}
