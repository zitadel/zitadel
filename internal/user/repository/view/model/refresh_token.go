package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	user_repo "github.com/caos/zitadel/internal/repository/user"
	usr_model "github.com/caos/zitadel/internal/user/model"
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
	ID                    string         `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate          time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate            time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner         string         `json:"-" gorm:"column:resource_owner"`
	Token                 string         `json:"-" gorm:"column:token"`
	UserID                string         `json:"-" gorm:"column:user_id"`
	ClientID              string         `json:"clientID" gorm:"column:client_id"`
	UserAgentID           string         `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience              pq.StringArray `json:"audience" gorm:"column:audience"`
	Scopes                pq.StringArray `json:"scopes" gorm:"column:scopes"`
	AuthMethodsReferences pq.StringArray `json:"authMethodsReference" gorm:"column:amr"`
	AuthTime              time.Time      `json:"authTime" gorm:"column:auth_time"`
	IdleExpiration        time.Time      `json:"-" gorm:"column:idle_expiration"`
	Expiration            time.Time      `json:"-" gorm:"column:expiration"`
	Sequence              uint64         `json:"-" gorm:"column:sequence"`
}

func RefreshTokenViewsToModel(tokens []*RefreshTokenView) []*usr_model.RefreshTokenView {
	result := make([]*usr_model.RefreshTokenView, len(tokens))
	for i, g := range tokens {
		result[i] = RefreshTokenViewToModel(g)
	}
	return result
}

func RefreshTokenViewToModel(token *RefreshTokenView) *usr_model.RefreshTokenView {
	return &usr_model.RefreshTokenView{
		ID:                    token.ID,
		CreationDate:          token.CreationDate,
		ChangeDate:            token.ChangeDate,
		ResourceOwner:         token.ResourceOwner,
		Token:                 token.Token,
		UserID:                token.UserID,
		ClientID:              token.ClientID,
		UserAgentID:           token.UserAgentID,
		Audience:              token.Audience,
		Scopes:                token.Scopes,
		AuthMethodsReferences: token.AuthMethodsReferences,
		AuthTime:              token.AuthTime,
		IdleExpiration:        token.IdleExpiration,
		Expiration:            token.Expiration,
		Sequence:              token.Sequence,
	}
}

func (t *RefreshTokenView) AppendEventIfMyRefreshToken(event *es_models.Event) (err error) {
	view := new(RefreshTokenView)
	switch eventstore.EventType(event.Type) {
	case user_repo.HumanRefreshTokenAddedType:
		view.setRootData(event)
		err = view.appendAddedEvent(event)
		if err != nil {
			return err
		}
	case user_repo.HumanRefreshTokenRenewedType:
		view.setRootData(event)
		err = view.appendRenewedEvent(event)
		if err != nil {
			return err
		}
	case user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType:
		view.appendRemovedEvent(event)
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
	case user_repo.HumanRefreshTokenAddedType:
		t.setRootData(event)
		return t.appendAddedEvent(event)
	case user_repo.HumanRefreshTokenRenewedType:
		t.setRootData(event)
		return t.appendRenewedEvent(event)
	}
	return nil
}

func (t *RefreshTokenView) setRootData(event *es_models.Event) {
	t.UserID = event.AggregateID
	t.ResourceOwner = event.ResourceOwner
}

func (t *RefreshTokenView) appendAddedEvent(event *es_models.Event) error {
	e := new(user_repo.HumanRefreshTokenAddedEvent)
	if err := json.Unmarshal(event.Data, e); err != nil {
		logging.Log("EVEN-Dbb31").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bbr42", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.CreationDate = event.CreationDate
	t.AuthMethodsReferences = e.AuthMethodsReferences
	t.AuthTime = e.AuthTime
	t.Audience = e.Audience
	t.ClientID = e.ClientID
	t.Expiration = event.CreationDate.Add(e.Expiration)
	t.IdleExpiration = event.CreationDate.Add(e.IdleExpiration)
	t.Scopes = e.Scopes
	t.Token = e.TokenID
	t.UserAgentID = e.UserAgentID
	return nil
}

func (t *RefreshTokenView) appendRenewedEvent(event *es_models.Event) error {
	e := new(user_repo.HumanRefreshTokenRenewedEvent)
	if err := json.Unmarshal(event.Data, e); err != nil {
		logging.Log("EVEN-Vbbn2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bbrn4", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.IdleExpiration = event.CreationDate.Add(e.IdleExpiration)
	t.Token = e.RefreshToken
	return nil
}

func (t *RefreshTokenView) appendRemovedEvent(event *es_models.Event) {
	t.Expiration = event.CreationDate
}
