package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
)

const (
	RefreshTokenKeyTokenID       = "id"
	RefreshTokenKeyUserID        = "user_id"
	RefreshTokenKeyApplicationID = "application_id"
	RefreshTokenKeyUserAgentID   = "user_agent_id"
	RefreshTokenKeyExpiration    = "expiration"
	RefreshTokenKeyResourceOwner = "resource_owner"
	RefreshTokenKeyInstanceID    = "instance_id"
)

type RefreshTokenView struct {
	ID                    string                     `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate          time.Time                  `json:"-" gorm:"column:creation_date"`
	ChangeDate            time.Time                  `json:"-" gorm:"column:change_date"`
	ResourceOwner         string                     `json:"-" gorm:"column:resource_owner"`
	Token                 string                     `json:"-" gorm:"column:token"`
	UserID                string                     `json:"-" gorm:"column:user_id"`
	ClientID              string                     `json:"clientID" gorm:"column:client_id"`
	UserAgentID           string                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience              database.TextArray[string] `json:"audience" gorm:"column:audience"`
	Scopes                database.TextArray[string] `json:"scopes" gorm:"column:scopes"`
	AuthMethodsReferences database.TextArray[string] `json:"authMethodsReference" gorm:"column:amr"`
	AuthTime              time.Time                  `json:"authTime" gorm:"column:auth_time"`
	IdleExpiration        time.Time                  `json:"-" gorm:"column:idle_expiration"`
	Expiration            time.Time                  `json:"-" gorm:"column:expiration"`
	Sequence              uint64                     `json:"-" gorm:"column:sequence"`
	InstanceID            string                     `json:"instanceID" gorm:"column:instance_id;primary_key"`
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

func (t *RefreshTokenView) AppendEventIfMyRefreshToken(event eventstore.Event) (err error) {
	view := new(RefreshTokenView)
	switch event.Type() {
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

func (t *RefreshTokenView) AppendEvent(event eventstore.Event) error {
	t.ChangeDate = event.CreatedAt()
	t.Sequence = event.Sequence()
	switch event.Type() {
	case user_repo.HumanRefreshTokenAddedType:
		t.setRootData(event)
		return t.appendAddedEvent(event)
	case user_repo.HumanRefreshTokenRenewedType:
		t.setRootData(event)
		return t.appendRenewedEvent(event)
	}
	return nil
}

func (t *RefreshTokenView) setRootData(event eventstore.Event) {
	t.UserID = event.Aggregate().ID
	t.ResourceOwner = event.Aggregate().ResourceOwner
	t.InstanceID = event.Aggregate().InstanceID
}

func (t *RefreshTokenView) appendAddedEvent(event eventstore.Event) error {
	e := new(user_repo.HumanRefreshTokenAddedEvent)
	if err := event.Unmarshal(e); err != nil {
		logging.Log("EVEN-Dbb31").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bbr42", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.CreationDate = event.CreatedAt()
	t.AuthMethodsReferences = e.AuthMethodsReferences
	t.AuthTime = e.AuthTime
	t.Audience = e.Audience
	t.ClientID = e.ClientID
	t.Expiration = event.CreatedAt().Add(e.Expiration)
	t.IdleExpiration = event.CreatedAt().Add(e.IdleExpiration)
	t.Scopes = e.Scopes
	t.Token = e.TokenID
	t.UserAgentID = e.UserAgentID
	return nil
}

func (t *RefreshTokenView) appendRenewedEvent(event eventstore.Event) error {
	e := new(user_repo.HumanRefreshTokenRenewedEvent)
	if err := event.Unmarshal(e); err != nil {
		logging.Log("EVEN-Vbbn2").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-Bbrn4", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.IdleExpiration = event.CreatedAt().Add(e.IdleExpiration)
	t.Token = e.RefreshToken
	return nil
}

func (t *RefreshTokenView) appendRemovedEvent(event eventstore.Event) {
	t.Expiration = event.CreatedAt()
}

func (t *RefreshTokenView) GetRelevantEventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user_repo.HumanRefreshTokenAddedType,
		user_repo.HumanRefreshTokenRenewedType,
		user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType,
	}
}
