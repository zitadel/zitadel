package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	"github.com/lib/pq"
)

const (
	RefreshTokenKeyTokenID       = "id"
	RefreshTokenKeyUserID        = "user_id"
	RefreshTokenKeyApplicationID = "application_id"
	RefreshTokenKeyUserAgentID   = "user_agent_id"
	RefreshTokenKeyExpiration    = "expiration"
	RefreshTokenKeyResourceOwner = "resource_owner"
)

type RefreshTokenView struct {
	ID                string         `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate      time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner     string         `json:"-" gorm:"column:resource_owner"`
	UserID            string         `json:"-" gorm:"column:user_id"`
	ApplicationID     string         `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string         `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          pq.StringArray `json:"audience" gorm:"column:audience"`
	Scopes            pq.StringArray `json:"scopes" gorm:"column:scopes"`
	IdleExpiration    time.Time      `json:"idle_expiration" gorm:"column:expiration"`
	Expiration        time.Time      `json:"expiration" gorm:"column:expiration"`
	Sequence          uint64         `json:"-" gorm:"column:sequence"`
	Token             string         `json:"-" gorm:"column:token"`
	PreferredLanguage string         `json:"preferredLanguage" gorm:"column:preferred_language"`
}

func RefreshTokenViewFromModel(token *usr_model.RefreshTokenView) *RefreshTokenView {
	return &RefreshTokenView{
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

func RefreshTokenViewToModel(token *RefreshTokenView) *usr_model.RefreshTokenView {
	return &usr_model.RefreshTokenView{
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

func (t *RefreshTokenView) AppendEventIfMyRefreshToken(event *es_models.Event) (err error) {
	view := new(RefreshTokenView)
	switch eventstore.EventType(event.Type) {
	case user_repo.UserRefreshTokenAddedType:
		view.setRootData(event)
		err = view.setData(event)
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

func (t *RefreshTokenView) AppendEvent(event *es_models.Event) error {
	t.ChangeDate = event.CreationDate
	t.Sequence = event.Sequence
	switch eventstore.EventType(event.Type) {
	case user_repo.UserRefreshTokenAddedType:
		t.setRootData(event)
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreationDate
	}
	return nil
}

func (t *RefreshTokenView) setRootData(event *es_models.Event) {
	t.UserID = event.AggregateID
	t.ResourceOwner = event.ResourceOwner
}

func (t *RefreshTokenView) setData(event *es_models.Event) error {
	if err := json.Unmarshal(event.Data, t); err != nil {
		logging.Log("EVEN-ADgn4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-BHgn3", "could not unmarshal event")
	}
	return nil
}
