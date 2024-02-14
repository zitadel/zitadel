package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TokenKeyTokenID        = "id"
	TokenKeyUserID         = "user_id"
	TokenKeyRefreshTokenID = "refresh_token_id"
	TokenKeyApplicationID  = "application_id"
	TokenKeyUserAgentID    = "user_agent_id"
	TokenKeyExpiration     = "expiration"
	TokenKeyResourceOwner  = "resource_owner"
	TokenKeyInstanceID     = "instance_id"
)

type TokenView struct {
	ID                string                     `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate      time.Time                  `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time                  `json:"-" gorm:"column:change_date"`
	ResourceOwner     string                     `json:"-" gorm:"column:resource_owner"`
	UserID            string                     `json:"-" gorm:"column:user_id"`
	ApplicationID     string                     `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          database.TextArray[string] `json:"audience" gorm:"column:audience"`
	Scopes            database.TextArray[string] `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time                  `json:"expiration" gorm:"column:expiration"`
	Sequence          uint64                     `json:"-" gorm:"column:sequence"`
	PreferredLanguage string                     `json:"preferredLanguage" gorm:"column:preferred_language"`
	RefreshTokenID    string                     `json:"refreshTokenID,omitempty" gorm:"refresh_token_id"`
	IsPAT             bool                       `json:"-" gorm:"is_pat"`
	Deactivated       bool                       `json:"-" gorm:"-"`
	InstanceID        string                     `json:"instanceID" gorm:"column:instance_id;primary_key"`
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
		RefreshTokenID:    token.RefreshTokenID,
		IsPAT:             token.IsPAT,
	}
}

func (t *TokenView) AppendEventIfMyToken(event eventstore.Event) (err error) {
	view := new(TokenView)
	switch event.Type() {
	case user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		view.setRootData(event)
		err = view.setData(event)
	case user_repo.UserTokenRemovedType:
		return t.appendTokenRemoved(event)
	case user_repo.HumanRefreshTokenRemovedType:
		return t.appendRefreshTokenRemoved(event)
	case user_repo.UserV1SignedOutType,
		user_repo.HumanSignedOutType:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		if t.UserAgentID == id {
			t.Deactivated = true
		}
		return nil
	case user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType:
		t.Deactivated = true
		return nil
	case user_repo.UserUnlockedType,
		user_repo.UserReactivatedType:
		if t.ID != "" && event.CreatedAt().Before(t.CreationDate) {
			t.Deactivated = false
		}
		return nil
	case user_repo.PersonalAccessTokenRemovedType:
		return t.appendPATRemoved(event)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	if view.ID == t.ID {
		return t.AppendEvent(event)
	}
	return nil
}

func (t *TokenView) AppendEvent(event eventstore.Event) error {
	t.ChangeDate = event.CreatedAt()
	t.Sequence = event.Sequence()
	switch event.Type() {
	case user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		t.setRootData(event)
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreatedAt()
		t.IsPAT = event.Type() == user_repo.PersonalAccessTokenAddedType
	}
	return nil
}

func (t *TokenView) setRootData(event eventstore.Event) {
	t.UserID = event.Aggregate().ID
	t.ResourceOwner = event.Aggregate().ResourceOwner
	t.InstanceID = event.Aggregate().InstanceID
}

func (t *TokenView) setData(event eventstore.Event) error {
	if err := event.Unmarshal(t); err != nil {
		logging.Log("EVEN-3Gm9s").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}

func agentIDFromSession(event eventstore.Event) (string, error) {
	session := make(map[string]interface{})
	if err := event.Unmarshal(&session); err != nil {
		logging.Log("EVEN-Ghgt3").WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-GBf32", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}

func (t *TokenView) appendTokenRemoved(event eventstore.Event) error {
	token, err := eventToMap(event)
	if err != nil {
		return err
	}
	if token["tokenId"] == t.ID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendRefreshTokenRemoved(event eventstore.Event) error {
	refreshToken, err := eventToMap(event)
	if err != nil {
		return err
	}
	if refreshToken["tokenId"] == t.RefreshTokenID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendPATRemoved(event eventstore.Event) error {
	pat, err := eventToMap(event)
	if err != nil {
		return err
	}
	if pat["tokenId"] == t.ID && t.IsPAT {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) GetRelevantEventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType,
		user_repo.UserTokenRemovedType,
		user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserV1SignedOutType,
		user_repo.HumanSignedOutType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType,
		user_repo.UserLockedType,
		user_repo.UserReactivatedType,
		user_repo.PersonalAccessTokenRemovedType,
	}
}

func eventToMap(event eventstore.Event) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := event.Unmarshal(&m); err != nil {
		logging.Log("EVEN-Dbffe").WithError(err).Error("could not unmarshal event data")
		return nil, zerrors.ThrowInternal(nil, "MODEL-SDAfw", "could not unmarshal data")
	}
	return m, nil
}
