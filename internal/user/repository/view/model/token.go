package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"

	caos_errs "github.com/caos/zitadel/internal/errors"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/lib/pq"
)

const (
	TokenKeyTokenID       = "id"
	TokenKeyUserID        = "user_id"
	TokenKeyApplicationID = "application_id"
	TokenKeyUserAgentID   = "user_agent_id"
	TokenKeyExpiration    = "expiration"
	TokenKeyResourceOwner = "resource_owner"
)

type TokenView struct {
	ID                string         `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate      time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner     string         `json:"-" gorm:"column:resource_owner"`
	UserID            string         `json:"-" gorm:"column:user_id"`
	ApplicationID     string         `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string         `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          pq.StringArray `json:"audience" gorm:"column:audience"`
	Scopes            pq.StringArray `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time      `json:"expiration" gorm:"column:expiration"`
	Sequence          uint64         `json:"-" gorm:"column:sequence"`
	PreferredLanguage string         `json:"preferredLanguage" gorm:"column:preferred_language"`
	Deactivated       bool           `json:"-" gorm:"-"`
}

func TokenViewFromModel(token *usr_model.TokenView) *TokenView {
	return &TokenView{
		ID:                token.ID,
		CreationDate:      token.CreationDate,
		ChangeDate:        token.ChangeDate,
		ResourceOwner:     token.ResourceOwner,
		UserID:            token.UserID,
		ApplicationID:     token.ApplicationID,
		UserAgentID:       token.UserAgentID,
		Audience:          token.Audience,
		Scopes:            token.Scopes,
		Expiration:        token.Expiration,
		Sequence:          token.Sequence,
		PreferredLanguage: token.PreferredLanguage,
	}
}

func TokenViewToModel(token *TokenView) *usr_model.TokenView {
	return &usr_model.TokenView{
		ID:                token.ID,
		CreationDate:      token.CreationDate,
		ChangeDate:        token.ChangeDate,
		ResourceOwner:     token.ResourceOwner,
		UserID:            token.UserID,
		ApplicationID:     token.ApplicationID,
		UserAgentID:       token.UserAgentID,
		Audience:          token.Audience,
		Scopes:            token.Scopes,
		Expiration:        token.Expiration,
		Sequence:          token.Sequence,
		PreferredLanguage: token.PreferredLanguage,
	}
}

func (t *TokenView) AppendEventIfMyToken(event *es_models.Event) (err error) {
	view := new(TokenView)
	switch event.Type {
	case usr_es_model.UserTokenAdded:
		view.setRootData(event)
		err = view.setData(event)
	case usr_es_model.SignedOut,
		usr_es_model.HumanSignedOut:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		if t.UserAgentID == id {
			t.Deactivated = true
		}
		return nil
	case usr_es_model.UserRemoved,
		usr_es_model.UserDeactivated,
		usr_es_model.UserLocked:
		t.Deactivated = true
		return nil
	case usr_es_model.UserUnlocked,
		usr_es_model.UserReactivated:
		if t.ID != "" && event.CreationDate.Before(t.CreationDate) {
			t.Deactivated = false
		}
		return nil
	default:
		return nil
	}
	if view.ID == t.ID {
		return t.AppendEvent(event)
	}
	return nil
}

func (t *TokenView) AppendEvent(event *es_models.Event) error {
	t.ChangeDate = event.CreationDate
	t.Sequence = event.Sequence
	switch event.Type {
	case usr_es_model.UserTokenAdded:
		t.setRootData(event)
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreationDate
	}
	return nil
}

func (t *TokenView) setRootData(event *es_models.Event) {
	t.UserID = event.AggregateID
	t.ResourceOwner = event.ResourceOwner
}

func (t *TokenView) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, t); err != nil {
		logging.Log("EVEN-3Gm9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}

func agentIDFromSession(event *es_models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &session); err != nil {
		logging.Log("EVEN-Ghgt3").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-GBf32", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}
